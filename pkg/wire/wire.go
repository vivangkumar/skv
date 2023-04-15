package wire

const (
	CmdDecodeErr   = "cmd_decode"
	CmdFailedErr   = "cmd_failed"
	CmdArgMismatch = "cmd_arg_mismatch"
)

type Operation int

const (
	Set Operation = iota
	Get
	Del
)

var ops = map[Operation]string{
	Set: "set",
	Get: "get",
	Del: "del",
}

var reverseOps = map[string]Operation{
	"set": Set,
	"get": Get,
	"del": Del,
}

var expectedArgsCount = map[Operation]int{
	Set: 2,
	Get: 1,
	Del: 1,
}

// String implements the Stringer interface
func (o Operation) String() string {
	return ops[o]
}

// Cmd represents a command that has been
// sent over the wire
type Cmd struct {
	// Type is the type of operation
	// set, get or del
	Type Operation

	// Args contains the arguments passed
	// to commands
	//
	// In case of set, the args will contain
	// the key followed by the value in the
	// order that they were passed in
	Args []string
}

type Reply int

const (
	OK Reply = iota
	Err
)

// CmdReply represents a reply to a command
// encoded onto the wire
type CmdReply struct {
	// Type is either OK or Err
	Type Reply

	// Args contains any arguments to be written
	// over the wire
	Args []string
}
