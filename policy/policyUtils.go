// policyUtils.go
package policy

import (
//	"fmt"
	"utils/patriciaDB"
)

type LocalDB struct {
	Prefix  patriciaDB.Prefix
	IsValid bool
	Precedence int
}
