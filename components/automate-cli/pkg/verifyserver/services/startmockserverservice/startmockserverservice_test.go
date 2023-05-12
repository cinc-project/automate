package startmockserverservice_test

import (
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/startmockserverservice"
	"github.com/stretchr/testify/require"
)

const (
	SERVER_KEY  = "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAr4V0UWr7B4Hu7gTRDYt7FGEBFXC6V39Cf8i4xgtUYZOKkvt/\nMXiYegFWgZPLah0TyhgdyjLpC8FKdpq+FHpksGgfkR8ARyBsp50b4FRfjCSpKNj+\nSpUhHzovYSywLJFykT9Wq00n/9M992Bxn94MJQIfmHfWUrpOl66o0fdv8viblik7\nEh+sPl+Lgg6miTg8ev529e0Lo//LrW5AtHHEEeBvBPWp1kwnWgatECYDZgXQAbvI\nYzNkqpYyE853hbCBzZ1HKsz4WK0poGg3MSl/AbJC36xi1ecP3rCnh5ae123p1V9d\n337YUDxL7Wst9wKulj+uVPgxNchg11JqgNNaawIDAQABAoIBAEV4SbCL6i1vhPTa\nHTACO8W2Gyq0QlytNtHCzTc9drlkHx3Lwuz+sULg0q9YotMuDQ4Y+3lzKwAHEgTd\nfEw4oS+dFplmrsJ4F+lDaqwgWOzr+bP8JrG4UrK8YdJRUK1jJ/hLHG+SizlbD5Sg\nrHg894uSSpUbIU3/BWpNq+3mxH1itofFanP7ePTCV78kpfGHR2Ok2UMHxyQHSk7y\nzii78grmbQObndx2b7bC7cswCa3/pjMN3X6B/Jjy/gn3b50RR77YFqpp7zzrnrZA\nSgx6CrJSBRwDExhCzMPBh3UjtYgohPHlVLjPluUSZheYDZQB1Gr8HwN3iL1dmf1w\n8QYQnCkCgYEA4GgAReXoOzjbjo06lr8TGogavWB2D3Zls3WUGO9FcD/59IJTi1rC\nh13V7gZYRFdc7HH660vqtfqsYtt5akV0FxYfXAMla9fuYxmr7h8OUfs8jTp77Uq8\ntSlQ3VPaas4YTzTReTD+wHvKVQ6uiDAYqjTqF9yJNkMcZ8veNzWtTG8CgYEAyDuN\nnNGY6VD2p5UA24y01yZqNH9WYAwuIYNAqCypZLzJBP8ABVrCOUNuPIaBCxmU8Orq\nYN9FVexY3nBdkDEjEwbXvnYyydoadIN9DMB3YhQ9etf3DuZFI65+DAskHE/cJrpP\neUM/XPpC73pf3U+eqSZwV+7FX8idgsK6Rxnnh8UCgYAyTN2SzV/qtmnwYBO76oR7\ns/pabJ7KBH3zZe2WUTu9V3nNptDXMbbc5NmpCt8KIpL/pOTbjR7FP7UYS53BhmPp\nMNpCo6nlrHcQ25ZAP9HT6n6+IVfZ7qCx8trfYYZZ3mxwhKRXh/Xya00FF89jU3ST\n4lx+kL5o3U4mrfnXYj7AHQKBgC5vERIS0SEaM3j9ZuuDH9TdBbgS55byfCgtZesa\nIFZKKVvNPtX/DBd3ebLzhi1qy01rTNsWK+AXJSzAZhIwMvAQoCt9AZ4pxATNEUzJ\nvWWzR+aa+qIr6FC0AGsOkls2cdlRT2jRnXoUVz1t5ZlPA346eccKih8CSPSv777Z\nVQX5AoGAf4rHCoMGGqlJZxaroCS2AgDtzq5Za9XdwT3RczYnJGCkv7nK4K0e4zxB\nqMxAiUTrzae1Hpf5OV0GXeoiN0taZ9L2Cc56OlaqPRL2BDpYU2K+3okrftXLgFCX\nClDVL9XxSbVT99AUEBw9eH9hYFVKjHIEyJ20Udw0zsR+qU4R6/0=\n-----END RSA PRIVATE KEY-----"
	SERVER_CERT = "-----BEGIN CERTIFICATE-----\nMIIDTTCCAjWgAwIBAgIUNidKDNanRMXILhrf1//SuK7DC3UwDQYJKoZIhvcNAQEL\nBQAwYzELMAkGA1UEBhMCVVMxEzARBgNVBAgMCldhc2hpbmd0b24xEDAOBgNVBAcM\nB1NlYXR0bGUxGjAYBgNVBAoMEUNoZWYgU29mdHdhcmUgSW5jMREwDwYDVQQDDAhw\ncm9ncmVzczAeFw0yMzA1MDkxMjAwMjZaFw0yNjA1MDgxMjAwMjZaMDYxCzAJBgNV\nBAYTAlVTMRMwEQYDVQQIDApDYWxpZm9ybmlhMRIwEAYDVQQDDAlsb2NhbGhvc3Qw\nggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCvhXRRavsHge7uBNENi3sU\nYQEVcLpXf0J/yLjGC1Rhk4qS+38xeJh6AVaBk8tqHRPKGB3KMukLwUp2mr4UemSw\naB+RHwBHIGynnRvgVF+MJKko2P5KlSEfOi9hLLAskXKRP1arTSf/0z33YHGf3gwl\nAh+Yd9ZSuk6XrqjR92/y+JuWKTsSH6w+X4uCDqaJODx6/nb17Quj/8utbkC0ccQR\n4G8E9anWTCdaBq0QJgNmBdABu8hjM2SqljITzneFsIHNnUcqzPhYrSmgaDcxKX8B\nskLfrGLV5w/esKeHlp7XbenVX13ffthQPEvtay33Aq6WP65U+DE1yGDXUmqA01pr\nAgMBAAGjJjAkMCIGA1UdEQQbMBmCCWxvY2FsaG9zdIIMaHR0cHMtc2VydmVyMA0G\nCSqGSIb3DQEBCwUAA4IBAQCxXspmv+BCRVFykb80embuZGCXMh7YmH3j5dJZhaKL\n/PPcUjgJTYRanDSSwt5IFyyYwiYG9qdUrRLxR5pgpdj0vNRaLdabG24UsVQQuK1Y\nrxb/6HF2AwWASiS5YLKoVMwg1sYiskpA7gJ23Xe34BVckqAd+Yoss4zDNR3d/yRM\nQYbnr/STpW4c+7jHL+vlpu/OdHwEtsTNrUG6lk0YO1lGH43a0rmvMzYOCgZEfr3h\nK4zbwy053Vq8PGIH24/bNu67pSfslgGI30bN9PmUdFMwEFuRC7rCFgQ8LpRbJNf1\nf5CUhHWn2nC05XOzKm+Kj/NHPtw5iJkrQvLNtsdiO92O\n-----END CERTIFICATE-----"
)

