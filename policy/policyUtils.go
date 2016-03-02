// policyUtils.go
package policy

import (
	"fmt"
	"utils/patriciaDB"
	"utils/policy/policyCommonDefs"
)
const (
	add = iota
	del
	delAll
	invalidate
)
const (
	Invalid = -1
	Valid = 0
)
type ConditionsAndActionsList struct {
	ConditionList []string
	ActionList    []string
}
type PolicyStmtMap struct {
	PolicyStmtMap map[string]ConditionsAndActionsList
}
type PolicyEngineFilterEntityParams struct {
	DestNetIp      string	//CIDR format 
	NextHopIp      string
	RouteProtocol  string
	CreatePath     bool
	DeletePath     bool
	PolicyList    []string
	PolicyHitCounter int
}
//struct sent to the application for updating its local maps/DBs 
type PolicyDetails struct {
	Policy            string
	PolicyStmt        string
	ConditionList     []string
	ActionList        []string
	EntityDeleted     bool    //whether this policy/stmt resulted in deleting the entity
}

type LocalDB struct {
	Prefix  patriciaDB.Prefix
	IsValid bool
	Precedence int
}
type LocalDBSlice []LocalDB

func (slice *LocalDBSlice )updateLocalDB(prefix patriciaDB.Prefix) {
	localDBRecord := LocalDB{Prefix:prefix, IsValid:true}
	if(slice == nil) {
		return
	} 
	*slice = append(*slice, localDBRecord)

}
type	 Policyfunc func(actionInfo interface{}, conditionInfo interface {}, params interface{})
type PolicyConditionCheckfunc func(entity PolicyEngineFilterEntityParams, condition PolicyCondition, policyStmt PolicyStmt) bool
type UndoActionfunc func(actionInfo interface {}, conditionList []string, params interface{}, policyStmt PolicyStmt)
type	 PolicyCheckfunc func(params interface{}) bool
type EntityUpdatefunc func(details PolicyDetails, params interface{})
type PolicyApplyfunc func(entity PolicyEngineFilterEntityParams, policyData interface{}, params interface{})
type EntityTraverseAndApplyPolicyfunc func(data interface{}, updatefunc PolicyApplyfunc )
type PolicyEntityMapIndex interface {}

type PolicyEngineDB struct {
	PolicyConditionsDB *patriciaDB.Trie
	LocalPolicyConditionsDB *LocalDBSlice
	PolicyActionsDB *patriciaDB.Trie
	LocalPolicyActionsDB *LocalDBSlice
	PolicyStmtDB *patriciaDB.Trie
	LocalPolicyStmtDB *LocalDBSlice
	PolicyDB *patriciaDB.Trie
	LocalPolicyDB *LocalDBSlice
    PolicyStmtPolicyMapDB map[string] []string //policies using this statement
    PrefixPolicyListDB *patriciaDB.Trie
    ProtocolPolicyListDB map[string][]string//policystmt names assoociated with every protocol type
    ImportPolicyPrecedenceMap map[int] string
    ExportPolicyPrecedenceMap map[int] string
    PolicyEntityMap map[PolicyEntityMapIndex]PolicyStmtMap
	DefaultImportPolicyActionFunc Policyfunc
	DefaultExportPolicyActionFunc Policyfunc
	IsEntityPresentFunc PolicyCheckfunc
	GetPolicyEntityMapIndex func(entity PolicyEngineFilterEntityParams, policy string) PolicyEntityMapIndex
	UpdateEntityDB EntityUpdatefunc
	ConditionCheckfuncMap map[int]PolicyConditionCheckfunc
	ActionfuncMap map[int]Policyfunc
	UndoActionfuncMap map[int]UndoActionfunc
	TraverseAndApplyPolicyFunc EntityTraverseAndApplyPolicyfunc
	TraverseAndReversePolicyFunc func(interface {})
}

func (db*PolicyEngineDB) buildPolicyConditionCheckfuncMap () {
	fmt.Println("buildPolicyConditionCheckfuncMap")
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeDstIpPrefixMatch] = db.DstIpPrefixMatchConditionfunc
	db.ConditionCheckfuncMap[policyCommonDefs.PolicyConditionTypeProtocolMatch] = db.ProtocolMatchConditionfunc
}
func NewPolicyEngineDB() (policyEngineDB *PolicyEngineDB) {
   policyEngineDB = &PolicyEngineDB{}
   policyEngineDB.PolicyActionsDB = patriciaDB.NewTrie()
   LocalPolicyActionsDB := make([]LocalDB,0)
   localActionSlice := LocalDBSlice(LocalPolicyActionsDB)
   policyEngineDB.LocalPolicyActionsDB = &localActionSlice

   policyEngineDB.PolicyConditionsDB = patriciaDB.NewTrie()
   LocalPolicyConditionsDB := make([]LocalDB,0)
   localConditionSlice := LocalDBSlice(LocalPolicyConditionsDB)
   policyEngineDB.LocalPolicyConditionsDB = &localConditionSlice

   policyEngineDB.PolicyStmtDB = patriciaDB.NewTrie()
   LocalPolicyStmtDB := make([]LocalDB,0)
   localStmtSlice := LocalDBSlice(LocalPolicyStmtDB)
   policyEngineDB.LocalPolicyStmtDB = &localStmtSlice

   policyEngineDB.PolicyDB = patriciaDB.NewTrie()
   LocalPolicyDB := make([]LocalDB,0)
   localPolicySlice := LocalDBSlice(LocalPolicyDB)
   policyEngineDB.LocalPolicyDB = &localPolicySlice

   policyEngineDB.PolicyStmtPolicyMapDB = make(map[string] []string) 
   policyEngineDB.PolicyEntityMap = make(map[PolicyEntityMapIndex] PolicyStmtMap) 
   policyEngineDB.PrefixPolicyListDB = patriciaDB.NewTrie()
   policyEngineDB.ProtocolPolicyListDB = make(map[string][]string)
   policyEngineDB.ImportPolicyPrecedenceMap = make(map[int] string)
   policyEngineDB.ExportPolicyPrecedenceMap = make(map[int] string)
   policyEngineDB.ConditionCheckfuncMap = make(map[int] PolicyConditionCheckfunc)
   policyEngineDB.buildPolicyConditionCheckfuncMap()
   policyEngineDB.ActionfuncMap = make(map[int]Policyfunc)
   policyEngineDB.UndoActionfuncMap = make(map[int]UndoActionfunc)
   return policyEngineDB
}

