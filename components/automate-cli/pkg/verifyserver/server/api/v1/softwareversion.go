package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber"
)

func (h *Handler) GetSoftwareVersion(c *fiber.Ctx) {
	services := h.StatusService.GetSoftwareVersion()
	c.JSON(response.BuildSuccessResponse(&models.ServiceDetails{
		Status:  "ok",
		Version: services,
	}))
}
