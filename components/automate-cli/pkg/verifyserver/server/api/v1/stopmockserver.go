package v1

import (
	"fmt"
	"net/http"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) StopMockServer(c *fiber.Ctx) error {

	reqBody := new(models.StopMockServerRequestBody)
	err := c.BodyParser(&reqBody)

	// If request body is invalid
	if err != nil {
		h.Logger.Error("Stop mock server request body parsing failed: ", err.Error())
		return fiber.NewError(http.StatusBadRequest, "Invalid Body Request")
	}

	// If port number is invalid
	if !isValidPort(reqBody.Port) {
		h.Logger.Error("Start mock-server request body contains invalid port number")
		return fiber.NewError(http.StatusBadRequest, "Invalid port number")
	}

	if h.MockServersService.IsMockServerRunningOnGivenPortAndProctocol(reqBody.Port, reqBody.Protocol) {
		if err := h.MockServersService.Stop(reqBody); err != nil {
			h.Logger.Error("Error while stoppping server: ", err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponse("Server stop successfully"))
	}
	// No Server is running on given port.
	resp := fmt.Sprintf("No Mock server is running on port %d with protocol %v", reqBody.Port, reqBody.Protocol)
	h.Logger.Debug(resp)
	return c.Status(fiber.StatusOK).JSON(response.BuildSuccessResponse(resp))
}
