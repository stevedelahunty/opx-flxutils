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

package arpdClntIntfs

import (
	"errors"
	"l3/arp/fsArpdClnt"
	"utils/clntUtils/clntDefs/arpdClntDefs"
)

type ArpdClntPluginName string

const (
	FS_ARPD_CLNT ArpdClntPluginName = "FSArpdClnt"
)

type ArpdClntIntf interface {
	GetBulkArpEntryState(fromIndex int, count int) (*arpdClntDefs.ArpEntryStateGetInfo, error)
	GetArpEntryState(ipAddr string) (*arpdClntDefs.ArpEntryState, error)
	GetBulkArpLinuxEntryState(fromIndex int, count int) (*arpdClntDefs.ArpLinuxEntryStateGetInfo, error)
	GetArpLinuxEntryState(ipAddr string) (*arpdClntDefs.ArpLinuxEntryState, error)
	ExecuteActionArpDeleteByIfName(cfg *arpdClntDefs.ArpDeleteByIfName) (bool, error)
	ExecuteActionArpDeleteByIPv4Addr(cfg *arpdClntDefs.ArpDeleteByIPv4Addr) (bool, error)
	ExecuteActionArpRefreshByIfName(cfg *arpdClntDefs.ArpRefreshByIfName) (bool, error)
	ExecuteActionArpRefreshByIPv4Addr(cfg *arpdClntDefs.ArpRefreshByIPv4Addr) (bool, error)
	CreateArpGlobal(cfg *arpdClntDefs.ArpGlobal) (bool, error)
	UpdateArpGlobal(origCfg *arpdClntDefs.ArpGlobal, newCfg *arpdClntDefs.ArpGlobal, attrset []bool, op []*arpdClntDefs.PatchOpInfo) (bool, error)
	DeleteArpGlobal(cfg *arpdClntDefs.ArpGlobal) (bool, error)

	ResolveArpIPv4(destNetIp string, ifIdx int32) (err error)
	DeleteResolveArpIPv4(NbrIP string) (err error)
	DeleteArpEntry(ipAddr string) (err error)
	SendGarp(ifName string, macAddr string, ipAddr string) (err error)
}

func NewArpdClntInit(clntPluginName ArpdClntPluginName, paramsFile string, arpdHdl arpdClntDefs.ArpdClientStruct) (ArpdClntIntf, error) {
	switch clntPluginName {
	case FS_ARPD_CLNT:
		return fsArpdClnt.NewArpdClntInit(paramsFile, arpdHdl)
	default:
		return nil, errors.New("Invalid Arpd Client Plugin Name")
	}
}
