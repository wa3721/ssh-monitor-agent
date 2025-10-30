package httpserver

import (
	"bytes"
	"context"
	"go.uber.org/zap"
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"net/http"
	"sshmonitor/config"
	"time"
)

type Consumer struct {
	controllerClient *http.Client
}

func NewConsumer() *Consumer {
	transport := &http.Transport{
		DisableKeepAlives: false,

		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		MaxConnsPerHost:     20,

		IdleConnTimeout:       90 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
	}

	return &Consumer{
		controllerClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
	}
}

func (c *Consumer) Consume(ctx context.Context, controllerAddr string) {
	for {
		select {
		case <-ctx.Done():
			config.GlobalLogger.Info(ctx.Err().Error())
			return
		default:
			msg := <-Catcher.CommandChan
			jsonData, err := json.Marshal(msg)
			if err != nil {
				config.GlobalLogger.Error("send to controller error:", zap.Error(err))
			}
			resp, err := c.controllerClient.Post("http://"+controllerAddr, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				config.GlobalLogger.Error("send to controller error:", zap.Error(err))
			}
			all, err := io.ReadAll(resp.Body)
			if err != nil {
				config.GlobalLogger.Error("send to controller error:", zap.Error(err))
			}

			config.GlobalLogger.Debug("controller response:", zap.Any("resp", string(all)))
			resp.Body.Close()
		}
	}

}
