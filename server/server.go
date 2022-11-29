package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	Addr    string
	Handler http.Handler
	TLS     *tls.Config
	Logger  *zap.Logger
}

func (s *Server) Listen() error {
	s.Logger = s.Logger.Named("server")
	s.Logger.Sugar().Infow("opening tcp socket", "addr", s.Addr)

	var lis net.Listener
	var err error

	if s.TLS != nil {
		lis, err = tls.Listen("tcp", s.Addr, s.TLS)
	} else {
		lis, err = net.Listen("tcp", s.Addr)
	}

	if err != nil {
		return err
	}
	defer lis.Close()

	for {
		conn, err := lis.Accept()
		if err != nil {
			s.Logger.Error("failed to accept connection, closing...")
			conn.Close()
			continue
		}
		now := time.Now()

		if err := conn.SetReadDeadline(now.Add(5 * time.Second)); err != nil {
			s.Logger.Error("failed to set read deadline, closing...")
			conn.Close()
			continue
		}
		if err := conn.SetWriteDeadline(now.Add(30 * time.Second)); err != nil {
			s.Logger.Error("failed to set write deadline, closing...")
			conn.Close()
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(c net.Conn) {
  defer c.Close()
	reader := bufio.NewReader(c)

	r, err := http.ReadRequest(reader)
	if err != nil {
		if err.Error() != "EOF" {
			s.Logger.Sugar().Errorf("failed to read request: %v", err.Error())
		}
		return
	}

	if res := validRequest(r); res != nil {
		if err := res.Write(c); err != nil {
			s.Logger.Sugar().Errorf("failed to write response: %v", err.Error())
			return
		}
		return
	}

	w := NewOSResponseWriter()
	s.Handler.ServeHTTP(w, r)

	res := NewResponse()
	res.StatusCode = w.statusCode
	res.Header = w.header
	res.Body = io.NopCloser(bytes.NewReader(w.body))
	res = SetBody(res, w.body)

	err = res.Write(c)
	if err != nil {
		s.Logger.Sugar().Errorf("failed to write response: %v", err.Error())
		return
	}
}
