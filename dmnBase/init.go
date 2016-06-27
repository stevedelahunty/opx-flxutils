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

package dmnBase

import (
	_ "utils/asicdClient/ovs"
	"utils/commonDefs"
)

type SwitchIntf interface {
	ConnectToServers() error
	CreateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	UpdateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (rv int32, err error)
	DeleteIPv4Neighbor(ipAddr string) (rv int32, err error)

	//GetBulkIPv4IntfState(curMark, count int) (*commonDefs.IPv4IntfStateGetInfo, error)
	GetAllIPv4IntfState() ([]*commonDefs.IPv4IntfState, error)
	GetAllPort(curMark, count int) (*commonDefs.PortGetInfo, error)
	GetAllPortState(curMark, count int) (*commonDefs.PortStateGetInfo, error)
	GetAllVlan(curMark, count int) (*commonDefs.VlanGetInfo, error)
	GetAllVlanState(curMark, count int) (*commonDefs.VlanStateGetInfo, error)
	DetermineRouterId() string
	Log(string, string)
}

const (
	INFO  = "info"
	ERR   = "error"
	DEBUG = "debug"
	ALERT = "alert"
)

// Need to return an interface object
func InitPlugin(dmnName, logPrefix, plugin string) SwitchIntf {
	switch plugin {
	case "Flexswitch":
		switchDmn := &FSDaemon{}
		switchDmn.FSBaseDmn = NewBaseDmn(dmnName, logPrefix)
		switchDmn.FSBaseDmn.Init()
		switchDmn.NewServer()
		return switchDmn
	case "OvsDB":
		// @TODO: for future
		switchDmn := &FSDaemon{} //&ovs.OvsAsicdClientMgr{}
		return switchDmn
	}
	return nil
}
