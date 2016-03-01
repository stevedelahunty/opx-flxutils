// policyEngine.go
package policy

import (
	 "utils/patriciaDB"
	  "utils/policy/policyCommonDefs"
	  "utils/netUtils"
	  "strings"
	 "reflect"
	 "sort"
	 "strconv"
//	"utils/commonDefs"
//	"net"
//	"asicdServices"
//	"asicd/asicdConstDefs"
	"bytes"
  //  "database/sql"
   "fmt"
)
func (db *PolicyEngineDB) ActionListHasAction(actionList []string, actionType int, action string) (match bool) {
	fmt.Println("ActionListHasAction for action ", action)
	return match
}
func (db *PolicyEngineDB) PolicyEngineCheck(route interface{}, policyType int) (actionList []string){
	fmt.Println("PolicyEngineTest to see if there are any policies  ")
	return nil
}
func (db *PolicyEngineDB) PolicyEngineUndoActionsPolicyStmt(policy Policy, policyStmt PolicyStmt, params interface{}, conditionsAndActionsList ConditionsAndActionsList) {
	fmt.Println("policyEngineUndoActionsPolicyStmt")
	if conditionsAndActionsList.ActionList == nil {
		fmt.Println("No actions")
		return
	}
	var i int
	for i=0;i<len(conditionsAndActionsList.ActionList);i++ {
	  fmt.Printf("Find policy action number %d name %s in the action database\n", i, conditionsAndActionsList.ActionList[i])
	  actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[i]))
	  if actionItem == nil {
	     fmt.Println("Did not find action ", conditionsAndActionsList.ActionList[i], " in the action database")	
		 continue
	  }
	  actionInfo := actionItem.(PolicyAction)
	  if db.UndoActionfuncMap[actionInfo.ActionType] != nil {
			db.UndoActionfuncMap[actionInfo.ActionType](actionItem,conditionsAndActionsList.ConditionList, params, policyStmt)	
	  }
	}
}
func (db *PolicyEngineDB) PolicyEngineUndoPolicyForEntity(entity PolicyEngineFilterEntityParams, policy Policy, params interface{}) {
	fmt.Println("policyEngineUndoPolicyForRoute - policy name ", policy.Name, "  route: ", entity.DestNetIp, " type:", entity.RouteProtocol)
    policyEntityIndex := db.GetPolicyEntityMapIndex(entity,policy.Name)
	if policyEntityIndex == nil {
		fmt.Println("policy entity map index nil")
		return
	}
	policyStmtMap := db.PolicyEntityMap[policyEntityIndex]
	if policyStmtMap.PolicyStmtMap == nil{
		fmt.Println("Unexpected:None of the policy statements of this policy have been applied on this route")
		return
	}
	for stmt,conditionsAndActionsList:=range policyStmtMap.PolicyStmtMap {
		fmt.Println("Applied policyStmtName ",stmt)
		policyStmt :=db. PolicyStmtDB.Get(patriciaDB.Prefix(stmt))
        if policyStmt == nil {
			fmt.Println("Invalid policyStmt")
			continue
		}
		db.PolicyEngineUndoActionsPolicyStmt(policy,policyStmt.(PolicyStmt), params, conditionsAndActionsList)
		//check if the route still exists - it may have been deleted by the previous statement action
	   if db.IsEntityPresentFunc != nil {
		if !(db.IsEntityPresentFunc(params)) {
			fmt.Println("This entity no longer exists")
			break
		}
	  }
	}
}
func (db *PolicyEngineDB) PolicyEngineImplementActions(entity PolicyEngineFilterEntityParams, policyStmt PolicyStmt, params interface {}) (actionList []string){
	fmt.Println("policyEngineImplementActions")
	if policyStmt.Actions == nil {
		fmt.Println("No actions")
		return actionList
	}
	var i int
	addActionToList := false
	for i=0;i<len(policyStmt.Actions);i++ {
	  addActionToList = false
	  fmt.Printf("Find policy action number %d name %s in the action database\n", i, policyStmt.Actions[i])
	  actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[i]))
	  if actionItem == nil {
	     fmt.Println("Did not find action ", policyStmt.Actions[i], " in the action database")	
		 continue
	  }
	  action := actionItem.(PolicyAction)
	  fmt.Printf("policy action number %d type %d\n", i, action.ActionType)
		switch action.ActionType {
		   case policyCommonDefs.PolicyActionTypeRouteDisposition:
		      fmt.Println("PolicyActionTypeRouteDisposition action to be applied")
	           addActionToList = true
			  if db.ActionfuncMap[policyCommonDefs.PolicyActionTypeRouteDisposition] != nil {
			     db.ActionfuncMap[policyCommonDefs.PolicyActionTypeRouteDisposition](action.ActionInfo,nil, params)	
			  }
			  break
		   case policyCommonDefs.PolicyActionTypeRouteRedistribute:
		      fmt.Println("PolicyActionTypeRouteRedistribute action to be applied")
			  if db.ActionfuncMap[policyCommonDefs.PolicyActionTypeRouteRedistribute] != nil {
			     db.ActionfuncMap[policyCommonDefs.PolicyActionTypeRouteRedistribute](action.ActionInfo,nil, params)	
			  }
	          addActionToList = true
			  break
		   default:
		      fmt.Println("UnknownInvalid type of action")
			  break
		}
		if addActionToList == true {
		   if actionList == nil {
		      actionList = make([]string,0)
		   }
	       actionList = append(actionList,action.Name)
		}
	}
    return actionList
}
func (db *PolicyEngineDB) FindPrefixMatch(ipAddr string, ipPrefix patriciaDB.Prefix, policyName string)(match bool){
    fmt.Println("Prefix match policy ", policyName)
	policyListItem := db.PrefixPolicyListDB.GetLongestPrefixNode(ipPrefix)
	if policyListItem == nil {
		fmt.Println("intf stored at prefix ", ipPrefix, " is nil")
		return false
	}
    if policyListItem != nil && reflect.TypeOf(policyListItem).Kind() != reflect.Slice {
		fmt.Println("Incorrect data type for this prefix ")
		 return false
	}
	policyListSlice := reflect.ValueOf(policyListItem)
	for idx :=0;idx < policyListSlice.Len();idx++ {
	   prefixPolicyListInfo := policyListSlice.Index(idx).Interface().(PrefixPolicyListInfo)
	   if prefixPolicyListInfo.policyName != policyName {
	      fmt.Println("Found a potential match for this prefix but the policy ", policyName, " is not what we are looking for")
		  continue
	   }
	   if prefixPolicyListInfo.lowRange == -1 && prefixPolicyListInfo.highRange == -1 {
          fmt.Println("Looking for exact match condition for prefix ", prefixPolicyListInfo.ipPrefix)
		  if bytes.Equal(ipPrefix, prefixPolicyListInfo.ipPrefix) {
			 fmt.Println(" Matched the prefix")
	         return true
		  }	else {
			 fmt.Println(" Did not match the exact prefix")
		     return false	
		  }
	   }
	   tempSlice:=strings.Split(ipAddr,"/")
	   maskLen,err:= strconv.Atoi(tempSlice[1])
	   if err != nil {
	       fmt.Println("err getting maskLen")
		   return false	
	   }
	   fmt.Println("Mask len = ", maskLen)
	   if maskLen < prefixPolicyListInfo.lowRange || maskLen > prefixPolicyListInfo.highRange {
	      fmt.Println("Mask range of the route ", maskLen , " not within the required mask range:", prefixPolicyListInfo.lowRange,"..", prefixPolicyListInfo.highRange)	
		  return false
	   } else {
	      fmt.Println("Mask range of the route ", maskLen , " within the required mask range:", prefixPolicyListInfo.lowRange,"..", prefixPolicyListInfo.highRange)	
		  return true
	   }
	} 
	return match
}
func (db *PolicyEngineDB) DstIpPrefixMatchConditionfunc (entity PolicyEngineFilterEntityParams, condition PolicyCondition, policyStmt PolicyStmt) (match bool) {
	fmt.Println("dstIpPrefixMatchConditionfunc")
	ipPrefix,err := netUtils.GetNetworkPrefixFromCIDR(entity.DestNetIp)
	if err != nil {
		fmt.Println("Invalid ipPrefix for the route ", entity.DestNetIp)
		return false
	}
	match = db.FindPrefixMatch(entity.DestNetIp, ipPrefix,policyStmt.Name)
	if match {
		fmt.Println("Found a match for this prefix")
	}
	return match
}
func (db *PolicyEngineDB) ProtocolMatchConditionfunc (entity PolicyEngineFilterEntityParams, condition PolicyCondition, policyStmt PolicyStmt) (match bool) {
	fmt.Println("protocolMatchConditionfunc")
	matchProto := condition.ConditionInfo.(string)
	if matchProto == entity.RouteProtocol {
	   fmt.Println("Protocol condition matches")
	   match = true
	} 
	return match
}
func (db *PolicyEngineDB) ConditionCheckValid(entity PolicyEngineFilterEntityParams,conditionsList []string, policyStmt PolicyStmt) (valid bool) {
   fmt.Println("conditionCheckValid")	
   valid = true
   if conditionsList == nil {
      fmt.Println("No conditions to match, so valid")
	  return true	
   }
   for i:=0;i<len(conditionsList);i++ {
	  fmt.Printf("Find policy condition number %d name %s in the condition database\n", i, policyStmt.Conditions[i])
	  conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(conditionsList[i]))
	  if conditionItem == nil {
	     fmt.Println("Did not find condition ", conditionsList[i], " in the condition database")	
		 continue
	  }
	  condition := conditionItem.(PolicyCondition)
	  fmt.Printf("policy condition number %d type %d\n", i, condition.ConditionType)
	  if db.ConditionCheckfuncMap[condition.ConditionType] != nil {
	      match := db.ConditionCheckfuncMap[condition.ConditionType](entity,condition,policyStmt)
		  if !match {
			fmt.Println("Condition does not match")
			return false
		  }
	  }
	}
   fmt.Println("returning valid= ", valid)
   return valid
}
func (db *PolicyEngineDB) PolicyEngineMatchConditions(entity PolicyEngineFilterEntityParams, policyStmt PolicyStmt) (match bool, conditionsList []string){
    fmt.Println("policyEngineMatchConditions")
	var i int
	allConditionsMatch := true
	anyConditionsMatch := false
	addConditiontoList := false
	for i=0;i<len(policyStmt.Conditions);i++ {
	  addConditiontoList = false
	  fmt.Printf("Find policy condition number %d name %s in the condition database\n", i, policyStmt.Conditions[i])
	  conditionItem := db.PolicyConditionsDB.Get(patriciaDB.Prefix(policyStmt.Conditions[i]))
	  if conditionItem == nil {
	     fmt.Println("Did not find condition ", policyStmt.Conditions[i], " in the condition database")	
		 continue
	  }
	  condition := conditionItem.(PolicyCondition)
	  fmt.Printf("policy condition number %d type %d\n", i, condition.ConditionType)
	  if db.ConditionCheckfuncMap[condition.ConditionType] != nil {
	      match = db.ConditionCheckfuncMap[condition.ConditionType](entity,condition,policyStmt)
		  if match {
			fmt.Println("Condition match found")
			anyConditionsMatch = true
			addConditiontoList = true
		  }
	  }
	  if addConditiontoList == true{
		if conditionsList == nil {
		   conditionsList = make([]string,0)
		}
		conditionsList = append(conditionsList,condition.Name)
	  }
	}
   if policyStmt.MatchConditions == "all" && allConditionsMatch == true {
	return true,conditionsList
   }
   if policyStmt.MatchConditions == "any" && anyConditionsMatch == true {
	return true,conditionsList
   }
    return match,conditionsList
}
func (db *PolicyEngineDB) PolicyEngineApplyPolicyStmt(entity *PolicyEngineFilterEntityParams, policy Policy, policyStmt PolicyStmt, policyPath int, params interface{}, hit *bool, deleted *bool) {
	fmt.Println("policyEngineApplyPolicyStmt - ", policyStmt.Name)
	var conditionList []string
	if policyStmt.Conditions == nil {
		fmt.Println("No policy conditions")
		*hit=true
	} else {
	   match,ret_conditionList := db.PolicyEngineMatchConditions(*entity, policyStmt)
	   fmt.Println("match = ", match)
	   *hit = match
	   if !match {
		   fmt.Println("Conditions do not match")
		   return
	   }
	   if ret_conditionList != nil {
		 if conditionList == nil {
			conditionList = make([]string,0)
		 }
		 for j:=0;j<len(ret_conditionList);j++ {
			conditionList =append(conditionList,ret_conditionList[j])
		 }
	   }
	}
	actionList := db.PolicyEngineImplementActions(*entity, policyStmt, params)
	if db.ActionListHasAction(actionList, policyCommonDefs.PolicyActionTypeRouteDisposition,"Reject") {
		fmt.Println("Reject action was applied for this entity")
		*deleted = true
	}
	//check if the route still exists - it may have been deleted by the previous statement action
	if db.IsEntityPresentFunc != nil {
		*deleted = !(db.IsEntityPresentFunc(params))
	}
	if db.UpdateEntityDB != nil {
		policyDetails := PolicyDetails{Policy:policy.Name, PolicyStmt:policyStmt.Name,ConditionList:conditionList,ActionList:actionList, EntityDeleted:*deleted}
		db.UpdateEntityDB(policyDetails,params)
	}
	if entity.CreatePath == true {
		db.AddPolicyEntityMapEntry(*entity, policy.Name, policyStmt.Name, conditionList, actionList)
	}
}

