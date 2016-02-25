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
   return policyEngineDB
}
