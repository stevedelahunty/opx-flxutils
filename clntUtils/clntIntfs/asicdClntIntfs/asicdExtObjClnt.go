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
	"models/objects"
	"utils/clntUtils/clntDefs"
	"utils/clntUtils/clntDefs/asicdClntDefs"
)

type AsicdExtObjClntIntf interface {
	GetBulkIPv4IntfState(fromIndex int, count int) (*asicdClntDefs.IPv4IntfStateGetInfo, error)
	GetIPv4IntfState(IntfRef string) (*objects.IPv4IntfState, error)

	GetBulkPort(fromIndex int, count int) (*asicdClntDefs.PortGetInfo, error)
	GetPort(IntfRef string) (*objects.Port, error)
	GetBulkPortState(fromIndex int, count int) (*asicdClntDefs.PortStateGetInfo, error)
	GetPortState(IntfRef string) (*objects.PortState, error)
	GetBulkVlanState(fromIndex int, count int) (*asicdClntDefs.VlanStateGetInfo, error)
	GetVlanState(VlanId int32) (*objects.VlanState, error)
	CreateVlan(cfg *objects.Vlan) (bool, error)
	UpdateVlan(origCfg, newCfg *objects.Vlan, attrset []bool, op []*clntDefs.PatchOpInfo) (bool, error)
	DeleteVlan(cfg *objects.Vlan) (bool, error)
}
