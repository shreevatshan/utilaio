package utilaio

import (
	"bufio"
	"net"
	"time"
)

const (
	PROTOCOL = "tcp"
	TIMEOUT  = 5
	TYPE_UNI = "uni"
	TYPE_BI  = "bi"
)

type TRequest struct {
	Timeout int
	Address string
	Type    string
	Delim   byte
	Body    []byte
}

type TResponse struct {
	Err  error
	Body []byte
}

func (req *TRequest) Send() TResponse {
	var res TResponse
	var err error

	conn, err := net.Dial(PROTOCOL, req.Address)
	if err != nil {
		res.Err = err
		return res
	}

	_, err = conn.Write(req.Body)
	if err != nil {
		res.Err = err
		conn.Close()
		return res
	}

	if req.Type != TYPE_UNI {

		if req.Timeout == 0 {
			req.Timeout = TIMEOUT
		}
		conn.SetReadDeadline(time.Now().Add(time.Duration(req.Timeout) * time.Second))

		buf, err := bufio.NewReader(conn).ReadBytes(req.Delim)
		if err != nil {
			res.Err = err
			conn.Close()
			return res
		}
		res.Body = buf
	}
	conn.Close()
	return res
}