func TestStartMockServer(t *testing.T) {

	t.Run("Start TCP server", func(t *testing.T) {
		servers := startmockserverservice.New()
		cfg := &models.StartMockServerRequestBody{
			Protocol: "tcp",
			Port:     8000,
		}

		err := servers.StartMockServer(cfg)
		server := servers.GetMockServers()[0]
		require.NoError(t, err)
		require.NotNil(t, server)

		// Check that the server is listening on the correct port
		require.Equal(t, cfg.Protocol, server.ListenerTCP.Addr().Network())
		require.Equal(t, cfg.Port, server.ListenerTCP.Addr().(*net.TCPAddr).Port)

		err = servers.StartMockServer(cfg)
		server2 := servers.GetMockServers()[1]
		require.Error(t, err)
		require.Nil(t, server2)
		// fmt.Print()
		// Stop the server and check that it was closed correctly
		err = server.ListenerTCP.Close()
		require.NoError(t, err)
		if server2 != nil {
			err = server2.ListenerTCP.Close()
			require.Error(t, err)
		}
	})

	t.Run("Start UDP server", func(t *testing.T) {
		servers := startmockserverservice.New()
		cfg := &models.StartMockServerRequestBody{
			Protocol: "udp",
			Port:     8001,
		}
		err := servers.StartMockServer(cfg)

		server := servers.GetMockServers()[0]
		require.NoError(t, err)
		require.NotNil(t, server)

		// Check that the server is listening on the correct port
		require.Equal(t, cfg.Protocol, server.ListenerUDP.LocalAddr().Network())
		require.Equal(t, cfg.Port, server.ListenerUDP.LocalAddr().(*net.UDPAddr).Port)

		fmt.Println(server.ListenerUDP.LocalAddr())
		// Stop the server and check that it was closed correctly
		err = server.ListenerUDP.Close()
		require.NoError(t, err)
	})

	t.Run("Start HTTPS server", func(t *testing.T) {
		servers := startmockserverservice.New()
		cfg := &models.StartMockServerRequestBody{
			Protocol: "https",
			Port:     8002,
			Cert:     SERVER_CERT,
			Key:      SERVER_KEY,
		}
		err := servers.StartMockServer(cfg)
		server := servers.GetMockServers()[0]
		require.NoError(t, err)
		require.NotNil(t, server)

		// Check that the server is listening on the correct port
		require.Equal(t, fmt.Sprintf(":%d", cfg.Port), server.ListenerHTTP.Addr)
		// require.Equal(t, cfg.Protocol, server.ListenerHTTP.Addr)

		// Stop the server and check that it was closed correctly
		err = server.ListenerHTTP.Close()
		require.NoError(t, err)
	})

	t.Run("Unsupported protocol", func(t *testing.T) {
		servers := startmockserverservice.New()
		err := servers.StartMockServer(&models.StartMockServerRequestBody{
			Port:     8003,
			Protocol: "http",
			Cert:     "",
			Key:      "",
		})

		require.Error(t, err)
		require.Equal(t, "unsupported protocol", err.Error())
	})
}

