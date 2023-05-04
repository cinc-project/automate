package externalpostgresqlservice

import "github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"

type MockExternalPostgresqlService struct {
	GetPgConnectionFunc func(models.ExternalPgRequest) models.ExternalPgConnectionDetails
}

func (mss *MockExternalPostgresqlService) GetPgConnection(req models.ExternalPgRequest)  models.ExternalPgConnectionDetails {
	return mss.GetPgConnectionFunc(req)
}
