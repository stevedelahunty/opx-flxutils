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

package asicdClntIntfs

import (
	"asicd/fsAsicdClnt"
	"errors"
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

type AsicdClntPluginName string

const (
	FS_ASICD_CLNT AsicdClntPluginName = "FSAsicdClnt"
)

type AsicdClntIntf interface {
	CreateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (int32, error)
	UpdateIPv4Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (int32, error)
	DeleteIPv4Neighbor(ipAddr string) (int32, error)
	GetBulkIPv4IntfState(curMark, count int) (*asicdClntDefs.IPv4IntfStateGetInfo, error)

	CreateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (int32, error)
	UpdateIPv6Neighbor(ipAddr string, macAddr string, vlanId int32, ifIdx int32) (int32, error)
	DeleteIPv6Neighbor(ipAddr string) (int32, error)

	GetBulkPort(curMark, count int) (*asicdClntDefs.PortGetInfo, error)
	GetBulkPortState(curMark, count int) (*asicdClntDefs.PortStateGetInfo, error)
	GetBulkVlan(curMark, count int) (*asicdClntDefs.VlanGetInfo, error)
	GetBulkVlanState(curMark, count int) (*asicdClntDefs.VlanStateGetInfo, error)

	CreateLag(ifname string, hashType int32, ports string) (int32, error)
	DeleteLag(ifIndex int32) (int32, error)
	UpdateLag(ifIndex, hashType int32, ports string) (int32, error)
	UpdateLagCfgIntfList(ifName string, ifIndexList []int32) (bool, error)
	GetBulkLag(curMark, count int) (*asicdClntDefs.LagGetInfo, error)

	GetAllSubIPv4IntfState() ([]*asicdClntDefs.SubIPv4IntfState, error)
	GetSysRsvdVlan() int
	GetIntfIdFromIfIndex(ifIndex int32) int
	GetIntfTypeFromIfIndex(ifIndex int32) int
	GetIfIndexFromIntfIdAndIntfType(ifId int, ifType int) int32
	GetMinSysPort() int
	IsLinuxOnlyPlugin() (bool, error)
}

func NewAsicdClntInit(clntPluginName AsicdClntPluginName, paramsFile string, asicdHdl asicdClntDefs.AsicdClientStruct) (AsicdClntIntf, error) {
	switch clntPluginName {
	case FS_ASICD_CLNT:
		return fsAsicdClnt.NewAsicdClntInit(paramsFile, asicdHdl)
	default:
		return nil, errors.New("Invalid Asicd Client Plugin Name")
	}
}
