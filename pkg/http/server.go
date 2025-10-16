package httpserver

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"os"
	"sync"
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

var pool sync.Pool

type SshCommandCatcher struct {
	CommandChan chan *SshCommand
}

type SshCommand struct {
	nodeIP      string //nodeIp 应该一开始就初始化 取自downwardApi 环境变量
	ExecuteTime string `json:"Time,omitempty"`
	User        string `json:"User,omitempty"`
	ClientIp    string `json:"IP,omitempty"`
	Port        string `json:"Port,omitempty"`
	Pwd         string `json:"PWD,omitempty"`
	Command     string `json:"Command,omitempty"`
	ExitCode    string `json:"ExitCode,omitempty"`
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
	initSshCommandCatcher(chanLength)
	//初始化sshCmd 对象池,对象复用使用
	newSshCommandPool()

	http.HandleFunc(s.router, s.handler)
	err := http.ListenAndServe(":"+s.port, nil)
	if err != nil {
		config.GlobalLogger.Panic(err.Error())
		return
	}
}

func initSshCommandCatcher(chanLength int) {
	Catcher.CommandChan = make(chan *SshCommand, chanLength)
}

func newSshCommandPool() {
	pool = sync.Pool{New: func() interface{} {
		return newSshCommand()
	}}
}
func newSshCommand() *SshCommand {
	nodeIP := os.Getenv(nodeIPEnvVar)
	if nodeIP == "" {
		nodeIP = "Unknown"
	}
	return &SshCommand{
		nodeIP: nodeIP,
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	sshCmd := pool.Get()
	defer pool.Put(sshCmd)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		config.GlobalLogger.Error("", zap.Error(err))
		return
	}
	defer r.Body.Close()
	config.GlobalLogger.Debug("", zap.ByteString("body", body))
	//1.反序列化body到sshcmd,得到完整的nodeip 和 其他数据
	err = json.Unmarshal(body, &sshCmd)
	if err != nil {
		config.GlobalLogger.Error("", zap.Error(err))
		return
	}
	//2.发送到全局变量的channel中
	switch sshCmd.(type) {
	case *SshCommand:
		Catcher.CommandChan <- sshCmd.(*SshCommand)
		if Catcher.CommandChan == nil {
			config.GlobalLogger.Error("CommandChan is nil, cannot send message")
		}
		config.GlobalLogger.Debug("messages send success", zap.Any("", sshCmd.(*SshCommand)))
		return
	default:
		config.GlobalLogger.Error("", zap.Error(fmt.Errorf("type not support")))
		return
	}

}
