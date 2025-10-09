package httpserver

import (
	"go.uber.org/zap"
	"io"
	"kubeauth/config"
	"net/http"
)

type handler func(http.ResponseWriter, *http.Request)

const (
	router      = "/command_log"
	defaultPort = "8080"
)

type Server struct {
	server  *http.Server
	handler handler
	router  string
	port    string
}

func NewServer() *Server {

	return &Server{
		server: &http.Server{
			ReadHeaderTimeout: 10,
			WriteTimeout:      10,
		},
		handler: func(w http.ResponseWriter, r *http.Request) {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				config.GlobalLogger.Error("", zap.Error(err))
				return
			}
			defer r.Body.Close()
			//
			config.GlobalLogger.Info("", zap.ByteString("body", body))
		},
		router: router,
		port:   defaultPort,
	}

}

func (s *Server) StartServer() {
	http.HandleFunc(s.router, s.handler)
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
}
