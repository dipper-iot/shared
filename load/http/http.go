package http

import (
	"fmt"
	"gitlab.com/dipper-iot/shared/cli"
	"gitlab.com/dipper-iot/shared/logger"
	"gitlab.com/dipper-iot/shared/service"
	"net"
	base "net/http"
)

type ServerHttp struct {
	post       string
	listen     net.Listener
	httpServer *base.Server
	mux        *base.ServeMux
}

func (s ServerHttp) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:         "http-post",
			Env:          []string{"HTTP_PORT"},
			Usage:        "Port start http server",
			DefaultValue: "8080",
		},
	}
}

func NewServerHttpPort(post string) *ServerHttp {
	server := &ServerHttp{
		post: post,
		mux:  base.NewServeMux(),
	}

	return server
}

func NewServerHttp() *ServerHttp {
	server := &ServerHttp{
		post: "",
		mux:  base.NewServeMux(),
	}

	return server
}

func (s ServerHttp) Name() string {
	return "server-http"
}

func (s ServerHttp) Priority() int {
	return 1
}

func (s *ServerHttp) Start(o *service.Options, c *cli.Context) error {
	if c.IsHelp() {
		return nil
	}
	var err error
	if s.post == "" {
		s.post = c.String("http-post")
	}
	s.listen, err = net.Listen("tcp", fmt.Sprintf(":%s", s.post))
	if err != nil {
		return err
	}

	s.httpServer = &base.Server{
		Handler: s.mux,
	}

	logger.Infof("HTTP Server listening at %v", s.listen.Addr())
	go func() {
		if err := s.httpServer.Serve(s.listen); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}()

	return nil
}

func (s *ServerHttp) Stop() error {
	if s.httpServer != nil {
		err := s.httpServer.Close()
		if err != nil {
			return err
		}
	}

	if s.listen != nil {
		err := s.listen.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ServerHttp) HandleFunc(path string, handler base.Handler) {
	s.mux.Handle(path, handler)
}
