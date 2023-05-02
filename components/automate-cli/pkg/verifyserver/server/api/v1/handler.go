package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/s3configservice"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/statusservice"
)

type Handler struct {
	Logger          logger.ILogger
	StatusService   statusservice.IStatusService
	S3ConfigService s3configservice.S3Config
}

func NewHandler(logger logger.ILogger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) AddStatusService(ss statusservice.IStatusService) *Handler {
	h.StatusService = ss
	return h
}
func (h *Handler) AddS3ConfigService(ss s3configservice.S3Config) *Handler {
	h.S3ConfigService = ss
	return h
}
