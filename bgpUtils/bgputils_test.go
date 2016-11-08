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

package bgpUtils

import (
	"fmt"
	"testing"
)

func TestEncodeExtCommunity(t *testing.T) {
	fmt.Println("****TestEncodeExtCommunity()****")
	var extComm ExtCommunity
	var err error
	var val string
	extComm = ExtCommunity{"Route-Target", "200:10"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "65535:4294967295"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "200.1.1.22:10"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "70000:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "65535.210:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)
	extComm = ExtCommunity{"Route-Target", "75535.210:300"}
	val, err = EncodeExtCommunity(extComm)
	fmt.Println("extComm bytes for inp:", extComm, " is - ", val, " error : ", err)

}
func TestGetCommunityValue(t *testing.T) {
	fmt.Println("****TestGetCommunityValue()****")
	comm := ""
	comm = "200:10"
	commVal, _ := GetCommunityValue(comm)
	fmt.Println("Comm bytes for inp:", comm, " is - ", commVal)
	comm = "65535:4294967295"
	commVal, _ = GetCommunityValue(comm)
	fmt.Println("Comm bytes for inp:", comm, " is - ", commVal)

}
