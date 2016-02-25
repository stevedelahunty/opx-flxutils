// policyApis.go
package policy

import (
	"errors"
	"utils/patriciaDB"
	"utils/netUtils"
	"utils/policy/policyCommonDefs"
	"strconv"
	"strings"
	"fmt"
	"reflect"
)

type PolicyStmt struct {				//policy engine uses this
	Name               string
	Precedence         int
	MatchConditions    string
	Conditions         []string
	Actions            []string
	LocalDBSliceIdx        int8  
}
type PolicyStmtConfig struct{
	Name string
	AdminState string 
	MatchConditions string
	Conditions []string
	Actions []string
}

type Policy struct {
	Name              string
	Precedence        int
	MatchType         string
	PolicyStmtPrecedenceMap map[int]string
	LocalDBSliceIdx        int8  
	ImportPolicy       bool
	ExportPolicy       bool  
	GlobalPolicy       bool
}

type PolicyDefinitionStmtPrecedence  struct {
	Precedence int
	Statement string
}
type PolicyDefinitionConfig struct{
	Name string
	Precedence int
	MatchType string
	PolicyDefinitionStatements []PolicyDefinitionStmtPrecedence
 	Export bool
	Import bool
	Global bool
}

type PrefixPolicyListInfo struct {
	ipPrefix  patriciaDB.Prefix
	policyName string
	lowRange   int
	highRange  int
}

var PolicyDB = patriciaDB.NewTrie()
var PolicyStmtDB = patriciaDB.NewTrie()
var PolicyStmtPolicyMapDB = make(map[string] []string) //policies using this statement
var PrefixPolicyListDB = patriciaDB.NewTrie()
var ProtocolPolicyListDB = make(map[string][]string)//policystmt names assoociated with every protocol type
var ImportPolicyPrecedenceMap = make(map[int] string)
var ExportPolicyPrecedenceMap = make(map[int] string)
var LocalPolicyStmtDB []LocalDB
var LocalPolicyDB []LocalDB

