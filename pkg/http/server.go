package httpserver

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"time"

	"net/http"
	"sshmonitor/config"
)

type handler func(http.ResponseWriter, *http.Request)

const (
	router       = "/command_log"
	defaultPort  = "8080"
	nodeIPEnvVar = "NODE_IP" //downwardApi
)

type Server struct {
	server  *http.Server
	handler handler
	router  string
	port    string
}

var Catcher SshCommandCatcher

type SshCommandCatcher struct {
	CommandChan chan *SshCommand
}

type SshCommand struct {
	nodeIP      string //nodeIp 应该一开始就初始化 取自downwardApi 环境变量
	ExecuteTime string `json:"time,omitempty"`
	User        string `json:"user,omitempty"`
	ClientIp    string `json:"ip,omitempty"`
	Port        string `json:"port,omitempty"`
	Pwd         string `json:"pwd,omitempty"`
	Command     string `json:"command,omitempty"`
	ExitCode    int32  `json:"exit_code,omitempty"`
}

func newSshCommand() (*SshCommand, error) {
	nodeIP := os.Getenv(nodeIPEnvVar)
	if nodeIP == "" {
		return nil, fmt.Errorf("%s environment variable not set", nodeIPEnvVar)
	}
	return &SshCommand{
		nodeIP: nodeIP,
	}, nil
}

func NewServer() *Server {

	return &Server{
		server: &http.Server{
			ReadHeaderTimeout: 10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       90 * time.Second,
		},
		handler: serverHandler,
		router:  router,
		port:    defaultPort,
	}

}

func (s *Server) StartServer(chanLength int) {
	//全局channel 之后从这个channel里取出数据发送到controller
	InitSshCommandCatcher(chanLength)

	http.HandleFunc(s.router, s.handler)
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
}

func InitSshCommandCatcher(chanLength int) {
	Catcher.CommandChan = make(chan *SshCommand, chanLength)
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	sshCmd, err := newSshCommand()
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		config.GlobalLogger.Error("", zap.Error(err))
		return
	}
	defer r.Body.Close()
	//反序列化逻辑没写
	config.GlobalLogger.Debug("", zap.ByteString("body", body))
	//将body转成json格式
	//1.反序列化body到sshcmd,得到完整的nodeip 和 其他数据
	//2.发送到全局变量的channel中

}

// 将收到的cmd转换成json之后才能反序列化
func stringToJson(string) (string error) {
	return nil
}
