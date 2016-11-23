//
//Copyright [2016] [SnapRoute Inc]
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//	 Unless required by applicable law or agreed to in writing, software
//	 distributed under the License is distributed on an "AS IS" BASIS,
//	 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	 See the License for the specific language governing permissions and
//	 limitations under the License.
//
// _______  __       __________   ___      _______.____    __    ____  __  .___________.  ______  __    __
// |   ____||  |     |   ____\  \ /  /     /       |\   \  /  \  /   / |  | |           | /      ||  |  |  |
// |  |__   |  |     |  |__   \  V  /     |   (----` \   \/    \/   /  |  | `---|  |----`|  ,----'|  |__|  |
// |   __|  |  |     |   __|   >   <       \   \      \            /   |  |     |  |     |  |     |   __   |
// |  |     |  `----.|  |____ /  .  \  .----)   |      \    /\    /    |  |     |  |     |  `----.|  |  |  |
// |__|     |_______||_______/__/ \__\ |_______/        \__/  \__/     |__|     |__|      \______||__|  |__|
//

package commonDefs

import (
	"git.apache.org/thrift.git/lib/go/thrift"
)

//L2 types
const (
	IfTypePort = iota
	IfTypeLag
	IfTypeVlan
	IfTypeP2P
	IfTypeBcast
	IfTypeLoopback
	IfTypeSecondary
	IfTypeVirtual
	IfTypeVtep
	IfTypeNull
)

var BGPWellKnownCommunitiesMap = map[string]uint32{
	"NO_EXPORT":           0xFFFFFF01,
	"NO_ADVERTISE":        0xFFFFFF02,
	"NO_EXPORT_SUBCONFED": 0xFFFFFF03,
}

func GetIfTypeName(ifType int) string {
	switch ifType {
	case IfTypePort:
		return "Port"
	case IfTypeLag:
		return "Lag"
	case IfTypeVlan:
		return "Vlan"
	case IfTypeVtep:
		return "Vtep"
	default:
		return "Unknown"
	}
}

const (
	MAX_JSON_LENGTH = 4096
)

const (
	// system wide common notifications
	_ = iota
	NOTIFY_IPV6_NEIGHBOR_CREATE
	NOTIFY_IPV6_NEIGHBOR_DELETE
)

// commond object for Neighbor Notifications
type Ipv6NeighborNotification struct {
	IpAddr  string
	IfIndex int32
}

type NdpNotification struct {
	MsgType uint8
	Msg     []byte
}

type ClientJson struct {
	Name string `json:Name`
	Port int    `json:Port`
}

type ClientBase struct {
	Address            string
	Transport          thrift.TTransport
	PtrProtocolFactory *thrift.TBinaryProtocolFactory
}

const (
	FLEXSWITCH_PLUGIN       = "flexswitch"
	MOCK_PLUGIN             = "mock"
	SUB_INTF_VIRTUAL_TYPE   = "virtual"
	SUB_INTF_SECONDARY_TYPE = "secondary"
)
