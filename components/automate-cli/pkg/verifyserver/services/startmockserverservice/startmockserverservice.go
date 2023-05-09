package startmockserverservice

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
)

// StartMockServerService provides functionality to start mock servers.
type StartMockServerService struct{}

// StartMockServer starts a mock server of the given type and port.
func (servers *StartMockServerService) StartMockServer(cfg models.StartMockServerRequestBody) (*models.Server, error) {
	switch cfg.Protocol {
	case "tcp":
		return servers.startTCPServer(cfg.Port)
	case "udp":
		return servers.startUDPServer(cfg.Port)
	case "https":
		return servers.startHTTPSServer(cfg.Port, cfg.Cert, cfg.Key)
	default:
		return nil, errors.New("unsupported protocol")
	}
}

func (s *StartMockServerService) startTCPServer(port int) (*models.Server, error) {
	// create a TCP listener on the specified port and
	// save the listener instance in the handler struct
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return nil, errors.New(err.Error())
	}
	// defer listener.Close()

	log.Printf("TCP server started on port %d", port)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println("Error accepting connection:", err)
				continue
			}
			go s.handleTCPRequest(conn)
			conn.Close()
		}
	}()

	return &models.Server{
		Port:        port,
		ListenerTCP: listener,
		ListenerUDP: nil,
		Protocol:    "tcp",
	}, nil
}

func (s *StartMockServerService) handleTCPRequest(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	fmt.Printf("%v", err)
	if err != nil {
		fmt.Println("Error reading data from connection:", err)
		return
	}

	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Your message is: %v. Received time: %v", string(buf[:]), time)
	conn.Write([]byte(responseStr))

	// close conn
	conn.Close()

}

func (s *StartMockServerService) startUDPServer(port int) (*models.Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	log.Printf("UDP server started on port %d", port)

	go func() {
		for {
			buf := make([]byte, 1024)
			n, addr, err := conn.ReadFromUDP(buf)
			if err != nil {
				log.Println("Error receiving message:", err)
				continue
			}
			go s.handleUDPRequest(conn, addr, buf[:n])
		}
	}()

	return &models.Server{
		Port:        port,
		ListenerTCP: nil,
		ListenerUDP: conn,
		Protocol:    "tcp",
	}, nil
}

func (s *StartMockServerService) handleUDPRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) {
	log.Printf("UDP request received from %v", addr)
	response := []byte("UDP response")
	conn.WriteToUDP(response, addr)
}

func (s *StartMockServerService) startHTTPSServer(port int, cert string, key string) (*models.Server, error) {

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
	}

	// Create the HTTPS server
	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		TLSConfig: config,
	}

	// Start the HTTPS server
	go func() {
		// err := http.ListenAndServeTLS(fmt.Sprintf(":%d", port), cert, key, nil)
		err = server.ListenAndServeTLS("", "")
		fmt.Println("Serv error")
		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting HTTPS server: ", err)
		}
	}()

	fmt.Printf("HTTPS server started on port %d\n", port)
	return &models.Server{
		Port:         port,
		ListenerTCP:  nil,
		ListenerUDP:  nil,
		ListenerHTTP: server,
		Protocol:     "tcp",
	}, nil
}
