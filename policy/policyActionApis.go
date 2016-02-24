// ribdPolicyActionApis.go
package policy

import (
	"errors"
	"utils/policy/policyCommonDefs"
	"utils/patriciaDB"
	"strconv"
)
type RedistributeActionInfo struct {
	Redistribute bool
	RedistributeTargetProtocol string
}
type PolicyAction struct {
	Name          string
	ActionType int
	ActionInfo interface {}
	PolicyStmtList []string
	ActionGetBulkInfo string
	LocalDBSliceIdx int
}
type PolicyRedistributionAction struct{
	Redistribute string
	RedistributeTargetProtocol string
}

type PolicyActionConfig struct{
	Name string
	ActionType string
	SetAdminDistanceValue int
	Accept bool
	Reject bool
	RedistributeActionInfo  PolicyRedistributionAction    
}

func updateLocalActionsDB(prefix patriciaDB.Prefix, localPolicyActionsDB []localDB) {
    localDBRecord := localDB{Prefix:prefix, IsValid:true}
    if(localPolicyActionsDB == nil) {
		localPolicyActionsDB = make([]LocalDB, 0)
	} 
	localPolicyActionsDB = append(localPolicyActionsDB, localDBRecord)
}
func CreatePolicyRouteDispositionAction(PolicyActionsDB *patriciaDB.Trie, localPolicyActionsDB []LocalDB, cfg PolicyActionConfig)(val bool, err error) {
	logger.Println("CreateRouteDispositionAction")
	policyAction := PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   logger.Println("Defining a new policy action with name ", cfg.Name)
	   routeDispositionAction := ""
	   if cfg.Accept == true {
	      routeDispositionAction = "Accept"	
	   } else if cfg.Reject == true {
	      routeDispositionAction = "Reject"	
	   } else {
	      logger.Println("User should set either one of accept/reject to true for this action type")
		  err = errors.New("User should set either one of accept/reject to true for this action type")
		  return val,err	
	   }
	   newPolicyAction := PolicyAction{Name:cfg.Name,actionType:policyCommonDefs.PolicyActionTypeRouteDisposition,ActionInfo:routeDispositionAction ,LocalDBSliceIdx:(len(localPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo =   routeDispositionAction
		if ok := PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			logger.Println(" return value not ok")
			return val, err
		}
	  updateLocalActionsDB(patriciaDB.Prefix(cfg.Name), localPolicyActionsDB)
	} else {
		logger.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}

func CreatePolicyAdminDistanceAction(PolicyActionsDB *patriciaDB.Trie, localPolicyActionsDB []LocalDB, cfg PolicyActionConfig) (val bool, err error) {
	logger.Println("CreatePolicyAdminDistanceAction")
	policyAction := PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   logger.Println("Defining a new policy action with name ", cfg.Name, "Setting admin distance value to ", cfg.SetAdminDistanceValue)
	   newPolicyAction := PolicyAction{Name:cfg.Name,actionType:policyCommonDefs.PoilcyActionTypeSetAdminDistance,ActionInfo:cfg.SetAdminDistanceValue ,LocalDBSliceIdx:(len(localPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo =  "Set admin distance to value "+strconv.Itoa(int(cfg.SetAdminDistanceValue))
		if ok := PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			logger.Println(" return value not ok")
			return val, err
		}
	  updateLocalActionsDB(patriciaDB.Prefix(cfg.Name), localPolicyActionsDB)
	} else {
		logger.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}

func CreatePolicyRedistributionAction(PolicyActionsDB *patriciaDB.Trie, localPolicyActionsDB []LocalDB, cfg *ribd.PolicyActionConfig) (val bool, err error) {
	logger.Println("CreatePolicyRedistributionAction")

	policyAction := PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   logger.Println("Defining a new policy action with name ", cfg.Name)
	   redistributeActionInfo := RedistributeActionInfo{ RedistributeTargetProtocol:cfg.RedistributeActionInfo.RedistributeTargetProtocol}
       if cfg.RedistributeActionInfo.Redistribute == "Allow" {
	      redistributeActionInfo.Redistribute = true	
	   } else if cfg.RedistributeActionInfo.Redistribute == "Block" {
	      redistributeActionInfo.Redistribute = false	
	   } else {
	      logger.Println("Invalid redistribute option ",cfg.RedistributeActionInfo.Redistribute," - should be either Allow/Block")	
          err = errors.New("Invalid redistribute option")
		  return val,err
	   }
	   newPolicyAction := PolicyAction{Name:cfg.Name,ActionType:policyCommonDefs.PolicyActionTypeRouteRedistribute,actionInfo:redistributeActionInfo ,LocalDBSliceIdx:(len(localPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo = cfg.RedistributeActionInfo.Redistribute + " Redistribute to Target Protocol " + cfg.RedistributeActionInfo.RedistributeTargetProtocol
		if ok := PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			logger.Println(" return value not ok")
			return val, err
		}
	    updateLocalActionsDB(patriciaDB.Prefix(cfg.Name), localPolicyActionsDB)
	} else {
		logger.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}
func CreatePolicyAction(PolicyActionsDB *patriciaDB.Trie, localPolicyActionsDB []LocalDB, cfg *ribd.PolicyActionConfig) ( err error) {
	logger.Println("CreatePolicyAction")
	switch cfg.ActionType {
		case "RouteDisposition":
		   CreatePolicyRouteDispositionAction(PolicyActionsDB, localPolicyActionsDB, cfg)
		   break
		case "Redistribution":
		   CreatePolicyRedistributionAction(PolicyActionsDB, localPolicyActionsDB,cfg)
		   break
        case "SetAdminDistance":
		   CreatePolicyAdminDistanceAction(PolicyActionsDB, localPolicyActionsDB,cfg)
		   break
		default:
		   logger.Println("Unknown action type ", cfg.ActionType)
		   err = errors.New("Unknown action type")
	}
	return err
}
