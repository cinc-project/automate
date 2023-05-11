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
		errString := fmt.Sprintf("start mock-server request body parsing failed: %v", err.Error())
		h.Logger.Error(fmt.Errorf(errString))
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body"})
		return
	}

	// If port number is invalid
	if reqBody.Port < 0 || reqBody.Port > 65535 {
		errString := fmt.Sprintf("start mock-server request body contains invalid port number: %v", err.Error())
		h.Logger.Error(fmt.Errorf(errString))
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid port number"})
		return
	}

	for _, s := range h.MockServerServices {
		if s.Port == reqBody.Port {
			// Server is already running in the port
			errString := fmt.Sprintf("start mock-server request body contains unavailable port: %v", reqBody.Port)
			h.Logger.Error(fmt.Errorf(errString))
			c.Status(fiber.ErrConflict.Code).JSON(fiber.Map{"port": reqBody.Port, "message": fmt.Sprintf(`"%s" server is already running on port %d`, s.Protocol, reqBody.Port)})
			return
		}
	}

	service := startmockserverservice.StartMockServerService{}
	mockServer, err := service.StartMockServer(*reqBody)

	if err != nil {
		errString := fmt.Sprintf("start mock-server request body contains unsupported protocol: %v", err.Error())
		h.Logger.Error(fmt.Errorf(errString))
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "UnSupported Protocol"})
		return
	}
	h.MockServerServices = append(h.MockServerServices, mockServer)
}
