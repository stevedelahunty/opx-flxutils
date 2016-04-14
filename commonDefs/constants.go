package commonDefs

//L2 types
const (
	L2RefTypePort = iota
	L2RefTypeLag
	L2RefTypeVlan
	L2RefTypeVtep
	IfTypeP2P
	IfTypeBcast
	IfTypeLoopback
	IfTypeSecondary
	IfTypeVirtual
	IfTypeVxlan
	IfTypeNull
)

func GetIfTypeName(ifType int) string {
	switch ifType {
	case L2RefTypePort:
		return "Port"
	case L2RefTypeLag:
		return "Lag"
	case L2RefTypeVlan:
		return "Vlan"
	default:
		return "Unknown"
	}
}
