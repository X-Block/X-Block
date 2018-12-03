
package operators

var (
	Call         = newPolymorphicOp(0x10, "call")
	CallIndirect = newPolymorphicOp(0x11, "call_indirect")
)
