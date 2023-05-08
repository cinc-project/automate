package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	nfsmountservice "github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/nfsmount"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/statusservice"
)

type Handler struct {
	Logger          logger.ILogger
	StatusService   statusservice.IStatusService
	NFSMountService nfsmountservice.INFSService
}

func NewHandler(logger logger.ILogger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) AddStatusService(ss statusservice.IStatusService) *Handler {
	h.StatusService = ss
	return h
}

func (h *Handler) AddNFSMountService(nm nfsmountservice.INFSService) *Handler {
	h.NFSMountService = nm
	return h
}
