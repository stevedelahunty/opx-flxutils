// policyUtils.go
package policy

import (
//	"fmt"
	"utils/patriciaDB"
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
type	 Policyfunc func(actionInfo interface{}, params interface{})
type	 PolicyCheckfunc func(params interface{}) bool
type EntityUpdatefunc func(details PolicyDetails, params interface{})

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
	DefaultImportPolicyActionFunc Policyfunc
	DefaultExportPolicyActionFunc Policyfunc
	IsEntityPresentFunc PolicyCheckfunc
	UpdateEntityDB EntityUpdatefunc
	ActionfuncMap map[int]Policyfunc
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
   policyEngineDB.PrefixPolicyListDB = patriciaDB.NewTrie()
   policyEngineDB.ProtocolPolicyListDB = make(map[string][]string)
   policyEngineDB.ImportPolicyPrecedenceMap = make(map[int] string)
   policyEngineDB.ExportPolicyPrecedenceMap = make(map[int] string)
   policyEngineDB.ActionfuncMap = make(map[int]Policyfunc)
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