func (db *PolicyEngineDB) PolicyEngineApplyPolicy(entity *PolicyEngineFilterEntityParams, policy Policy, policyPath int,params interface{}, hit *bool) {
	fmt.Println("policyEngineApplyPolicy - ", policy.Name)
     var policyStmtKeys []int
	 deleted := false
	 for k:=range policy.PolicyStmtPrecedenceMap {
		fmt.Println("key k = ", k)
		policyStmtKeys = append(policyStmtKeys,k)
	}
	sort.Ints(policyStmtKeys)
	for i:=0;i<len(policyStmtKeys);i++ {
		fmt.Println("Key: ", policyStmtKeys[i], " policyStmtName ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
		policyStmt := db.PolicyStmtDB.Get((patriciaDB.Prefix(policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])))
        if policyStmt == nil {
			fmt.Println("Invalid policyStmt")
			continue
		}
		db.PolicyEngineApplyPolicyStmt(entity,policy,policyStmt.(PolicyStmt),policyPath, params, hit, &deleted)
		if deleted == true {
			fmt.Println("Entity was deleted as a part of the policyStmt ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
             break
		}
		if *hit == true {
			if policy.MatchType == "any" {
				fmt.Println("Match type for policy ", policy.Name, " is any and the policy stmt ", (policyStmt.(PolicyStmt)).Name, " is a hit, no more policy statements will be executed")
				break
			}
		}
	}
}
func (db *PolicyEngineDB) PolicyEngineApplyForEntity(entity PolicyEngineFilterEntityParams, policyData interface{}, params interface{}) {
   fmt.Println("policyEngineApplyForEntity" )	
   policy := policyData.(Policy)
   policyHit := false
     if len(entity.PolicyList) == 0 {
	  fmt.Println("This route has no policy applied to it so far, just apply the new policy")
      db.PolicyEngineApplyPolicy(&entity, policy, policyCommonDefs.PolicyPath_All,params, &policyHit)
     } else {
      fmt.Println("This route already has policy applied to it - len(route.PolicyList) - ", len(entity.PolicyList))
    
	  for i:=0;i<len(entity.PolicyList);i++ {
		 fmt.Println("policy at index ", i)
	     policyInfo := db.PolicyDB.Get(patriciaDB.Prefix(entity.PolicyList[i]))
	     if policyInfo == nil {
		    fmt.Println("Unexpected: Invalid policy in the route policy list")
	     } else {
	       oldPolicy := policyInfo.(Policy)
		   if !isPolicyTypeSame(oldPolicy, policy) {
			 fmt.Println("The policy type applied currently is not the same as new policy, so apply new policy")
              db.PolicyEngineApplyPolicy(&entity, policy, policyCommonDefs.PolicyPath_All,params, &policyHit)
		   } else if oldPolicy.Precedence < policy.Precedence {
			 fmt.Println("The policy types are same and precedence of the policy applied currently is lower than the new policy, so do nothing")
			 return 
		   } else {
			fmt.Println("The new policy's precedence is lower, so undo old policy's actions and apply the new policy")
			db.PolicyEngineUndoPolicyForEntity(entity, oldPolicy, params)
			db.PolicyEngineApplyPolicy(&entity, policy, policyCommonDefs.PolicyPath_All,params, &policyHit)
		   }
		}
	  }	
    }
}

func (db *PolicyEngineDB) PolicyEngineApplyGlobalPolicyStmt(policy Policy, policyStmt PolicyStmt) {
	fmt.Println("policyEngineApplyGlobalPolicyStmt - ", policyStmt.Name)
    var conditionItem interface{}=nil
//global policies can only have statements with 1 condition and 1 action
	if policyStmt.Actions == nil {
		fmt.Println("No policy actions defined")
		return
	}
	if policyStmt.Conditions == nil {
		fmt.Println("No policy conditions")
	} else {
		if len(policyStmt.Conditions) > 1 {
			fmt.Println("only 1 condition allowed for global policy stmt")
			return
		}
		conditionItem = db.PolicyConditionsDB.Get(patriciaDB.Prefix(policyStmt.Conditions[0]))
		if conditionItem == nil {
			fmt.Println("Condition ", policyStmt.Conditions[0]," not found")
			return
		}
		actionItem := db.PolicyActionsDB.Get(patriciaDB.Prefix(policyStmt.Actions[0]))
		if actionItem == nil {
			fmt.Println("Action ", policyStmt.Actions[0]," not found")
			return
		}
		actionInfo := actionItem.(PolicyAction)
		if db.ActionfuncMap[actionInfo.ActionType] != nil {
			db.ActionfuncMap[actionInfo.ActionType](actionItem,conditionItem,nil)	
		}
	}
}

func (db *PolicyEngineDB) PolicyEngineApplyGlobalPolicy(policy Policy) {
	fmt.Println("policyEngineApplyGlobalPolicy")
     var policyStmtKeys []int
	 for k:=range policy.PolicyStmtPrecedenceMap {
		fmt.Println("key k = ", k)
		policyStmtKeys = append(policyStmtKeys,k)
	}
	sort.Ints(policyStmtKeys)
	for i:=0;i<len(policyStmtKeys);i++ {
		fmt.Println("Key: ", policyStmtKeys[i], " policyStmtName ", policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])
		policyStmt := db.PolicyStmtDB.Get((patriciaDB.Prefix(policy.PolicyStmtPrecedenceMap[policyStmtKeys[i]])))
        if policyStmt == nil {
			fmt.Println("Invalid policyStmt")
			continue
		}
		db.PolicyEngineApplyGlobalPolicyStmt(policy,policyStmt.(PolicyStmt))
	}
}

