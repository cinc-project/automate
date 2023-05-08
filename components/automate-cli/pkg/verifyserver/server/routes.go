package server

import (
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/utils/fiberutils"
)

func (vs *VerifyServer) SetupRoutes() {
	// Status
	vs.App.Get("/status", vs.Handler.GetStatus)

	apiGroup := vs.App.Group("/api")
	apiV1Group := apiGroup.Group("/v1")
	apiFetchGroup := apiV1Group.Group("/fetch")
	apiFetchGroup.Post("/nfs-mount-loc", vs.Handler.NFSMountLoc)

	apiChecksGroup := apiV1Group.Group("/checks")
	apiChecksGroup.Get("/fqdn", vs.Handler.CheckFqdn)
	apiChecksGroup.Post("/nfs-mount", vs.Handler.NFSMount)

	fiberutils.LogResgisteredRoutes(vs.App.Stack(), vs.Log)
}
