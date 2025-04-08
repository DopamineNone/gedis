package reply

import (
	"bytes"
	"strconv"
)

var CRLF = []byte("\r\n")

type Reply interface {
	ToBytes() []byte // "bytes" => "$5\r\n bytes\r\n"
}

type BulkReply struct {
	Args []byte
}

func (r *BulkReply) ToBytes() []byte {
	if len(r.Args) == 0 {
		return nullBulkBytes
	}
	buf := new(bytes.Buffer)
	num := strconv.Itoa(len(r.Args))

	buf.Grow(5 + len(r.Args) + len(num))

	buf.WriteByte('$')
	buf.WriteString(num)
	buf.Write(CRLF)
	buf.Write(r.Args)
	buf.Write(CRLF)
	return buf.Bytes()
}

func NewBulkReply(args []byte) Reply {
	return &BulkReply{
		Args: args,
	}
}

type MultiBulkReply struct {
	Args [][]byte
}

func (r *MultiBulkReply) ToBytes() []byte {
	if len(r.Args) == 0 {
		return emptyMultiBulkBytes
	}
	buf := new(bytes.Buffer)

	buf.WriteByte('*')
	buf.WriteString(strconv.Itoa(len(r.Args)))
	buf.Write(CRLF)

	for _, bulk := range r.Args {
		if bulk == nil {
			buf.Write(nullBulkBytes)
			buf.Write(CRLF)
			continue
		}
		buf.WriteByte('$')
		buf.WriteString(strconv.Itoa(len(bulk)))
		buf.Write(CRLF)
		buf.Write(bulk)
		buf.Write(CRLF)
	}

	return buf.Bytes()
}

func NewMultiBulkReply(args [][]byte) Reply {
	return &MultiBulkReply{
		Args: args,
	}
}

type StatusReply struct {
	Status string
}

func (r *StatusReply) ToBytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte('+')
	buf.WriteString(r.Status)
	buf.Write(CRLF)
	return buf.Bytes()
}

func NewStatusReply(status string) Reply {
	return &StatusReply{
		Status: status,
	}
}

type ErrReply struct {
	Err string
}

func (r *ErrReply) ToBytes() []byte {
	buf := new(bytes.Buffer)
	buf.WriteByte('-')
	buf.WriteString(r.Err)
	buf.Write(CRLF)
	return buf.Bytes()
}

func NewErrReply(err string) Reply {
	return &ErrReply{
		Err: err,
	}
}

type IntReply struct {
	Code int
}

func (r *IntReply) ToBytes() []byte {
	return append([]byte(":"+strconv.Itoa(r.Code)), CRLF...)
}

func NewIntReply(code int) Reply {
	return &IntReply{
		Code: code,
	}
}
