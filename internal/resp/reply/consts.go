package reply



var (
	_PongRelpy           = new(PongReply)
	_OkReply             = new(OkReply)
	_NullBulkReply       = new(NullBulkReply)
	_EmptyMultiBulkReply = new(EmptyMultiBulkReply)
	_NoRelpy             = new(NoReply)

	pongBytes           = []byte("+PONG\r\n")
	okBytes             = []byte("+OK\r\n")
	nullBulkBytes       = []byte("$-1\r\n")
	emptyMultiBulkBytes = []byte("*0\r\n")
	noBytes             = []byte("")
)

// PongReply when server received a 'ping' request, return a pong response
type PongReply struct{}

func (r *PongReply) ToBytes() []byte {
	return pongBytes
}

func NewPongReply() Reply {
	return _PongRelpy
}

type OkReply struct{}

func (r *OkReply) ToBytes() []byte {
	return okBytes
}

func NewOkReply() Reply {
	return _OkReply
}

// NullBulkReply retuns '$-1\r\n'
type NullBulkReply struct{}

func (r *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func NewNullBulkReply() Reply {
	return _NullBulkReply
}

// EmptyMultiBulkReply returns '*0\r\n'
type EmptyMultiBulkReply struct{}

func (r *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

func NewEmptyBulkReply() Reply {
	return _EmptyMultiBulkReply
}

// NoReply return ""
type NoReply struct{}

func (r *NoReply) ToBytes() []byte {
	return noBytes
}

func NewNoReply() Reply {
	return _NoRelpy
}