func validMatchConditions(matchConditionStr string) (valid bool) {
    fmt.Println("validMatchConditions for string ", matchConditionStr)
	if matchConditionStr == "any" || matchConditionStr == "all"{
		fmt.Println("valid")
		valid = true
	}
	return valid
}
func updateProtocolPolicyTable(protoType string, name string, op int) {
	fmt.Printf("updateProtocolPolicyTable for protocol %d policy name %s op %d\n", protoType, name, op)
    var i int
    policyList := ProtocolPolicyListDB[protoType]
	if(policyList == nil) {
		if (op == del) {
			fmt.Println("Cannot find the policy map for this protocol, so cannot delete")
			return
		}
		policyList = make([]string, 0)
	}
    if op == add {
	   policyList = append(policyList, name)
	}
	found :=false
	if op == del {
		for i =0; i< len(policyList);i++ {
			if policyList[i] == name {
				fmt.Println("Found the policy in the protocol policy table, deleting it")
				found = true
				break
			}
		}
		if found {
		   policyList = append(policyList[:i], policyList[i+1:]...)
		}
	}
	ProtocolPolicyListDB[protoType] = policyList
}
func updatePrefixPolicyTableWithPrefix(ipAddr string, name string, op int, lowRange int, highRange int){
	fmt.Println("updatePrefixPolicyTableWithPrefix ", ipAddr)
	var i int
       ipPrefix, err := netUtils.GetNetworkPrefixFromCIDR(ipAddr)
	   if err != nil {
		fmt.Println("ipPrefix invalid ")
		return 
	   }
	var policyList []PrefixPolicyListInfo
	var prefixPolicyListInfo PrefixPolicyListInfo
	policyListItem:= PrefixPolicyListDB.Get(ipPrefix)
	if policyListItem != nil && reflect.TypeOf(policyListItem).Kind() != reflect.Slice {
		fmt.Println("Incorrect data type for this prefix ")
		return
	}
	if(policyListItem == nil) {
		if (op == del) {
			fmt.Println("Cannot find the policy map for this prefix, so cannot delete")
			return
		}
		policyList = make([]PrefixPolicyListInfo, 0)
	} else {
	   policyListSlice := reflect.ValueOf(policyListItem)
	   policyList = make([]PrefixPolicyListInfo,0)
	   for i = 0;i<policyListSlice.Len();i++ {
	      policyList = append(policyList, policyListSlice.Index(i).Interface().(PrefixPolicyListInfo))	
	   }
	}
    if op == add {
	   prefixPolicyListInfo.ipPrefix = ipPrefix
	   prefixPolicyListInfo.policyName = name
	   prefixPolicyListInfo.lowRange = lowRange
	   prefixPolicyListInfo.highRange = highRange
	   policyList = append(policyList, prefixPolicyListInfo)
	}
	found :=false
	if op == del {
		for i =0; i< len(policyList);i++ {
			if policyList[i].policyName == name {
				fmt.Println("Found the policy in the prefix policy table, deleting it")
				break
			}
		}
		if found {
		   policyList = append(policyList[:i], policyList[i+1:]...)
		}
	}
	PrefixPolicyListDB.Set(ipPrefix, policyList)
}
func updatePrefixPolicyTableWithMaskRange(ipAddr string, masklength string, name string, op int){
	fmt.Println("updatePrefixPolicyTableWithMaskRange")
	    maskList := strings.Split(masklength,"..")
		if len(maskList) !=2 {
			fmt.Println("Invalid masklength range")
			return 
		}
        lowRange,err := strconv.Atoi(maskList[0])
		if err != nil {
			fmt.Println("maskList[0] not valid")
			return
		}
		highRange,err := strconv.Atoi(maskList[1])
		if err != nil {
			fmt.Println("maskList[1] not valid")
			return
		}
		fmt.Println("lowRange = ", lowRange, " highrange = ", highRange)
		updatePrefixPolicyTableWithPrefix(ipAddr, name, op,lowRange,highRange)
/*		for idx := lowRange;idx<highRange;idx ++ {
			ipMask:= net.CIDRMask(idx, 32)
			ipMaskStr := net.IP(ipMask).String()
			fmt.Println("idx ", idx, "ipMaskStr = ", ipMaskStr)
			ipPrefix, err := getNetowrkPrefixFromStrings(ipAddrStr, ipMaskStr)
			if err != nil {
				fmt.Println("Invalid prefix")
				return 
			}
			updatePrefixPolicyTableWithPrefix(ipPrefix, name, op,lowRange,highRange)
		}*/
}
func updatePrefixPolicyTableWithPrefixSet(prefixSet string, name string, op int) {
	fmt.Println("updatePrefixPolicyTableWithPrefixSet")
}
func updatePrefixPolicyTable(conditionInfo interface{}, name string, op int) {
    condition := conditionInfo.(MatchPrefixConditionInfo)
	fmt.Printf("updatePrefixPolicyTable for prefixSet %s prefix %s policy name %s op %d\n", condition.PrefixSet, condition.Prefix, name, op)
    if condition.UsePrefixSet {
		fmt.Println("Need to look up Prefix set to get the prefixes")
		updatePrefixPolicyTableWithPrefixSet(condition.PrefixSet, name, op)
	} else {
	   if condition.Prefix.MasklengthRange == "exact" {
       /*ipPrefix, err := getNetworkPrefixFromCIDR(condition.prefix.IpPrefix)
	   if err != nil {
		fmt.Println("ipPrefix invalid ")
		return 
	   }*/
	   updatePrefixPolicyTableWithPrefix(condition.Prefix.IpPrefix, name, op,-1,-1)
	 } else {
		fmt.Println("Masklength= ", condition.Prefix.MasklengthRange)
		updatePrefixPolicyTableWithMaskRange(condition.Prefix.IpPrefix, condition.Prefix.MasklengthRange, name, op)
	 }
   }
}
func updateStatements(policy  string, stmt string, op int) (err error){
   fmt.Println("updateStatements stmt ", stmt, " with policy ", policy)
   var i int
    policyList := PolicyStmtPolicyMapDB[stmt]
	if(policyList == nil) {
		if (op == del) {
			fmt.Println("Cannot find the policy map for this stmt, so cannot delete")
            err = errors.New("Cannot find the policy map for this stmt, so cannot delete")
			return err
		}
		policyList = make([]string, 0)
	}
    if op == add {
	   policyList = append(policyList, policy)
	}
	found :=false
	if op == del {
		for i =0; i< len(policyList);i++ {
			if policyList[i] == policy {
				fmt.Println("Found the policy in the policy stmt table, deleting it")
                 found = true
				break
			}
		}
		if found {
		   policyList = append(policyList[:i], policyList[i+1:]...)
		}
	}
	PolicyStmtPolicyMapDB[stmt] = policyList
	return err
}
func updateConditions(policyStmt PolicyStmt, conditionName string, op int) (err error){
	fmt.Println("updateConditions for condition ", conditionName)
	conditionItem := PolicyConditionsDB.Get(patriciaDB.Prefix(conditionName))
	if(conditionItem != nil) {
		condition := conditionItem.(PolicyCondition)
		switch condition.ConditionType {
			case policyCommonDefs.PolicyConditionTypeProtocolMatch:
			   fmt.Println("PolicyConditionTypeProtocolMatch")
			   updateProtocolPolicyTable(condition.ConditionInfo.(string), policyStmt.Name, op)
			   break
			case policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch:
			   fmt.Println("PolicyConditionTypeDstIpPrefixMatch")
			   updatePrefixPolicyTable(condition.ConditionInfo, policyStmt.Name, op)
			   break
		}
		if condition.PolicyStmtList == nil {
			condition.PolicyStmtList = make([]string,0)
		}
        condition.PolicyStmtList = append(condition.PolicyStmtList, policyStmt.Name)
		fmt.Println("Adding policy ", policyStmt.Name, "to condition ", conditionName)
		PolicyConditionsDB.Set(patriciaDB.Prefix(conditionName), condition)
	} else {
		fmt.Println("Condition name ", conditionName, " not defined")
		err = errors.New("Condition name not defined")
	}
	return err
}

