package commonDefs

//L2 types
const (
	IfTypePort = iota
	IfTypeLag
	IfTypeVlan
	IfTypeP2P
	IfTypeBcast
	IfTypeLoopback
	IfTypeSecondary
	IfTypeVirtual
	IfTypeNull
)

func GetIfTypeName(ifType int) string {
	switch ifType {
	case IfTypePort:
		return "Port"
	case IfTypeLag:
		return "Lag"
	case IfTypeVlan:
		return "Vlan"
	default:
		return "Unknown"
	}
}
