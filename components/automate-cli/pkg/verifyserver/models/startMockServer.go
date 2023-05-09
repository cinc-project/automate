package models

import (
	"net"
	"net/http"
)

// StartMockServerRequestBody contains the configuration for starting a mock server.
type StartMockServerRequestBody struct {
	Port     int
	Protocol string
	Cert     string
	Key      string
}

// NewStartMockServerRequestBody creates a new StartMockServerRequestBody object.
func NewStartMockServerRequestBody(port int, protocol string, cert string, key string) *StartMockServerRequestBody {
	return &StartMockServerRequestBody{
		Port:     port,
		Protocol: protocol,
		Cert:     cert,
		Key:      key,
	}
}

type Server struct {
	Port         int
	ListenerTCP  net.Listener
	ListenerUDP  net.PacketConn
	ListenerHTTP *http.Server
	Protocol     string
}
