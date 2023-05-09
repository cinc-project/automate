package v1

import (
	"net"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/statusservice"
)

type Server struct {
	port        int
	listenerTCP net.Listener
	listenerUDP net.PacketConn
	protocol    string
}

type Handler struct {
	Logger        logger.ILogger
	StatusService statusservice.IStatusService
	servers       []*Server
}

func NewHandler(logger logger.ILogger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) AddStatusService(ss statusservice.IStatusService) *Handler {
	h.StatusService = ss
	return h
}
