package stopmockserverservice

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/lib/logger"
)

type IStopMockServerService interface {
	StopMockServer(server *models.Server) error
	StopTCPServer(server *models.Server) error
	StopUDPServer(server *models.Server) error
	StopHTTPSServer(server *models.Server) error
}

// StopMockServerService provides functionality to stop mock servers.
type StopMockServerService struct {
	Logger logger.Logger
}

func NewStopMockServerService(logger logger.Logger) IStopMockServerService {
	return &StopMockServerService{
		Logger: logger,
	}
}

// StopMockServer stops a mock server of the given type and port.
func (sm *StopMockServerService) StopMockServer(server *models.Server) error {
	switch server.Protocol {
	case "tcp":
		return sm.StopTCPServer(server)
	case "udp":
		return sm.StopUDPServer(server)
	case "https":
		return sm.StopHTTPSServer(server)
	}
	return nil
}

func (sm *StopMockServerService) StopTCPServer(server *models.Server) error {

	server.SignalChan <- true
	if err := server.ListenerTCP.Close(); err != nil {
		return err
	}

	return nil
}

func (sm *StopMockServerService) StopUDPServer(server *models.Server) error {

	server.SignalChan <- true
	if err := server.ListenerUDP.Close(); err != nil {
		return err
	}

	return nil
}

func (sm *StopMockServerService) StopHTTPSServer(server *models.Server) error {

	//server.SignalChan <- true
	if err := server.ListenerHTTP.Close(); err != nil {
		return err
	}

	return nil
}
