package v1

import (
	"fmt"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/startmockserverservice"
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

	for _, s := range h.MockServerServices {
		if s.Port == reqBody.Port {
			// Server is already running
			c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"port": reqBody.Port, "message": fmt.Sprintf("Server is already running on port %d", reqBody.Port), "listener": s.ListenerTCP})
			return
		}
	}

	service := startmockserverservice.MockServerService{}
	mockServer, err := service.StartMockServer(*reqBody)

	fmt.Printf("Error: %v", err)
	h.MockServerServices = append(h.MockServerServices, mockServer)
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
