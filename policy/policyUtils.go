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
type PolicyEngineDB struct {
	PolicyConditionsDB *patriciaDB.Trie
	localPolicyConditionsDB *LocalDBSlice
	PolicyActionsDB *patriciaDB.Trie
	localPolicyActionsDB *LocalDBSlice
	PolicyStmtDB *patriciaDB.Trie
	localPolicyStmtDB *LocalDBSlice
	PolicyDB *patriciaDB.Trie
	localPolicyDB *LocalDBSlice
    PolicyStmtPolicyMapDB map[string] []string //policies using this statement
    PrefixPolicyListDB *patriciaDB.Trie
    ProtocolPolicyListDB map[string][]string//policystmt names assoociated with every protocol type
    ImportPolicyPrecedenceMap map[int] string
    ExportPolicyPrecedenceMap map[int] string
}

func NewPolicyEngineDB() (policyEngineDB *PolicyEngineDB) {
   policyEngineDB = &PolicyEngineDB{}
   policyEngineDB.PolicyActionsDB = patriciaDB.NewTrie()
   localPolicyActionsDB := make([]LocalDB,0)
   localActionSlice := LocalDBSlice(localPolicyActionsDB)
   policyEngineDB.localPolicyActionsDB = &localActionSlice

   policyEngineDB.PolicyConditionsDB = patriciaDB.NewTrie()
   localPolicyConditionsDB := make([]LocalDB,0)
   localConditionSlice := LocalDBSlice(localPolicyConditionsDB)
   policyEngineDB.localPolicyConditionsDB = &localConditionSlice

   policyEngineDB.PolicyStmtDB = patriciaDB.NewTrie()
   localPolicyStmtDB := make([]LocalDB,0)
   localStmtSlice := LocalDBSlice(localPolicyStmtDB)
   policyEngineDB.localPolicyStmtDB = &localStmtSlice

   policyEngineDB.PolicyDB = patriciaDB.NewTrie()
   localPolicyDB := make([]LocalDB,0)
   localPolicySlice := LocalDBSlice(localPolicyDB)
   policyEngineDB.localPolicyDB = &localPolicySlice

   policyEngineDB.PolicyStmtPolicyMapDB = make(map[string] []string) 
   policyEngineDB.PrefixPolicyListDB = patriciaDB.NewTrie()
   policyEngineDB.ProtocolPolicyListDB = make(map[string][]string)
   policyEngineDB.ImportPolicyPrecedenceMap = make(map[int] string)
   policyEngineDB.ExportPolicyPrecedenceMap = make(map[int] string)
   return policyEngineDB
}
