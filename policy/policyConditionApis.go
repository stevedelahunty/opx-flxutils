// PolicyConditionApis.go
package policy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"utils/netUtils"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)

type PolicyPrefix struct {
	IpPrefix        string //CIDR eg: 1.1.1.2/24
	MasklengthRange string //exact or a specific range 21..24
}
type PolicyDstIpMatchPrefixSetCondition struct {
	PrefixSet string
	Prefix    PolicyPrefix
}

type MatchPrefixConditionInfo struct {
	UsePrefixSet bool
	PrefixSet    string
	DstIpMatch   bool
	SrcIpMatch   bool
	Prefix       PolicyPrefix
	IpPrefix     patriciaDB.Prefix //network prefix
	LowRange     int
	HighRange    int
}
type PolicyConditionConfig struct {
	Name                          string
	ConditionType                 string
	MatchProtocolConditionInfo    string
	MatchDstIpPrefixConditionInfo PolicyDstIpMatchPrefixSetCondition
	//MatchNeighborConditionInfo   PolicyMatchNeighborSetCondition
	//MatchTagConditionInfo   PolicyMatchTagSetCondition
}

type PolicyCondition struct {
	Name                 string
	ConditionType        int
	ConditionInfo        interface{}
	PolicyStmtList       []string
	ConditionGetBulkInfo string
	LocalDBSliceIdx      int
}

