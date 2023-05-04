package externalpostgresqlservice_test

import (
	"testing"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/externalpostgresqlservice"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/utils/db"
	"github.com/stretchr/testify/assert"
)

func TestExternalPostgresqlService(t *testing.T) {
	cs := externalpostgresqlservice.NewExternalPostgresqlService(db.NewDBImpl())
	services := cs.GetPgConnection(models.ExternalPgRequest{})
	assert.Equal(t, models.ExternalPgConnectionDetails{}, services)
}