func updateActions(policyStmt PolicyStmt, actionName string, op int) (err error) {
	fmt.Println("updateActions for action ", actionName)
	actionItem := PolicyActionsDB.Get(patriciaDB.Prefix(actionName))
	if(actionItem != nil) {
		action := actionItem.(PolicyAction)
		if action.PolicyStmtList == nil {
			action.PolicyStmtList = make([]string,0)
		}
        action.PolicyStmtList = append(action.PolicyStmtList, policyStmt.Name)
		PolicyActionsDB.Set(patriciaDB.Prefix(actionName), action)
	} else {
		fmt.Println("action name ", actionName, " not defined")
		err = errors.New("action name not defined")
	}
	return err
}

func CreatePolicyStatement(cfg PolicyStmtConfig) (err error) {
	fmt.Println("CreatePolicyStatement")
	policyStmt := PolicyStmtDB.Get(patriciaDB.Prefix(cfg.Name))
	var i int
	if(policyStmt == nil) {
	   fmt.Println("Defining a new policy statement with name ", cfg.Name)
	   var newPolicyStmt PolicyStmt
	   newPolicyStmt.Name = cfg.Name
	   if !validMatchConditions(cfg.MatchConditions) {
	      fmt.Println("Invalid match conditions - try any/all")
		  err = errors.New("Invalid match conditions - try any/all")	
		  return  err
	   }
	   newPolicyStmt.MatchConditions = cfg.MatchConditions
	   if len(cfg.Conditions) > 0 {
	      fmt.Println("Policy Statement has %d ", len(cfg.Conditions)," number of conditions")	
		  newPolicyStmt.Conditions = make([] string, 0)
		  for i=0;i<len(cfg.Conditions);i++ {
			newPolicyStmt.Conditions = append(newPolicyStmt.Conditions, cfg.Conditions[i])
			err = updateConditions(newPolicyStmt, cfg.Conditions[i], add)
			if err != nil {
				fmt.Println("updateConditions returned err ", err)
				return err
			}
		}
	   }
	   if len(cfg.Actions) > 0 {
	      fmt.Println("Policy Statement has %d ", len(cfg.Actions)," number of actions")	
		  newPolicyStmt.Actions = make([] string, 0)
		  for i=0;i<len(cfg.Actions);i++ {
			newPolicyStmt.Actions = append(newPolicyStmt.Actions,cfg.Actions[i])
			err = updateActions(newPolicyStmt, cfg.Actions[i], add)
			if err != nil {
				fmt.Println("updateActions returned err ", err)
				return err
			}
		}
	   }
        newPolicyStmt.LocalDBSliceIdx = int8(len(LocalPolicyStmtDB))
		if ok := PolicyStmtDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyStmt); ok != true {
			fmt.Println(" return value not ok")
			return err
		}
        localDBRecord := LocalDB{Prefix:patriciaDB.Prefix(cfg.Name), IsValid:true}
		if(LocalPolicyStmtDB == nil) {
			LocalPolicyStmtDB = make([]LocalDB, 0)
		} 
	    LocalPolicyStmtDB = append(LocalPolicyStmtDB, localDBRecord)
	} else {
		fmt.Println("Duplicate Policy definition name")
		err = errors.New("Duplicate policy definition")
		return err
	}
	return err
}