func TestHandleTCPRequest(t *testing.T) {
	// create a new StartMockServerService
	service := &startmockserverservice.MockServerService{}

	t.Run("For healthy connetion", func(t *testing.T) {
		// create a new TCP listener
		listener, err := net.Listen("tcp", "localhost:8080")
		if err != nil {
			t.Errorf("Error creating TCP listener: %v", err)
			return
		}

		go func() {
			for {
				conn, err := listener.Accept()
				if err != nil {
					t.Errorf("Error Listerning TCP connection: %v", err)
					return
				}
				go service.HandleTCPRequest(conn)
				// conn.Close()
			}
		}()

		// create a new TCP connection
		conn1, err := net.Dial("tcp", listener.Addr().String())
		if err != nil {
			t.Errorf("Error dialing TCP connection: %v", err)
			return
		}

		// write a message to the connection
		message := "test message"
		_, err = conn1.Write([]byte(message))
		if err != nil {
			t.Errorf("Error writing message to connection: %v", err)
			return
		}

		// read the response from the connection
		responseBuf := make([]byte, 1024)
		_, err = conn1.Read(responseBuf)
		if err != nil {
			t.Errorf("Error reading response from connection: %v", err)
			return
		}
		response := string(responseBuf[:])

		// verify that the response contains the message
		if !strings.Contains(response, message) {
			t.Errorf("Unexpected response. Expected message should contain \"%v\" \nActual message received: %v\n", message, response)
		}

		// close the connection and #listener
		defer conn1.Close()
		// defer listener.Close()
	})

	// t.Run("For unhealthy connetion", func(t *testing.T) {
	// 	// create a new TCP listener
	// 	listener, err := net.Listen("tcp", "localhost:8080")
	// 	if err != nil {
	// 		t.Errorf("Error creating TCP listener: %v", err)
	// 		return
	// 	}

	// 	conn, err := listener.Accept()
	// 	go func() {
	// 		// for {
	// 		if err != nil {
	// 			t.Errorf("Error Listerning TCP connection: %v", err)
	// 			return
	// 		}
	// 		service.HandleTCPRequest(conn)
	// 		conn.Close()
	// 		// }
	// 	}()

	// 	// service.HandleTCPRequest(wConn)

	// 	// assert that the connection is closed
	// 	if checkConnectionClosed(conn) {
	// 		fmt.Printf("%v", checkConnectionClosed(conn))
	// 		t.Errorf("Expected connection to be closed, but it was not")
	// 	}

	// })

}

func checkConnectionClosed(conn net.Conn) bool {
	// Try to read 0 bytes from the connection
	buf := make([]byte, 0)
	_, err := conn.Read(buf)

	if err != nil {
		// If the error is "use of closed network connection",
		// the connection is closed.
		if strings.Contains(err.Error(), "use of closed network connection") {
			return true
		}
		// Otherwise, there was a different error reading from the
		// connection, so assume it is still open.
		return false
	}

	// If there was no error, assume the connection is still open.
	return false
}

// func TestHandleUDPRequest(t *testing.T) {
// 	// create a new StartMockServerService
// 	service := &startmockserverservice.StartMockServerService{}

// 	// create a new UDP listener
// 	addr, _ := net.ResolveUDPAddr("udp", ":8080")
// 	udpServer, err := net.ListenUDP("udp", addr)
// 	if err != nil {
// 		t.Errorf("Error creating UDP listener: %v", err)
// 		return
// 	}

// 	go func() {
// 		for {
// 			buf := make([]byte, 1024)
// 			n, addr1, err := udpServer.ReadFromUDP(buf)

// 			if err != nil {
// 				t.Errorf("Error Listerning UDP connection: %v", err)
// 				return
// 			}
// 			go service.HandleUDPRequest(udpServer, addr1, buf[:n])
// 			// conn.Close()
// 		}
// 	}()

// 	// create a new UDP connection
// 	conn1, err := net.DialUDP("udp", nil, addr)
// 	if err != nil {
// 		t.Errorf("Error dialing UDP connection: %v", err)
// 		return
// 	}

// 	// write a message to the connection
// 	message := "test message"
// 	_, err = conn1.Write([]byte(message))
// 	if err != nil {
// 		t.Errorf("Error writing message to connection: %v", err)
// 		return
// 	}

// 	// read the response from the connection
// 	responseBuf := make([]byte, 1024)
// 	_, err = conn1.Read(responseBuf)
// 	if err != nil {
// 		t.Errorf("Error reading response from connection: %v", err)
// 		return
// 	}
// 	response := string(responseBuf[:])

// 	// verify that the response contains the message
// 	if !strings.Contains(response, message) {
// 		t.Errorf("Unexpected response. Expected message should contain \"%v\" \nActual message received: %v\n", message, response)
// 	}

// 	// close the connection and #listener
// 	defer conn1.Close()
// 	// defer listener.Close()
// }
