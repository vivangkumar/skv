package backend

// The below types represent non wire
// types for various operations supported
// by the store.

// SetCmd is the non wire type for 'set'.
type SetCmd struct {
	Key   string
	Value string
}

// GetCmd is the non wire type for 'get'.
type GetCmd struct {
	Key string
}

// DelCmd is the non wire type for 'delete'.
type DelCmd struct {
	Key string
}
