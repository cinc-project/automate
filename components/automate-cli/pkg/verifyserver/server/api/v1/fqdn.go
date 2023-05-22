package v1

import (
	"net/http"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber"
)

func (h *Handler) CheckFqdn(c *fiber.Ctx) {
	req := new(models.FqdnRequest)
	if err := c.BodyParser(req); err != nil {
		h.Logger.Error(err.Error())
		c.Next(&fiber.Error{Code: http.StatusBadRequest, Message: "Invalid Body Request"})
		return
	}

	res := h.FqdnService.CheckFqdnReachability(*req)
	c.JSON(response.BuildSuccessResponse(res))
}
