package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/gofiber/fiber"
)

func (h *Handler)CheckExternalPostgresql(c *fiber.Ctx)  {
	externalPostgresqlRequest := new(models.ExternalPgRequest)
	if err := c.BodyParser(&externalPostgresqlRequest); err != nil {
		c.Status(503).Send(err)
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
