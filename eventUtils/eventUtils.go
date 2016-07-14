//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//       Unless required by applicable law or agreed to in writing, software
//       distributed under the License is distributed on an "AS IS" BASIS,
//       WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//       See the License for the specific language governing permissions and
//       limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package eventUtils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"models/events"
	"sort"
	"strconv"
	"strings"
	"time"
	"utils/dbutils"
	"utils/logging"
	"utils/typeConv"
)

type EventDetails struct {
	Enable      bool
	OwnerId     events.OwnerId
	OwnerName   string
	EventName   string
	Description string
	SrcObjName  string
}

type EventBase struct {
	EventDetails
	EvtId     events.EventId
	TimeStamp time.Time
}

type Event struct {
	EventBase
	SrcObjKey interface{}
}

var EventMap map[events.EventId]EventDetails

type FaultDetail struct {
	RaiseFault       bool
	ClearingEventId  int
	ClearingDaemonId int
	AlarmSeverity    string
}

type EventStruct struct {
	EventId     int
	EventName   string
	Description string
	SrcObjName  string
	EventEnable bool
	IsFault     bool
	Fault       FaultDetail
}

type DaemonEvent struct {
	DaemonId          int
	DaemonName        string
	DaemonEventEnable bool
	EventList         []EventStruct
}

type EventJson struct {
	DaemonEvents []DaemonEvent
}

type PubIntf interface {
	Publish(string, interface{}, interface{})
}

type KeyObj struct {
	Key   string
	UTime int64
}

type KeyObjSlice []KeyObj

var GlobalEventEnable bool = true
var Logger logging.LoggerIntf
var PubHdl PubIntf
var DbHdl dbutils.DBIntf

const (
	EventDir string = "/etc/flexswitch/"
)

func initEventDetails(ownerName string) error {
	var evtJson EventJson
	eventsFile := EventDir + "events.json"
	bytes, err := ioutil.ReadFile(eventsFile)
	if err != nil {
		Logger.Err(fmt.Sprintln("Error in reading ", eventsFile, " file."))
		err := errors.New(fmt.Sprintln("Error in reading ", eventsFile, " file."))
		return err
	}

	err = json.Unmarshal(bytes, &evtJson)
	if err != nil {
		Logger.Err(fmt.Sprintln("Errors in unmarshalling json file : ", eventsFile))
		err := errors.New(fmt.Sprintln("Errors in unmarshalling json file: ", eventsFile))
		return err
	}

	Logger.Debug(fmt.Sprintln("Owner Name :", ownerName, "evtJson:", evtJson))
	for _, daemon := range evtJson.DaemonEvents {
		Logger.Debug(fmt.Sprintln("OwnerName:", ownerName, "daemon.DaemonName:", daemon.DaemonName))
		if daemon.DaemonName == ownerName {
			GlobalEventEnable = daemon.DaemonEventEnable
			for _, evt := range daemon.EventList {
				evtId := events.EventId(evt.EventId)
				evtEnt, _ := EventMap[evtId]
				evtEnt.EventName = evt.EventName
				evtEnt.Description = evt.Description
				evtEnt.SrcObjName = evt.SrcObjName
				evtEnt.OwnerId = events.OwnerId(daemon.DaemonId)
				evtEnt.OwnerName = ownerName
				evtEnt.Enable = evt.EventEnable
				EventMap[evtId] = evtEnt
			}
			continue
		}
	}

	return nil
}

func InitEvents(ownerName string, dbHdl dbutils.DBIntf, pubHdl PubIntf, logger logging.LoggerIntf) error {

	EventMap = make(map[events.EventId]EventDetails)
	Logger = logger
	PubHdl = pubHdl
	DbHdl = dbHdl
	Logger.Info(fmt.Sprintln("Initializing Owner Name :", ownerName))
	err := initEventDetails(ownerName)
	if err != nil {
		return err
	}

	Logger.Info(fmt.Sprintln("EventMap:", EventMap))
	return nil
}

