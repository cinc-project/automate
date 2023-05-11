package v1

import (
	"fmt"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber"
)

func (h *Handler)CheckExternalPostgresql(c *fiber.Ctx)  {
	externalPostgresqlRequest := new(models.ExternalPgRequest)
	if err := c.BodyParser(&externalPostgresqlRequest); err != nil {
		errString := fmt.Sprintf("External postgresql config request body parsing failed: %v", err.Error())
		h.Logger.Error(fmt.Errorf(errString))
		c.Status(fiber.StatusBadRequest).JSON(response.BuildFailedResponse(&fiber.Error{Message: err.Error()}))
		return
	}
	pgConnection := h.ExternalPostgresqlService.GetPgConnection(*externalPostgresqlRequest)
	c.JSON(response.BuildSuccessResponse(&models.ExternalPgResponse{
		Passed : true,
		Checks : []models.ExternalPgConnectionDetails{
			pgConnection,
		},
	}))

}
