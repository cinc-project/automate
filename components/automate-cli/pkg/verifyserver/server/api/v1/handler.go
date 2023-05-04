package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/externalpostgresqlservice"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/statusservice"
)

type Handler struct {
	Logger        logger.ILogger
	StatusService statusservice.IStatusService
	ExternalPostgresqlService externalpostgresqlservice.ISExternalPostgresqlService
}

func NewHandler(logger logger.ILogger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) AddStatusService(ss statusservice.IStatusService) *Handler {
	h.StatusService = ss
	return h
}

func (h *Handler) AddExternalPostgresqlService(pg externalpostgresqlservice.ISExternalPostgresqlService) *Handler {
	h.ExternalPostgresqlService = pg
	return h
}