func DeletePolicyStatement(cfg PolicyStmtConfig) (err error) {
	fmt.Println("DeletePolicyStatement for name ", cfg.Name)
	ok := PolicyStmtDB.Match(patriciaDB.Prefix(cfg.Name))
	if !ok {
		err = errors.New("No policy statement with this name found")
		return err
	}
	policyStmtInfoGet := PolicyStmtDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyStmtInfoGet != nil) {
       //invalidate localPolicyStmt 
	   policyStmtInfo := policyStmtInfoGet.(PolicyStmt)
	   if policyStmtInfo.LocalDBSliceIdx < int8(len(LocalPolicyStmtDB)) {
          fmt.Println("local DB slice index for this policy stmt is ", policyStmtInfo.LocalDBSliceIdx)
		  LocalPolicyStmtDB[policyStmtInfo.LocalDBSliceIdx].IsValid = false		
	   }
	  // PolicyEngineTraverseAndReverse(policyStmtInfo)
	   fmt.Println("Deleting policy statement with name ", cfg.Name)
		if ok := PolicyStmtDB.Delete(patriciaDB.Prefix(cfg.Name)); ok != true {
			fmt.Println(" return value not ok for delete PolicyDB")
			return err
		}
	   //update other tables
	   if len(policyStmtInfo.Conditions) > 0 {
	      for i:=0;i<len(policyStmtInfo.Conditions);i++ {
			updateConditions(policyStmtInfo, policyStmtInfo.Conditions[i],del)
		}	
	   }
	   if len(policyStmtInfo.Actions) > 0 {
	      for i:=0;i<len(policyStmtInfo.Actions);i++ {
			updateActions(policyStmtInfo, policyStmtInfo.Actions[i],del)
		}	
	   }
	} 
	return err
}

func CreatePolicyDefinition(cfg PolicyDefinitionConfig) (err error) {
	fmt.Println("CreatePolicyDefinition")
	if cfg.Import && ImportPolicyPrecedenceMap != nil {
	   _,ok:=ImportPolicyPrecedenceMap[int(cfg.Precedence)]
	   if ok {
		fmt.Println("There is already a import policy with this precedence.")
		err =  errors.New("There is already a import policy with this precedence.")
         return err
	   }
	} else if cfg.Export && ExportPolicyPrecedenceMap != nil {
	   _,ok:=ExportPolicyPrecedenceMap[int(cfg.Precedence)]
	   if ok {
		fmt.Println("There is already a export policy with this precedence.")
		err =  errors.New("There is already a export policy with this precedence.")
         return err
	   }
	} else if cfg.Global {
		fmt.Println("This is a global policy")
	}
	policy := PolicyDB.Get(patriciaDB.Prefix(cfg.Name))
	var i int
	if(policy == nil) {
	   fmt.Println("Defining a new policy with name ", cfg.Name)
	   var newPolicy Policy
	   newPolicy.Name = cfg.Name
	   newPolicy.Precedence = cfg.Precedence
	   newPolicy.MatchType = cfg.MatchType
       if cfg.Export == false && cfg.Import == false && cfg.Global == false {
			fmt.Println("Need to set import, export or global to true")
			return err
	   }	  
	   newPolicy.ExportPolicy = cfg.Export
	   newPolicy.ImportPolicy = cfg.Import
	   newPolicy.GlobalPolicy = cfg.Global
	   fmt.Println("Policy has %d ", len(cfg.PolicyDefinitionStatements)," number of statements")
	   newPolicy.PolicyStmtPrecedenceMap = make(map[int]string)	
	   for i=0;i<len(cfg.PolicyDefinitionStatements);i++ {
		  fmt.Println("Adding statement ", cfg.PolicyDefinitionStatements[i].Statement, " at precedence id ", cfg.PolicyDefinitionStatements[i].Precedence)
          newPolicy.PolicyStmtPrecedenceMap[int(cfg.PolicyDefinitionStatements[i].Precedence)] = cfg.PolicyDefinitionStatements[i].Statement 
		  err = updateStatements(newPolicy.Name, cfg.PolicyDefinitionStatements[i].Statement, add)
		  if err != nil {
			fmt.Println("updateStatements returned err ", err)
			return err
		  }
	   }
       for k:=range newPolicy.PolicyStmtPrecedenceMap {
		fmt.Println("key k = ", k)
	   }
       newPolicy.LocalDBSliceIdx = int8(len(LocalPolicyDB))
	   if ok := PolicyDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicy); ok != true {
			fmt.Println(" return value not ok")
			return err
		}
        localDBRecord := LocalDB{Prefix:patriciaDB.Prefix(cfg.Name), IsValid:true}
		if(LocalPolicyDB == nil) {
			LocalPolicyDB = make([]LocalDB, 0)
		} 
	    LocalPolicyDB = append(LocalPolicyDB, localDBRecord)
		if cfg.Import {
		   fmt.Println("Adding ", newPolicy.Name, " as import policy")
		   if ImportPolicyPrecedenceMap == nil {
	          ImportPolicyPrecedenceMap = make(map[int]string)	
		   }
		   ImportPolicyPrecedenceMap[int(cfg.Precedence)]=cfg.Name
		} else if cfg.Export {
		   fmt.Println("Adding ", newPolicy.Name, " as export policy")
		   if ExportPolicyPrecedenceMap == nil {
	          ExportPolicyPrecedenceMap = make(map[int]string)	
		   }
		   ExportPolicyPrecedenceMap[int(cfg.Precedence)]=cfg.Name
		}
	     PolicyEngineTraverseAndApplyPolicy(newPolicy)
	} else {
		fmt.Println("Duplicate Policy definition name")
		err = errors.New("Duplicate policy definition")
		return err
	}
	return err
}

