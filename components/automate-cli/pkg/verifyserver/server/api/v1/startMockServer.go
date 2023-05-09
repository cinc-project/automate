package v1

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/gofiber/fiber"
)

func (h *Handler) StartMockServer(c *fiber.Ctx) {
	// Specify the server type and port
	// parse the port number from the request body

	reqBody := new(models.StartMockServerRequestBody)
	err := c.BodyParser(&reqBody)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
		return
	}

	if reqBody.Port < 0 || reqBody.Port > 65535 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid port number"})
		return
	}

	portInt := reqBody.Port

	for _, s := range h.servers {
		if s.port == portInt {
			// Server is already running
			c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"port": portInt, "message": fmt.Sprintf("Server is already running on port %d", portInt), "listener": s.listenerTCP})
			return
		}
	}

	var (
		HOST = "localhost"
	)

	// create a TCP listener on the specified port and
	// save the listener instance in the handler struct
	var ln net.Listener
	var lnUDP net.PacketConn
	var s *Server

	switch reqBody.Protocol {
	case "tcp":
		ln, err = net.Listen(reqBody.Protocol, fmt.Sprintf("%s:%d", HOST, portInt))

		s = &Server{port: portInt, listenerTCP: ln, protocol: reqBody.Protocol}
	case "udp":
		lnUDP, err = net.ListenPacket("udp", ":"+strconv.Itoa(portInt))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
			return
		}
		// fmt.Print("addr: ")
		// fmt.Print(reqBody.Protocol)
		// fmt.Print(addr)
		// lnUDP, err := net.ListenUDP("udp", addr)
		s = &Server{port: portInt, listenerUDP: lnUDP, protocol: reqBody.Protocol}
		if err != nil {
			fmt.Printf("Error: %v", err)
		}
		// defer lnUDP.Close()
	}

	if err != nil {
		// Port is already consumed by external Server
		c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"port": portInt, "message": fmt.Sprintf("Port %d is already consumed by other/external Server", portInt), "error": err})
	}

	h.servers = append(h.servers, s)
	fmt.Printf("HI Arvi: %v", h.servers)
	switch reqBody.Protocol {
	case "tcp":
		go func() {
			for {
				conn, err := ln.Accept()
				// fmt.Printf("%v", err)
				if err != nil {
					log.Fatal(err)
					os.Exit(1)
				}
				go handleRequest(conn)
				conn.Close()
			}
		}()
	case "udp":
		go func() {
			var wg sync.WaitGroup
			fmt.Printf("UDP Socket: %v", lnUDP)
			for {
				buffer := make([]byte, 1024)
				_, addr, err := lnUDP.ReadFrom(buffer)
				wg.Add(1)
				if err != nil {
					fmt.Println("Error receiving message: ", err.Error())
					continue
				}
				// fmt.Println("Message received from ", addr.String(), ": ", string(buffer[:n]))
				go handleUDPRequest(&wg, lnUDP, addr, buffer)
			}
			wg.Wait()
		}()

		fmt.Printf("UDP server started on port %d\n", portInt)
		// return nil

	}
}

// func (h *Handler) StopMockServer(c *fiber.Ctx) {
// 	if server != nil {
// 		// close the listener and stop the server
// 		err := h.listener.Close()
// 		if err != nil {
// 			c.Status(500).SendString("Failed to stop TCP server")
// 			return
// 		}

// 		// reset the listener instance in the handler struct
// 		h.listener = nil
// 	}

// 	// return success response
// 	c.Status(200).SendString("TCP server stopped")
// }

func handleRequest(conn net.Conn) {
	// defer conn.Close()

	// Read data from connection
	fmt.Println("I am here")
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	fmt.Printf("%v", err)
	if err != nil {
		fmt.Println("Error reading data from connection:", err)
		return
	}

	// Echo data back to client
	// _, err = conn.Write(buf[:n])
	// if err != nil {
	// 	fmt.Println("Error writing data to connection:", err)
	// }
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("Your message is: %v. Received time: %v", string(buf[:]), time)
	conn.Write([]byte(responseStr))

	// close conn
	conn.Close()

}

func handleUDPRequest(wg *sync.WaitGroup, udpServer net.PacketConn, addr net.Addr, buf []byte) {
	defer wg.Done()
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))

	udpServer.WriteTo([]byte(responseStr), addr)
}
