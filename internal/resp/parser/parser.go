package parser

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/DopamineNone/gedis/internal/resp/reply"
)

type Payload struct {
	Data reply.Reply
	Err  error
}

type readState struct {
	readingMutiLine   bool
	msgType           byte
	args              [][]byte
	expectedArgsCount int64
	bulkLen           int64
}

func (rs *readState) reset() {
	rs.readingMutiLine = false
	rs.msgType = '\n'
	rs.args = nil
	rs.expectedArgsCount = 0
	rs.bulkLen = 0
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse(reader, ch)
	return ch
}

func parse(reader io.Reader, out chan<- *Payload) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()
	bufReader, state := bufio.NewReader(reader), new(readState)

	for {
		msg, ioErr, err := readLine(bufReader, state)
		switch {
		case ioErr:
			out <- &Payload{
				Err: err,
			}
			close(out)
			return
		case err != nil:
			out <- &Payload{
				Err: err,
			}
			state.reset()
			continue
		}

		if !state.readingMutiLine {
			switch msg[0] {
			case '*':
				err := parseMultiBulkHeader(msg, state)
				if err != nil {
					out <- &Payload{
						Err: errors.New("unknown protocol: " + string(msg)),
					}
					state.reset()
					continue
				}
			case '$':
				err := parseBulkHeader(msg, state)
				if err != nil {
					out <- &Payload{
						Err: errors.New("unknown protocol: " + string(msg)),
					}
					state.reset()
					continue
				}
			case '+', '-', ':':
				result, err := parseSingleLineReply(msg)
				out <- &Payload{
					Data: result,
					Err:  err,
				}
				state.reset()
				continue
			default:
				out <- &Payload{
					Data: nil,
					Err:  errors.New("unknown protocol: " + string(msg)),
				}
			}
		} else {
			err := parseMultiLineBody(msg, state)
			if err != nil {
				out <- &Payload{
					Err: err,
				}
				state.reset()
				continue
			}
			if state.isFinished() {
				switch state.msgType {
				case '*':
					out <- &Payload{
						Data: reply.NewMultiBulkReply(state.args),
					}
				case '$':
					out <- &Payload{
						Data: reply.NewBulkReply(state.args[0]),
					}
				}
				state.reset()
			}
		}
	}
}

func parseSingleLineReply(msg []byte) (result reply.Reply, err error) {
	if len(msg) < 3 {
		return nil, errors.New("protocol error: " + string(msg))
	}
	content := strings.TrimSuffix(string(msg), "\r\n")
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(content)
	case '-':
		result = reply.NewErrReply(content)
	case ':':
		val, err := strconv.Atoi(string(content))
		if err != nil {
			return nil, errors.New("protocol error: " + string(msg))
		}
		result = reply.NewIntReply(val)
	}
	return
}

func parseBulkHeader(header []byte, state *readState) error {
	if len(header) < 3 {
		return errors.New("protocol error: " + string(header))
	}
	expectedLen, err := strconv.ParseInt(string(header[1:len(header)-2]), 10, 64)
	if err != nil || expectedLen < -1 {
		return errors.New("protocol error: " + string(header))
	}

	state.bulkLen = expectedLen
	if expectedLen == -1 {
		return nil
	}
	state.expectedArgsCount = 1
	state.msgType = header[0]
	state.readingMutiLine = true
	state.args = make([][]byte, 0, 1)
	return nil
}

func parseMultiBulkHeader(header []byte, state *readState) error {
	if len(header) < 3 {
		return errors.New("protocol error: " + string(header))
	}
	expectedLine, err := strconv.ParseInt(string(header[1:len(header)-2]), 10, 64)
	if err != nil || expectedLine < 0 {
		return errors.New("protocol error: " + string(header))
	}
	if expectedLine > 0 {
		state.expectedArgsCount = expectedLine
		state.msgType = header[0]
		state.readingMutiLine = true
		state.args = make([][]byte, 0, expectedLine)
	}
	return nil
}

func parseMultiLineBody(body []byte, state *readState) error {
	// start with '$'
	if body[0] == '$' {
		if len(body) < 3 {
			return errors.New("protocol error: " + string(body))
		}
		expectedLine, err := strconv.ParseInt(string(body[1:len(body)-2]), 10, 64)
		if err != nil {
			return errors.New("protocol error: " + string(body))
		}
		if expectedLine <= 0 {
			state.args = append(state.args, []byte{})
		} else {
			state.bulkLen = expectedLine
		}
	} else {
		if len(body) < 2 || !bytes.Equal(body[len(body)-2:], reply.CRLF) {
			return errors.New("protocol error: " + string(body))
		}
		state.args = append(state.args, body[:len(body)-2])
	}
	return nil
}

func readLine(reader *bufio.Reader, state *readState) (msg []byte, ioErr bool, err error) {
	if state.bulkLen == 0 {
		msg, err = reader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
	} else {
		msg = make([]byte, state.bulkLen+2)
		_, err = io.ReadFull(reader, msg)
		if err != nil {
			return nil, true, err
		}
	}
	// check if format is correct
	if len(msg) < 2 || msg[len(msg)-1] != '\n' || msg[len(msg)-2] != '\r' {
		return nil, false, errors.New("protocol error: " + string(msg))
	}
	state.bulkLen = 0
	return msg, false, nil
}

func (rs *readState) isFinished() bool {
	return rs.expectedArgsCount > 0 && rs.expectedArgsCount == int64(len(rs.args))
}
