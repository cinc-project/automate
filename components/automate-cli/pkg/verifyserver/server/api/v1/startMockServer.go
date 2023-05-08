package v1

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/gofiber/fiber"
)

func (h *Handler) StartMockServer(c *fiber.Ctx) {
	// Specify the server type and port
	// parse the port number from the request body
	// type RequestBody struct {
	// 	Port     int    `json:"port"`
	// 	Protocol string `json:"protocol"`
	// }
	// var reqBody RequestBody
	// err := c.BodyParser(&reqBody)
	// if err != nil {
	// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
	// 	return
	// }

	// if reqBody.Port < 0 || reqBody.Port > 65535 {
	// 	c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid port number"})
	// 	return
	// }

	// portInt := reqBody.Port
	// // serverType := "TCP"

	// if h.listener != nil {
	// 	// Server is already running
	// 	c.JSON(fiber.Map{"port": portInt, "message": fmt.Sprintf("Server is already running on port %d", portInt), "list": h.listener})
	// 	return
	// }

	var (
		HOST = "localhost"
		PORT = c.Params("port")
		TYPE = "tcp"
	)

	// fmt.Println(portInt)
	// Start TCP server
	// addr := net.TCPAddr{
	// 	IP:   net.IPv4(127, 0, 0, 1),
	// 	Port: portInt,
	// }

	// create a TCP listener on the specified port and
	// save the listener instance in the handler struct

	// fmt.Printf("\n%s:%d\n", HOST, PORT)

	// _, err = net.ListenTCP(TYPE, &addr)
	// listen, err := net.Listen(TYPE, fmt.Sprintf("%s:%d", HOST, PORT))
	fmt.Println(PORT)
	ln, err := net.Listen(TYPE, HOST+":"+PORT)
	fmt.Printf("My error: %v \n", err)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// close listen
	// fmt.Printf("%v", listen)
	// defer listen.Close()
	fmt.Println("Also Got here")
	fmt.Printf("%v", ln)
	// defer h.listener.Close()
	// go func() {
	for {
		fmt.Println("Also Got here-1")
		conn, err := ln.Accept()
		fmt.Printf("%v", err)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		go handleRequest(conn)
		conn.Close()
	}
	// }()
}

// func (h *Handler) StopMockServer(c *fiber.Ctx) {
// 	if h.listener != nil {
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