func (db *PolicyEngineDB) PolicyEngineTraverseAndApplyPolicy(policy Policy) {
	fmt.Println("PolicyEngineTraverseAndApplyPolicy -  apply policy ", policy.Name)
    if policy.ExportPolicy || policy.ImportPolicy{
	   fmt.Println("Applying import/export policy to all routes")
	   if db.TraverseAndApplyPolicyFunc != nil {
	      fmt.Println("Calling TraverseAndApplyPolicyFunc function")
	      db.TraverseAndApplyPolicyFunc(policy, db.PolicyEngineApplyForEntity)	
	   }
	} else if policy.GlobalPolicy {
		fmt.Println("Need to apply global policy")
		db.PolicyEngineApplyGlobalPolicy(policy)
	}
}

func (db *PolicyEngineDB) PolicyEngineTraverseAndReversePolicy(policy Policy){
	fmt.Println("PolicyEngineTraverseAndReversePolicy -  reverse policy ", policy.Name)
    if policy.ExportPolicy || policy.ImportPolicy{
	   fmt.Println("Reversing import/export policy ")
	   db.TraverseAndReversePolicyFunc(policy)
	} else if policy.GlobalPolicy {
		fmt.Println("Need to reverse global policy")
		//policyEngineReverseGlobalPolicy(policy)
	}
	
}
func (db *PolicyEngineDB) PolicyEngineFilter(entity PolicyEngineFilterEntityParams, policyPath int, params interface{}) {
	fmt.Println("PolicyEngineFilter")
	var policyPath_Str string
	if policyPath == policyCommonDefs.PolicyPath_Import {
	   policyPath_Str = "Import"
	} else if policyPath == policyCommonDefs.PolicyPath_Export {
	   policyPath_Str = "Export"
	} else if policyPath == policyCommonDefs.PolicyPath_All {
		policyPath_Str = "ALL"
		fmt.Println("policy path ", policyPath_Str, " unexpected in this function")
		return
	}
	fmt.Println("PolicyEngineFilter for policypath ", policyPath_Str, "create = ", entity.CreatePath, " delete = ", entity.DeletePath, " route: ", entity.DestNetIp, " protocol type: ", entity.RouteProtocol)
    var policyKeys []int
	var policyHit bool
	idx :=0
	var policyInfo interface{}
	if policyPath == policyCommonDefs.PolicyPath_Import{
	   for k:=range db.ImportPolicyPrecedenceMap {
	      policyKeys = append(policyKeys,k)
	   }
	} else if policyPath == policyCommonDefs.PolicyPath_Export{
	   for k:=range db.ExportPolicyPrecedenceMap {
	      policyKeys = append(policyKeys,k)
	   }
	}
	sort.Ints(policyKeys)
	for ;; {
		if entity.DeletePath == true {		//policyEngineFilter called during delete
			if entity.PolicyList != nil {
             if idx >= len(entity.PolicyList) {
				break
			 } 		
		     fmt.Println("getting policy ", idx, " from entity.PolicyList")
	         policyInfo = 	db.PolicyDB.Get(patriciaDB.Prefix(entity.PolicyList[idx]))
		     idx++
			 if policyInfo.(Policy).ExportPolicy && policyPath == policyCommonDefs.PolicyPath_Import || policyInfo.(Policy).ImportPolicy && policyPath == policyCommonDefs.PolicyPath_Export {
				fmt.Println("policy ", policyInfo.(Policy).Name, " not the same type as the policypath -", policyPath_Str)
				continue
			 } 
	        } else {
		      fmt.Println("PolicyList empty and this is a delete operation, so break")
               break
	        }		
	    }  else if entity.CreatePath == true{ //policyEngine filter called during create 
			fmt.Println("idx = ", idx, " len(policyKeys):", len(policyKeys))
            if idx >= len(policyKeys) {
				break
			}		
			policyName := ""
            if policyPath == policyCommonDefs.PolicyPath_Import {
               policyName = db.ImportPolicyPrecedenceMap[policyKeys[idx]]
			} else if policyPath == policyCommonDefs.PolicyPath_Export {
               policyName = db.ExportPolicyPrecedenceMap[policyKeys[idx]]
			}
		    fmt.Println("getting policy  ", idx, " policyKeys[idx] = ", policyKeys[idx]," ", policyName," from PolicyDB")
             policyInfo = db.PolicyDB.Get((patriciaDB.Prefix(policyName)))
			idx++
	      }
	      if policyInfo == nil {
	        fmt.Println("Nil policy")
		    continue
	      }
	      policy := policyInfo.(Policy)
		  localPolicyDB := *db.LocalPolicyDB
	      if localPolicyDB != nil && localPolicyDB[policy.LocalDBSliceIdx].IsValid == false {
	        fmt.Println("Invalid policy at localDB slice idx ", policy.LocalDBSliceIdx)
		    continue	
	      }		
	      db.PolicyEngineApplyPolicy(&entity, policy, policyPath, params, &policyHit)
	      if policyHit {
	         fmt.Println("Policy ", policy.Name, " applied to the route")	
		     break
	      }
	}
	if entity.PolicyHitCounter == 0{
		fmt.Println("Need to apply default policy, policyPath = ", policyPath, "policyPath_Str= ", policyPath_Str)
		if policyPath == policyCommonDefs.PolicyPath_Import {
		   fmt.Println("Applying default import policy")
			if db.DefaultImportPolicyActionFunc != nil {
				db.DefaultImportPolicyActionFunc(nil,nil,params)
			}
		} else if policyPath == policyCommonDefs.PolicyPath_Export {
			fmt.Println("Applying default export policy")
			if db.DefaultExportPolicyActionFunc != nil {
				db.DefaultExportPolicyActionFunc(nil,nil,params)
			}
		}
	}
	if entity.DeletePath == true {
		db.DeletePolicyEntityMapEntry(entity,"")
	}
}


