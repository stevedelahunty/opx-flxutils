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

package asicdClntDefs

import (
	"models/objects"
)

type IPv4IntfStateGetInfo struct {
	StartIdx          int32
	EndIdx            int32
	Count             int32
	More              bool
	IPv4IntfStateList []*objects.IPv4IntfState
}

type PortGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	PortList []*objects.Port
}

type PortStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	PortStateList []*objects.PortState
}

/* Internal Vlan Object */
type Vlan struct {
	VlanId           int32
	IfIndexList      []int32
	UntagIfIndexList []int32
	VlanName         string
}

type VlanGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	VlanList []*Vlan
}

type VlanStateGetInfo struct {
	StartIdx      int32
	EndIdx        int32
	Count         int32
	More          bool
	VlanStateList []*objects.VlanState
}

type Lag struct {
	LagIfIndex  int32
	HashType    int32
	IfIndexList []int32
	LagName     string
}

type LagGetInfo struct {
	StartIdx int32
	EndIdx   int32
	Count    int32
	More     bool
	LagList  []*Lag
}