func DeletePolicyDefinition(cfg PolicyDefinitionConfig) (err error) {
	fmt.Println("DeletePolicyDefinition for name ", cfg.Name)
	ok := PolicyDB.Match(patriciaDB.Prefix(cfg.Name))
	if !ok {
		err = errors.New("No policy with this name found")
		return err
	}
	policyInfoGet := PolicyDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyInfoGet != nil) {
       //invalidate localPolicy 
	   policyInfo := policyInfoGet.(Policy)
	   if policyInfo.LocalDBSliceIdx < int8(len(LocalPolicyDB)) {
          fmt.Println("local DB slice index for this policy is ", policyInfo.LocalDBSliceIdx)
		  LocalPolicyDB[policyInfo.LocalDBSliceIdx].IsValid = false		
	   }
	   PolicyEngineTraverseAndReversePolicy(policyInfo)
	   fmt.Println("Deleting policy with name ", cfg.Name)
		if ok := PolicyDB.Delete(patriciaDB.Prefix(cfg.Name)); ok != true {
			fmt.Println(" return value not ok for delete PolicyDB")
			return err
		}
		for _,v:=range policyInfo.PolicyStmtPrecedenceMap {
		  err = updateStatements(policyInfo.Name, v, del)
		  if err != nil {
			fmt.Println("updateStatements returned err ", err)
			return err
		  }
		}
		if policyInfo.ExportPolicy{
			if ExportPolicyPrecedenceMap != nil {
				delete(ExportPolicyPrecedenceMap,int(policyInfo.Precedence))
			}
		}
		if policyInfo.ImportPolicy{
			if ImportPolicyPrecedenceMap != nil {
				delete(ImportPolicyPrecedenceMap,int(policyInfo.Precedence))
			}
		}
	} 
	return err
}

func GetPolicyStmtDB() (db *patriciaDB.Trie, err error) { 
	if PolicyStmtDB == nil {
		fmt.Println("policyStmt nil")
		err := errors.New("policyStmt nil")
		return nil,err
	}
	return PolicyStmtDB, err
}
func GetLocalPolicyStmtDB()(db []LocalDB, err error) { 
	if LocalPolicyStmtDB == nil {
		fmt.Println("local policyStmt nil")
		err := errors.New("local policyStmt nil")
		return nil,err
	}
	return LocalPolicyStmtDB, err
}

func GetPolicyDB() (db *patriciaDB.Trie, err error) { 
	if PolicyDB == nil {
		fmt.Println("policy nil")
		err := errors.New("policy nil")
		return nil,err
	}
	return PolicyDB, err
}
func GetLocalPolicyDB()(db []LocalDB, err error) { 
	if LocalPolicyDB == nil {
		fmt.Println("local policy nil")
		err := errors.New("local policy nil")
		return nil,err
	}
	return LocalPolicyDB, err
}
