package v1_test

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/logger"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/server"
	v1 "github.com/chef/automate/components/automate-cli/pkg/verifyserver/server/api/v1"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/services/externalpostgresqlservice"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/utils/fiberutils"
	"github.com/gofiber/fiber"
	"github.com/stretchr/testify/assert"
)

func SetupMockExternalPostgresqlService() externalpostgresqlservice.ISExternalPostgresqlService {
	return &externalpostgresqlservice.MockExternalPostgresqlService{
		GetPgConnectionFunc: func(models.ExternalPgRequest) models.ExternalPgConnectionDetails {
			return models.ExternalPgConnectionDetails{
				Passed:        true,
				Title:         "S3 connection test",
				SuccessMsg:    "Machine is able to connect with S3 using the provided access key and secret key",
				ErrorMsg:      "",
				ResolutionMsg: "",
			}
		},
	}
}

func SetupExternalPostgresqlHandlers(pg externalpostgresqlservice.ISExternalPostgresqlService) *fiber.App {
	log := logger.NewLogger(true)
	fconf := &fiber.Settings{
		ServerHeader: server.SERVICE,
		ErrorHandler: fiberutils.CustomErrorHandler,
	}
	app := fiber.New(fconf)
	handler := v1.NewHandler(log).
	AddExternalPostgresqlService(pg)
	vs := &server.VerifyServer{
		Port:    server.DEFAULT_PORT,
		Log:     log,
		App:     app,
		Handler: handler,
	}
	vs.Setup()
	return vs.App
}

func TestExternalPostgresql(t *testing.T) {
	tests := []struct {
		description  string
		expectedCode int
		expectedBody string
	}{
		{
			description:  "200:success status route",
			expectedCode: 200,

			expectedBody: "{\"passed\": \"true\",\"checks\": [{\"title\": \"Connection successfully tested\",\"passed\": true,\"status\": \"PASS\",\"success_msg\": \"Machine is able to connect with External Managed Postgres\",\"error_msg\": \"\",\"resolution_msg\": \"\",\"debug_msg\": \"\"}]}}",
		},
	}	
	statusEndpoint := "/api/v1/checks/external-postgresql"
	// Setup the app as it is done in the main function
	app := SetupExternalPostgresqlHandlers(SetupMockExternalPostgresqlService())

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			req := httptest.NewRequest("POST", statusEndpoint, nil)
			req.Header.Add("Content-Type", "application/json")
			res, err := app.Test(req, -1)
			assert.NoError(t, err)
			body, err := ioutil.ReadAll(res.Body)
			t.Log(body, "body")
			assert.NoError(t, err, test.description)
			assert.Contains(t, string(body), test.expectedBody)
			assert.Equal(t, res.StatusCode, test.expectedCode)
		})
	}
}