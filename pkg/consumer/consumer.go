package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"sshmonitor/config"
	httpserver "sshmonitor/pkg/http"
	"time"
)

type controllerClient struct {
	client         *http.Client
	controllerAddr string
}

func newContorllerClient(controllerAddr string) *controllerClient {
	return &controllerClient{
		controllerAddr: controllerAddr,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:      100,
				IdleConnTimeout:   90 * time.Second,
				DisableKeepAlives: false,
			},
		},
	}
}

var globalClient *controllerClient

func Consume(addr string) {
	globalClient = newContorllerClient(addr)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for {
		select {
		case msg := <-httpserver.Catcher.CommandChan:
			err := globalClient.postController(ctx, msg)
			if err != nil {
				config.GlobalLogger.Error("consumer process error", zap.Error(err))
			}
		default:
		}
	}
}

//请求controller 创建cr

func (c *controllerClient) postController(ctx context.Context, commandMsg *httpserver.SshCommand) error {
	jsonData, err := json.Marshal(commandMsg)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		c.controllerAddr,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	config.GlobalLogger.Debug("http code = ", zap.Int("httpcode", resp.StatusCode))
	config.GlobalLogger.Debug("send request successful.", zap.String("BODY", string(jsonData)))
	return nil
}
