type=$1 
test -n "$type" || exit 1
cat <<EOF
// convert a reference from the C $type to the Go $type
func as$type(cRef *C.$type) *$type {
	return (*$type)(unsafe.Pointer(cRef.goRef))
}

//export newGo$type
func newGo${type}(cRef *C.${type}) unsafe.Pointer {
	goRef := &${type}{cRef}
	cRef.goRef = unsafe.Pointer(goRef)
	return cRef.goRef
}
EOF
