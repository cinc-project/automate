package v1

import (
	"net/http"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber"
)

// var i int

func (h *Handler) NFSMountLoc(c *fiber.Ctx) {
	// resp := models.NFSMountLocResponse{
	// 	Address:       "10.0.0.11",
	// 	Nfs:           "10.0.0.11:/backup_share",
	// 	MountLocation: "/mnt/automate_backups",
	// }
	// resp2 := models.NFSMountLocResponse{
	// 	Address:       "10.0.0.11",
	// 	Nfs:           "10.0.0.11:/backup_shared",
	// 	MountLocation: "/mnt/automate_backups",
	// }
	// if i%2 == 0 {
	// 	c.JSON(response.BuildSuccessResponse(resp2))
	// } else {
	// 	c.JSON(response.BuildSuccessResponse(resp))
	// }
	// i++

	resp := models.NFSMountLocResponse{
		Address:       "10.0.0.11",
		Nfs:           "10.0.0.11:/backup_share",
		MountLocation: "/mnt/automate_backups",
	}

	c.JSON(response.BuildSuccessResponse(resp))
	// value := make(chan int)
	// resp := math.Inf(1)

	// c.JSON(response.BuildFailedResponse(&fiber.Error{}))
	// c.JSON(response.BuildSuccessResponse(resp))
}

func (h *Handler) NFSMount(c *fiber.Ctx) {
	reqBody := models.NFSMountRequest{}
	if err := c.BodyParser(&reqBody); err != nil {
		h.Logger.Error(err.Error())
		c.JSON(response.BuildFailedResponse(&fiber.Error{Code: http.StatusBadRequest, Message: "Invalid Body Request"}))
		return
	}
	if len(reqBody.AutomateNodeIPs) == 0 || len(reqBody.ChefInfraServerNodeIPs) == 0 || len(reqBody.PostgresqlNodeIPs) == 0 || len(reqBody.OpensearchNodeIPs) == 0 || reqBody.MountLocation == "" {
		c.JSON(response.BuildFailedResponse(&fiber.Error{Code: http.StatusBadRequest, Message: "Give all the required body parameters"}))
		return
	}

	NFSMountDetails := h.NFSMountService.GetNFSMountDetails(reqBody, false)
	c.JSON(response.BuildSuccessResponse(NFSMountDetails))
}
