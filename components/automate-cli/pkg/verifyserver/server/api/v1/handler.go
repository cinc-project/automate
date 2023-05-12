package v1

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/opensearchbackupservice"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/statusservice"
)

type Handler struct {
	Logger          logger.ILogger
	StatusService   statusservice.IStatusService
	OSBackupService opensearchbackupservice.IOSS3BackupService
}

func NewHandler(logger logger.ILogger) *Handler {
	return &Handler{Logger: logger}
}

func (h *Handler) AddStatusService(ss statusservice.IStatusService) *Handler {
	h.StatusService = ss
	return h
}

func (h *Handler) AddOSS3BackupService(ss opensearchbackupservice.IOSS3BackupService) *Handler {
	h.OSBackupService = ss
	return h

}