func (db*PolicyEngineDB) SetDefaultImportPolicyActionFunc(defaultfunc Policyfunc){
	db.DefaultImportPolicyActionFunc = defaultfunc
}
func (db*PolicyEngineDB) SetDefaultExportPolicyActionFunc(defaultfunc Policyfunc){
	db.DefaultExportPolicyActionFunc = defaultfunc
}
func (db*PolicyEngineDB) SetIsEntityPresentFunc(IsPresent PolicyCheckfunc) {
	db.IsEntityPresentFunc = IsPresent
}
func (db*PolicyEngineDB) SetEntityUpdateFunc(updatefunc EntityUpdatefunc) {
	db.UpdateEntityDB = updatefunc
}
func (db *PolicyEngineDB) SetActionFunc(action int, setfunc Policyfunc) {
	db.ActionfuncMap[action] = setfunc
}
func (db *PolicyEngineDB) SetUndoActionFunc(action int, setfunc UndoActionfunc) {
	db.UndoActionfuncMap[action] = setfunc
}
func (db *PolicyEngineDB) SetTraverseAndApplyPolicyFunc(updatefunc EntityTraverseAndApplyPolicyfunc) {
	db.TraverseAndApplyPolicyFunc = updatefunc
}
func (db *PolicyEngineDB) SetTraverseAndReversePolicyFunc(updatefunc func(policyItem interface{})) {
	db.TraverseAndReversePolicyFunc = updatefunc
}
func (db *PolicyEngineDB) SetGetPolicyEntityMapIndexFunc(getfunc func(entity PolicyEngineFilterEntityParams, policy string) PolicyEntityMapIndex) {
	db.GetPolicyEntityMapIndex = getfunc
}
func isPolicyTypeSame(oldPolicy Policy, policy Policy) (same bool){
	if oldPolicy.ExportPolicy == policy.ExportPolicy && oldPolicy.ImportPolicy==policy.ImportPolicy {
		same = true
	}
	return same
}
func (db *PolicyEngineDB) AddPolicyEntityMapEntry(entity PolicyEngineFilterEntityParams, policy string, policyStmt string, conditionList []string, actionList []string) {
	fmt.Println("AddPolicyEntityMapEntry")
	var policyStmtMap PolicyStmtMap
	var conditionsAndActionsList ConditionsAndActionsList
	if db.PolicyEntityMap == nil {
		db.PolicyEntityMap = make(map[PolicyEntityMapIndex]PolicyStmtMap)
	}
    if db.GetPolicyEntityMapIndex == nil {
		return
	}
	policyEntityMapIndex := db.GetPolicyEntityMapIndex(entity, policy)
	if policyEntityMapIndex == nil {
		fmt.Println("policyEntityMapKey nil")
		return
	}
	policyStmtMap, ok:= db.PolicyEntityMap[policyEntityMapIndex]
	if !ok {
		policyStmtMap.PolicyStmtMap = make(map[string]ConditionsAndActionsList)
	}
	_, ok = policyStmtMap.PolicyStmtMap[policyStmt]
	if ok {
		fmt.Println("policy statement map for statement ", policyStmt, " already in place for policy ", policy)
		return
	} 
	conditionsAndActionsList.ConditionList = make([]string,0)
	conditionsAndActionsList.ActionList = make([]string,0)
	for i:=0;conditionList != nil && i<len(conditionList);i++ {
		conditionsAndActionsList.ConditionList = append(conditionsAndActionsList.ConditionList,conditionList[i])
	}
	for i:=0;actionList != nil && i<len(actionList);i++ {
		conditionsAndActionsList.ActionList = append(conditionsAndActionsList.ActionList,actionList[i])
	}
	policyStmtMap.PolicyStmtMap[policyStmt]=conditionsAndActionsList
	db.PolicyEntityMap[policyEntityMapIndex]=policyStmtMap
}
func (db *PolicyEngineDB) DeletePolicyEntityMapEntry(entity PolicyEngineFilterEntityParams, policy string) {
	fmt.Println("DeletePolicyEntityMapEntry for policy ", policy)
	if db.PolicyEntityMap == nil {
		fmt.Println("PolicyEntityMap empty")
		return
	}
    if db.GetPolicyEntityMapIndex == nil {
		return
	}
	policyEntityMapIndex := db.GetPolicyEntityMapIndex(entity, policy)
	if policyEntityMapIndex == nil {
		fmt.Println("policyEntityMapIndex nil")
		return
	}
	//PolicyRouteMap[policyRouteIndex].policyStmtMap=nil
	delete(db.PolicyEntityMap,policyEntityMapIndex)
}