func PublishEvents(eventId events.EventId, key interface{}) error {
	var err error
	if GlobalEventEnable == false {
		return nil
	}
	evt := new(Event)
	evtEnt, exist := EventMap[eventId]
	if !exist {
		err := errors.New(fmt.Sprintln("Unable to find the event corresponding to given eventId: ", eventId))
		return err
	}

	if evtEnt.Enable == false {
		return nil
	}
	//Store raw event in DB
	evt.OwnerId = evtEnt.OwnerId
	evt.OwnerName = evtEnt.OwnerName
	evt.EvtId = eventId
	evt.EventName = evtEnt.EventName
	evt.TimeStamp = time.Now()
	evt.Description = evtEnt.Description
	evt.SrcObjName = evtEnt.SrcObjName
	evt.SrcObjKey = key
	msg, _ := json.Marshal(*evt)
	var unmarshalMsg Event
	err = json.Unmarshal(msg, &unmarshalMsg)
	keyMap, _ := events.EventKeyMap[evt.OwnerName]
	obj, _ := keyMap[evt.SrcObjName]
	obj = unmarshalMsg.SrcObjKey
	str := fmt.Sprintf("%v", obj)
	keyString := strings.TrimPrefix(str, "map[")
	strKey := strings.Split(keyString, "]")
	Logger.Info(fmt.Sprintln("Events to be published: ", evt, strKey[0]))
	keyStr := fmt.Sprintf("Events#%s#%s#%s#%s#%s#%d#", evt.OwnerName, evt.EventName, evt.SrcObjName, strKey[0], evt.TimeStamp.String(), evt.TimeStamp.UnixNano())
	Logger.Debug(fmt.Sprintln("Key Str :", keyStr))
	err = DbHdl.StoreValInDb(keyStr, evt.Description, "Desc")
	if err != nil {
		Logger.Err(fmt.Sprintln("Storing Events in database failed, err:", err))
	}
	//Store event stats in DB
	var statObj events.EventStats
	statObj.EventId = eventId
	dbObj, err := DbHdl.GetEventObjectFromDb(statObj, statObj.GetKey())
	if err != nil {
		//Event stat does not exist in db. Create one.
		statObj.EventName = evtEnt.EventName
		statObj.NumEvents = uint32(1)
		statObj.LastEventTime = evt.TimeStamp.String()
	} else {
		//Update DB entry
		statObj = dbObj.(events.EventStats)
		statObj.NumEvents += 1
		statObj.LastEventTime = evt.TimeStamp.String()
	}
	err = DbHdl.StoreEventObjectInDb(statObj)
	if err != nil {
		Logger.Err(fmt.Sprintln("Storing Event Stats in database failed, err:", err))
	}
	//Publish event
	PubHdl.Publish("PUBLISH", evt.OwnerName, msg)
	return nil
}

func GetEvents(evtObj events.EventObj, dbHdl dbutils.DBIntf, logger logging.LoggerIntf) (evt []events.EventObj, err error) {
	switch evtObj.(type) {
	case events.Event:
		qPattern := constructQueryPattern(evtObj.(events.Event))
		fmt.Println("Pattern Query:", qPattern)
		keys, err := typeConv.ConvertToStrings(dbHdl.GetAllKeys(qPattern))
		if err != nil {
			logger.Err(fmt.Sprintln("Error querying for keys:", err))
		}
		keySlice := constructKeySlice(keys)
		if keySlice == nil {
			logger.Err("Key slice is nil")
		}
		sort.Sort(KeyObjSlice(keySlice))
		for _, keyObj := range keySlice {
			desc, err := typeConv.ConvertToString(dbHdl.GetValFromDB(keyObj.Key, "Desc"))
			if err != nil {
				logger.Err(fmt.Sprintln("Error getting the value from DB", err))
				continue
			}
			str := strings.Split(keyObj.Key, "#")
			obj := events.Event{
				OwnerName:   str[1],
				EventName:   str[2],
				TimeStamp:   str[5],
				Description: desc,
				SrcObjName:  str[3],
				SrcObjKey:   str[4],
			}
			evt = append(evt, obj)
		}
	case events.EventStats:
		var statObj events.EventStats
		evt, _ = dbHdl.GetAllEventObjFromDb(statObj)
	default:
	}
	return evt, err
}

func constructQueryPattern(evtObj events.Event) string {
	pattern := "Events#"
	if evtObj.OwnerName == "" {
		pattern = pattern + "*#"
	} else {
		pattern = pattern + strings.ToUpper(evtObj.OwnerName) + "#"
	}
	if evtObj.EventName == "" {
		pattern = pattern + "*#"
	} else {
		pattern = pattern + evtObj.EventName + "#"
	}
	if evtObj.SrcObjName == "" {
		pattern = pattern + "*#"
	} else {
		pattern = pattern + evtObj.SrcObjName + "#"
	}
	if evtObj.SrcObjKey == "" {
		pattern = pattern + "*#"
	} else {
		pattern = pattern + evtObj.SrcObjKey + "#"
	}
	pattern = pattern + "*"
	return pattern
}

func constructKeySlice(keys []string) []KeyObj {
	var kObjSlice []KeyObj
	for _, key := range keys {
		str := strings.Split(key, "#")
		uTime, err := strconv.ParseInt(str[len(str)-2], 10, 64)
		if err != nil {
			fmt.Println("Unable to Parse Int64")
			continue
		}
		kObj := KeyObj{
			Key:   key,
			UTime: uTime,
		}
		kObjSlice = append(kObjSlice, kObj)
	}
	return kObjSlice
}

func (kObjSlice KeyObjSlice) Less(i, j int) bool {
	return kObjSlice[i].UTime > kObjSlice[j].UTime
}

func (kObjSlice KeyObjSlice) Swap(i, j int) {
	kObjSlice[i], kObjSlice[j] = kObjSlice[j], kObjSlice[i]
}

func (kObjSlice KeyObjSlice) Len() int {
	return len(kObjSlice)
}