func (db *PolicyEngineDB) CreatePolicyDstIpMatchPrefixSetCondition(inCfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyDstIpMatchPrefixSetCondition"))
	cfg := inCfg.MatchDstIpPrefixConditionInfo
	var conditionInfo MatchPrefixConditionInfo
	conditionInfo.HighRange = -1
	conditionInfo.LowRange = -1
	var conditionGetBulkInfo string
	if len(cfg.PrefixSet) == 0 && len(cfg.Prefix.IpPrefix) == 0 {
		db.Logger.Err(fmt.Sprintln("Empty prefix set/nil prefix"))
		err = errors.New("Empty prefix set/nil prefix")
		return false, err
	}
	if len(cfg.PrefixSet) != 0 && len(cfg.Prefix.IpPrefix) != 0 {
		db.Logger.Err(fmt.Sprintln("Cannot provide both prefix set and individual prefix"))
		err = errors.New("Cannot provide both prefix set and individual prefix")
		return false, err
	}
	if len(cfg.Prefix.IpPrefix) != 0 {
		conditionGetBulkInfo = "match destination Prefix " + cfg.Prefix.IpPrefix + "MasklengthRange " + cfg.Prefix.MasklengthRange
		conditionInfo.UsePrefixSet = false
		conditionInfo.Prefix.IpPrefix = cfg.Prefix.IpPrefix
		conditionInfo.Prefix.MasklengthRange = cfg.Prefix.MasklengthRange
		conditionInfo.IpPrefix, err = netUtils.GetNetworkPrefixFromCIDR(conditionInfo.Prefix.IpPrefix)
		if err != nil {
			db.Logger.Err(fmt.Sprintln("ipPrefix invalid "))
			return
		}
		if cfg.Prefix.MasklengthRange == "exact" {
		} else {
			maskList := strings.Split(conditionInfo.Prefix.MasklengthRange, "-")
			if len(maskList) != 2 {
				db.Logger.Err(fmt.Sprintln("Invalid masklength range"))
				return
			}
			conditionInfo.LowRange, err = strconv.Atoi(maskList[0])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("lowRange mask not valid"))
				return
			}
			conditionInfo.HighRange, err = strconv.Atoi(maskList[1])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("highRange mask not valid"))
				return
			}
			db.Logger.Info(fmt.Sprintln("lowRange = ", conditionInfo.LowRange, " highrange = ", conditionInfo.HighRange))
		}
	} else if len(cfg.PrefixSet) != 0 {
		conditionInfo.UsePrefixSet = true
		conditionInfo.PrefixSet = cfg.PrefixSet
		conditionGetBulkInfo = "match destination Prefix " + cfg.PrefixSet
	}
	conditionInfo.DstIpMatch = true
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", inCfg.Name))
		newPolicyCondition := PolicyCondition{Name: inCfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch, ConditionInfo: conditionInfo, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = conditionGetBulkInfo
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(inCfg.Name), newPolicyCondition); ok != true {
			db.Logger.Err(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(inCfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}

func (db *PolicyEngineDB) CreatePolicyMatchProtocolCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyMatchProtocolCondition"))

	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if policyCondition == nil {
		db.Logger.Info(fmt.Sprintln("Defining a new policy condition with name ", cfg.Name, " to match on protocol ", cfg.MatchProtocolConditionInfo))
		matchProto := cfg.MatchProtocolConditionInfo
		newPolicyCondition := PolicyCondition{Name: cfg.Name, ConditionType: policyCommonDefs.PolicyConditionTypeProtocolMatch, ConditionInfo: matchProto, LocalDBSliceIdx: (len(*db.LocalPolicyConditionsDB))}
		newPolicyCondition.ConditionGetBulkInfo = "match Protocol " + matchProto
		if ok := db.PolicyConditionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyCondition); ok != true {
			db.Logger.Info(fmt.Sprintln(" return value not ok"))
			err = errors.New("Error creating condition in the DB")
			return false, err
		}
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), add)
	} else {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return false, err
	}
	return true, err
}
func (db *PolicyEngineDB) ValidateConditionConfigCreate(inCfg PolicyConditionConfig) (err error) {
	db.Logger.Info(fmt.Sprintln("ValidateConditionConfigCreate"))
	policyCondition := db.PolicyConditionsDB.Get(patriciaDB.Prefix(inCfg.Name))
	if policyCondition != nil {
		db.Logger.Err(fmt.Sprintln("Duplicate Condition name"))
		err = errors.New("Duplicate policy condition definition")
		return err
	}
	switch inCfg.ConditionType {
	case "MatchDstIpPrefix":
	    cfg := inCfg.MatchDstIpPrefixConditionInfo
	    if len(cfg.PrefixSet) == 0 && len(cfg.Prefix.IpPrefix) == 0 {
		    db.Logger.Err(fmt.Sprintln("Empty prefix set/nil prefix"))
		    err = errors.New("Empty prefix set/nil prefix")
		    return err
	    }
	    if len(cfg.PrefixSet) != 0 && len(cfg.Prefix.IpPrefix) != 0 {
		    db.Logger.Err(fmt.Sprintln("Cannot provide both prefix set and individual prefix"))
		    err = errors.New("Cannot provide both prefix set and individual prefix")
		    return err
	    }
	    if len(cfg.Prefix.IpPrefix) != 0 {
		_, err = netUtils.GetNetworkPrefixFromCIDR(cfg.Prefix.IpPrefix)
		if err != nil {
			db.Logger.Err(fmt.Sprintln("ipPrefix invalid "))
			return errors.New("ipPrefix invalid") 
		}
		if cfg.Prefix.MasklengthRange == "exact" {
		} else {
			maskList := strings.Split(cfg.Prefix.MasklengthRange, "-")
			if len(maskList) != 2 {
				db.Logger.Err(fmt.Sprintln("Invalid masklength range"))
				return errors.New("Invalid masklength range")
			}
			_, err = strconv.Atoi(maskList[0])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("lowRange mask not valid"))
				return errors.New("lowRange mask not valid")
			}
			_, err = strconv.Atoi(maskList[1])
			if err != nil {
				db.Logger.Err(fmt.Sprintln("highRange mask not valid"))
				return errors.New("highRange mask not valid")
			}
		}
	}		
	default:
		db.Logger.Err(fmt.Sprintln("Unknown condition type ", inCfg.ConditionType))
		err = errors.New("Unknown condition type")
		return err
	}
	return err
}
func (db *PolicyEngineDB) CreatePolicyCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("CreatePolicyCondition"))
	switch cfg.ConditionType {
	case "MatchDstIpPrefix":
		val, err = db.CreatePolicyDstIpMatchPrefixSetCondition(cfg)
		break
	case "MatchProtocol":
		val, err = db.CreatePolicyMatchProtocolCondition(cfg)
		break
	default:
		db.Logger.Err(fmt.Sprintln("Unknown condition type ", cfg.ConditionType))
		err = errors.New("Unknown condition type")
		return false, err
	}
	return val, err
}
func (db *PolicyEngineDB) DeletePolicyCondition(cfg PolicyConditionConfig) (val bool, err error) {
	db.Logger.Info(fmt.Sprintln("DeletePolicyCondition"))
	conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if conditionItem == nil {
		db.Logger.Err(fmt.Sprintln("Condition ", cfg.Name, "not found in the DB"))
		err = errors.New("Condition not found")
		return false, err
	}
	condition := conditionItem.(PolicyCondition)
	if len(condition.PolicyStmtList) != 0 {
		db.Logger.Err(fmt.Sprintln("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition"))
		err = errors.New("This condition is currently being used by a policy statement. Try deleting the stmt before deleting the condition")
		return false, err
	}
	deleted := db.PolicyConditionsDB.Delete(patriciaDB.Prefix(cfg.Name))
	if deleted {
		db.Logger.Info(fmt.Sprintln("Found and deleted condition ", cfg.Name))
		db.LocalPolicyConditionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name), del)
	}
	return true, err
}
