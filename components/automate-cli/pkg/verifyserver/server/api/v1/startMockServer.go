package v1

import (
	"fmt"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/startmockserverservice"
	"github.com/gofiber/fiber"
)

func (h *Handler) StartMockServer(c *fiber.Ctx) {

	reqBody := new(models.StartMockServerRequestBody)
	err := c.BodyParser(&reqBody)

	// If request body is invalid
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
		return
	}

	// If port number is invalid
	if reqBody.Port < 0 || reqBody.Port > 65535 {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid port number"})
		return
	}

	for _, s := range h.MockServerServices {
		if s.Port == reqBody.Port {
			// Server is already running in the port
			c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"port": reqBody.Port, "message": fmt.Sprintf(`"%s" server is already running on port %d`, s.Protocol, reqBody.Port), "listener": s.ListenerTCP})
			return
		}
	}

	service := startmockserverservice.StartMockServerService{}
	mockServer, err := service.StartMockServer(*reqBody)

	fmt.Printf("Error: %v", err)
	h.MockServerServices = append(h.MockServerServices, mockServer)
}
