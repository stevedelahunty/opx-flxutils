// policyEngine.go
package policy

import (
//	 "utils/patriciaDB"
//	  "utils/policy/policyCommonDefs"
//	 "reflect"
//	 "sort"
//	 "strconv"
//	"utils/commonDefs"
//	"net"
//	"asicdServices"
//	"asicd/asicdConstDefs"
//	"bytes"
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
func (db *PolicyEngineDB) PolicyEngineTraverseAndApplyPolicy(policy Policy) {
	fmt.Println("PolicyEngineTraverseAndApplyPolicy -  apply policy ", policy.Name)
    if policy.ExportPolicy || policy.ImportPolicy{
	   fmt.Println("Applying import/export policy to all routes")
	  // PolicyEngineTraverseAndApply(policy)
	} else if policy.GlobalPolicy {
		fmt.Println("Need to apply global policy")
		//policyEngineApplyGlobalPolicy(policy)
	}
}

func (db *PolicyEngineDB) PolicyEngineTraverseAndReversePolicy(policy Policy){
	fmt.Println("PolicyEngineTraverseAndReversePolicy -  reverse policy ", policy.Name)
    if policy.ExportPolicy || policy.ImportPolicy{
	   fmt.Println("Reversing import/export policy ")
	   //PolicyEngineTraverseAndReverse(policy)
	} else if policy.GlobalPolicy {
		fmt.Println("Need to reverse global policy")
		//policyEngineReverseGlobalPolicy(policy)
	}
	
}
func (db *PolicyEngineDB) PolicyEngineFilter(route interface{}, policyPath int, params interface{}) {
	fmt.Println("PolicyEngineFilter")
}


