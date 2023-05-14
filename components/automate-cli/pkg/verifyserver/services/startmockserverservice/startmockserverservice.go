package startmockserverservice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/constants"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
)

// MockServerService provides functionality to start mock servers.
type MockServerService struct {
	MockServers []*models.Server
}

type IStartMockServersService interface {
	StartMockServer(cfg *models.StartMockServerRequestBody) error
	GetMockServers() []*models.Server
	SetMockServers(servers []*models.Server)
}

func New() IStartMockServersService {
	return &MockServerService{}
}

// func to get mock server list
func (s *MockServerService) GetMockServers() []*models.Server {
	return s.MockServers
}

func (s *MockServerService) SetMockServers(servers []*models.Server) {
	s.MockServers = servers
}

// StartMockServer starts a mock server of the given type and port.
func (servers *MockServerService) StartMockServer(cfg *models.StartMockServerRequestBody) error {
	var myServer *models.Server
	var err error
	switch cfg.Protocol {
	case constants.TCP:
		myServer, err = servers.StartTCPServer(cfg.Port)
	case constants.UDP:
		myServer, err = servers.StartUDPServer(cfg.Port)
	case constants.HTTPS:
		myServer, err = servers.StartHTTPSServer(cfg.Port, cfg.Cert, cfg.Key)
	default:
		err = errors.New("unsupported protocol")
	}
	if err != nil {
		return err
	}
	servers.MockServers = append(servers.MockServers, myServer)
	return nil
}

func (s *MockServerService) StartTCPServer(port int) (*models.Server, error) {
	// create a TCP listener on the specified port and
	// save the listener instance in the handler struct
	listener, err := net.Listen(constants.TCP, fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return nil, errors.New(err.Error())
	}

	log.Printf("TCP server started on port %d", port)

	//Create a channel to listen for close signal
	signalChan := make(chan bool, 1)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-signalChan:
					log.Println("Stopping accepting incoming connections....")
					return
				default:
					log.Println("Error accepting connection:", err)
				}
			}
			go s.HandleTCPRequest(conn)
			conn.Close()
		}
	}()

	return &models.Server{
		Port:        port,
		ListenerTCP: listener,
		ListenerUDP: nil,
		Protocol:    constants.TCP,
		SignalChan:  signalChan,
	}, nil
}

func (s *MockServerService) HandleTCPRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data from connection:", err)
		return
	}

	// time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Your message is: %v.", string(buf[:]))
	conn.Write([]byte(responseStr))

	// close conn
	conn.Close()

}

func (s *MockServerService) StartUDPServer(port int) (*models.Server, error) {
	addr, err := net.ResolveUDPAddr(constants.UDP, fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP(constants.UDP, addr)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	//Create a channel to listen for close signal
	signalChan := make(chan bool, 1)

	log.Printf("UDP server started on port %d", port)

	go func() {
		for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				select {
				case <-signalChan:
					log.Println("Stopping UDP server....")
					return
				default:
					log.Println("Error receiving message:", err)
				}
			}
			go s.HandleUDPRequest(conn, addr, buf[:n])
		}
	}()

	return &models.Server{
		Port:        port,
		ListenerTCP: nil,
		ListenerUDP: conn,
		Protocol:    constants.UDP,
		SignalChan:  signalChan,
	}, nil
}

func (s *MockServerService) HandleUDPRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	log.Printf("UDP request received from %v", addr)
	responseStr := []byte(fmt.Sprintf("Your message is: %v.", string(buf[:])))
	conn.WriteToUDP(responseStr, addr)
}

func (s *MockServerService) StartHTTPSServer(port int, cert string, key string) (*models.Server, error) {

	// Load the TLS certificate and private key
	tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		fmt.Println("Cert error")
		fmt.Printf("%v", err)
		return nil, err
	}

	// Create the TLS configuration for the server
	config := &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		MinVersion:   tls.VersionTLS13,
	}

	// Create the HTTPS server
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: config,
	}

	// Start the HTTPS server
	go func() {
		err = server.ListenAndServeTLS("", "")
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting HTTPS server: ", err)
			return
		}
	}()

	fmt.Printf("HTTPS server started on port %d\n", port)
	return &models.Server{
		Port:         port,
		ListenerTCP:  nil,
		ListenerUDP:  nil,
		ListenerHTTP: server,
		Protocol:     constants.HTTPS,
	}, nil
}
