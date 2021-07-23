package server_test

import (
	"github.com/chef/automate/components/infra-proxy-service/server"
	"github.com/go-chef/chef"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPolicyFile(t *testing.T) {
	policyFileResponse := chef.PolicyGetResponse{
		"testPolicy": chef.PolicyRevision{
			"4f2cbe355786b19360128bb32c003545535e2c96": chef.PolicyRevisionDetail{},
		},
	}

	t.Run("transform policy response to api response", func(t *testing.T) {
		formattedData := server.FromAPIIncludedPolicyfileRevisions(policyFileResponse)
		assert.NotNil(t, formattedData)
		assert.Equal(t, 1, len(formattedData))
		assert.Equal(t, formattedData[0].RevisionId, "4f2cbe355786b19360128bb32c003545535e2c96")
	})
}
