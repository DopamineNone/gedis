package reply

var (
	_UnknownErrReply   = new(UnknowErrReply)
	_SyntaxErrReply    = new(SyntaxErrReply)
	_WrongTypeErrReply = new(WrongTypeErrReply)

	unknownErrBytes   = []byte("-ERR unknown\r\n")
	syntaxErrBytes    = []byte("-ERR syntax error\r\n")
	wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key holding the wrong kind of value\r\n")
)

type UnknowErrReply struct{}

func (r *UnknowErrReply) ToBytes() []byte {
	return unknownErrBytes
}

func NewUnknownErrReply() Reply {
	return _UnknownErrReply
}

type ArgNumErrReply struct {
	Cmd string
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + r.Cmd + "' command\r\n")
}

func NewArgNumErrReply(cmd string) Reply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

type SyntaxErrReply struct{}

func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func NewSyntaxErrReply() Reply {
	return _SyntaxErrReply
}

type WrongTypeErrReply struct{}

func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func NewWrongTypeErrReply() Reply {
	return _WrongTypeErrReply
}

type ProtocolErrReply struct{}

func (r *ProtocolErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func NewProtocolErrReply() Reply {
	return _SyntaxErrReply
}
