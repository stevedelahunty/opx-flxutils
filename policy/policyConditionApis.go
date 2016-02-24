// PolicyConditionApis.go
package policy

import (
	"errors"
	"utils/policy/policyCommonDefs"
	"utils/patriciaDB"
	"fmt"
)
type PolicyPrefix struct{
	IpPrefix string //CIDR eg: 1.1.1.2/24
	MasklengthRange string //exact or a specific range 21..24
}
type PolicyDstIpMatchPrefixSetCondition struct{
	PrefixSet string
	Prefix PolicyPrefix
}

type MatchPrefixConditionInfo struct {
	UsePrefixSet bool
	PrefixSet string
	DstIpMatch     bool
	SrcIpMatch     bool
	MatchPrefix PolicyPrefix
}
type PolicyConditionConfig struct {
	Name string
	ConditionType string
	MatchProtocolConditionInfo string          
    MatchDstIpPrefixConditionInfo PolicyDstIpMatchPrefixSetCondition       
    //MatchNeighborConditionInfo   PolicyMatchNeighborSetCondition        
	//MatchTagConditionInfo   PolicyMatchTagSetCondition             
}

type PolicyCondition struct {
	Name          string
	ConditionType int
	ConditionInfo interface {}
	PolicyStmtList    [] string
	ConditionGetBulkInfo string
	LocalDBSliceIdx int
}
func updateLocalConditionsDB(prefix patriciaDB.Prefix, localPolicyConditionsDB []LocalDB) {
	localDBRecord := LocalDB{Prefix:prefix, IsValid:true}
	if(localPolicyConditionsDB == nil) {
		localPolicyConditionsDB = make([]LocalDB, 0)
	} 
	localPolicyConditionsDB = append(localPolicyConditionsDB, localDBRecord)

}
func CreatePolicyDstIpMatchPrefixSetCondition(PolicyConditionsDB *patriciaDB.Trie, localPolicyConditionsDB []LocalDB, inCfg PolicyConditionConfig) (val bool, err error) {
	fmt.Println("CreatePolicyDstIpMatchPrefixSetCondition")
	cfg := inCfg.MatchDstIpPrefixConditionInfo
	var conditionInfo MatchPrefixConditionInfo
	var conditionGetBulkInfo string
    if len(cfg.PrefixSet) == 0 && len(cfg.Prefix.IpPrefix) == 0 {
		fmt.Println("Empty prefix set")
		err = errors.New("Empty prefix set")
		return val, err
	}
    if len(cfg.PrefixSet) != 0 && len(cfg.Prefix.IpPrefix) != 0 {
		fmt.Println("Cannot provide both prefix set and individual prefix")
		err = errors.New("Cannot provide both prefix set and individual prefix")
		return val, err
	}
    if len(cfg.Prefix.IpPrefix) != 0 {
	   conditionInfo.UsePrefixSet = false
       conditionInfo.MatchPrefix.IpPrefix = cfg.Prefix.IpPrefix
	   conditionInfo.MatchPrefix.MasklengthRange = cfg.Prefix.MasklengthRange
	   conditionGetBulkInfo = "match destination Prefix " + cfg.Prefix.IpPrefix + "MasklengthRange " + cfg.Prefix.MasklengthRange
	} else if len(cfg.PrefixSet) != 0 {
		conditionInfo.UsePrefixSet = true
		conditionInfo.PrefixSet = cfg.PrefixSet
	    conditionGetBulkInfo = "match destination Prefix " + cfg.PrefixSet
	}
	conditionInfo.DstIpMatch = true
	policyCondition := PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if(policyCondition == nil) {
	   fmt.Println("Defining a new policy condition with name ", inCfg.Name)
	   newPolicyCondition := PolicyCondition{Name:inCfg.Name,ConditionType:policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch,ConditionInfo:conditionInfo ,LocalDBSliceIdx:(len(localPolicyConditionsDB))}
       newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo 
	   if ok := PolicyConditionsDB.Insert(patriciaDB.Prefix(inCfg.Name), newPolicyCondition); ok != true {
	   fmt.Println(" return value not ok")
	   return val, err
	}
	updateLocalConditionsDB(patriciaDB.Prefix(inCfg.Name), localPolicyConditionsDB)
    } else {
		fmt.Println("Duplicate Condition name")
		err = errors.New("Duplicate policy condition definition")
		return val, err
	}	
	return val, err
}

func CreatePolicyMatchProtocolCondition(PolicyConditionsDB *patriciaDB.Trie, localPolicyConditionsDB[] LocalDB, cfg PolicyConditionConfig) (val bool, err error) {
	fmt.Println("CreatePolicyMatchProtocolCondition")

	policyCondition := PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyCondition == nil) {
	   fmt.Println("Defining a new policy condition with name ", cfg.Name)
	   matchProto := cfg.MatchProtocolConditionInfo
	   newPolicyCondition := PolicyCondition{Name:cfg.Name,ConditionType:policyCommonDefs.PolicyConditionTypeProtocolMatch,ConditionInfo:matchProto ,LocalDBSliceIdx:(len(localPolicyConditionsDB))}
       newPolicyCondition.ConditionGetBulkInfo = "match Protocol " + matchProto
		if ok := PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			fmt.Println(" return value not ok")
			return val, err
		}
	    updateLocalConditionsDB(patriciaDB.Prefix(cfg.Name), localPolicyConditionsDB)
	} else {
		fmt.Println("Duplicate Condition name")
		err = errors.New("Duplicate policy condition definition")
		return val, err
	}
	return val, err
}
func CreatePolicyCondition(ConditionDB *patriciaDB.Trie, localConditionsDB[]LocalDB, cfg PolicyConditionConfig) (err error) {
	fmt.Println("CreatePolicyCondition")
	switch cfg.ConditionType {
		case "MatchDstIpPrefix":
		   CreatePolicyDstIpMatchPrefixSetCondition(ConditionDB, localConditionsDB, cfg)
		   break
		case "MatchProtocol":
		   CreatePolicyMatchProtocolCondition(ConditionDB, localConditionsDB, cfg)
		   break
		default:
		   fmt.Println("Unknown condition type ", cfg.ConditionType)
		   err = errors.New("Unknown condition type")
	}
	return err
}

