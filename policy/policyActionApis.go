// ribdPolicyActionApis.go
package policy

import (
	"errors"
	"fmt"
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

type PolicyActionConfig struct{
	Name string
	ActionType string
	SetAdminDistanceValue int
	Accept bool
	Reject bool
	RedistributeAction  string
	RedistributeTargetProtocol string
}

func (db * PolicyEngineDB) CreatePolicyRouteDispositionAction(cfg PolicyActionConfig)(val bool, err error) {
	fmt.Println("CreateRouteDispositionAction")
	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   fmt.Println("Defining a new policy action with name ", cfg.Name)
	   routeDispositionAction := ""
	   if cfg.Accept == true {
	      routeDispositionAction = "Accept"	
	   } else if cfg.Reject == true {
	      routeDispositionAction = "Reject"	
	   } else {
	      fmt.Println("User should set either one of accept/reject to true for this action type")
		  err = errors.New("User should set either one of accept/reject to true for this action type")
		  return val,err	
	   }
	   newPolicyAction := PolicyAction{Name:cfg.Name,ActionType:policyCommonDefs.PolicyActionTypeRouteDisposition,ActionInfo:routeDispositionAction ,LocalDBSliceIdx:(len(*db.LocalPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo =   routeDispositionAction
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			fmt.Println(" return value not ok")
			return val, err
		}
	  db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name))
	} else {
		fmt.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}

func (db * PolicyEngineDB) CreatePolicyAdminDistanceAction(cfg PolicyActionConfig) (val bool, err error) {
	fmt.Println("CreatePolicyAdminDistanceAction")
	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   fmt.Println("Defining a new policy action with name ", cfg.Name, "Setting admin distance value to ", cfg.SetAdminDistanceValue)
	   newPolicyAction := PolicyAction{Name:cfg.Name,ActionType:policyCommonDefs.PoilcyActionTypeSetAdminDistance,ActionInfo:cfg.SetAdminDistanceValue ,LocalDBSliceIdx:(len(*db.LocalPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo =  "Set admin distance to value "+strconv.Itoa(int(cfg.SetAdminDistanceValue))
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			fmt.Println(" return value not ok")
			return val, err
		}
	  db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name))
	} else {
		fmt.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}

func (db * PolicyEngineDB) CreatePolicyRedistributionAction(cfg PolicyActionConfig) (val bool, err error) {
	fmt.Println("CreatePolicyRedistributionAction")

	policyAction := db.PolicyActionsDB.Get(patriciaDB.Prefix(cfg.Name))
	if(policyAction == nil) {
	   fmt.Println("Defining a new policy action with name ", cfg.Name)
	   redistributeActionInfo := RedistributeActionInfo{ RedistributeTargetProtocol:cfg.RedistributeTargetProtocol}
       if cfg.RedistributeAction == "Allow" {
	      redistributeActionInfo.Redistribute = true	
	   } else if cfg.RedistributeAction == "Block" {
	      redistributeActionInfo.Redistribute = false	
	   } else {
	      fmt.Println("Invalid redistribute option ",cfg.RedistributeAction," - should be either Allow/Block")	
          err = errors.New("Invalid redistribute option")
		  return val,err
	   }
	   newPolicyAction := PolicyAction{Name:cfg.Name,ActionType:policyCommonDefs.PolicyActionTypeRouteRedistribute,ActionInfo:redistributeActionInfo ,LocalDBSliceIdx:(len(*db.LocalPolicyActionsDB))}
       newPolicyAction.ActionGetBulkInfo = cfg.RedistributeAction + " Redistribute to Target Protocol " + cfg.RedistributeTargetProtocol
		if ok := db.PolicyActionsDB.Insert(patriciaDB.Prefix(cfg.Name), newPolicyAction); ok != true {
			fmt.Println(" return value not ok")
			return val, err
		}
	    db.LocalPolicyActionsDB.updateLocalDB(patriciaDB.Prefix(cfg.Name))
	} else {
		fmt.Println("Duplicate action name")
		err = errors.New("Duplicate policy action definition")
		return val, err
	}
	return val, err
}
func (db *PolicyEngineDB) CreatePolicyAction(cfg PolicyActionConfig) ( err error) {
	fmt.Println("CreatePolicyAction")
	switch cfg.ActionType {
		case "RouteDisposition":
		   db.CreatePolicyRouteDispositionAction(cfg)
		   break
		case "Redistribution":
		   db.CreatePolicyRedistributionAction(cfg)
		   break
        case "SetAdminDistance":
		   db.CreatePolicyAdminDistanceAction(cfg)
		   break
		default:
		   fmt.Println("Unknown action type ", cfg.ActionType)
		   err = errors.New("Unknown action type")
	}
	return err
}
/*
func GetPolicyActionsDB() (db *patriciaDB.Trie, err error) { 
	if PolicyActionsDB == nil {
		fmt.Println("policyActions nil")
		err := errors.New("policyActions nil")
		return nil,err
	}
	return PolicyActionsDB, err
}
func GetLocalPolicyActionsDB()(db []LocalDB, err error) { 
	if LocalPolicyActionsDB == nil {
		fmt.Println("local policyActions nil")
		err := errors.New("local policyActions nil")
		return nil,err
	}
	return LocalPolicyActionsDB, err
}
*/