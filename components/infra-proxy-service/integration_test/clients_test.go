package integration_test

import (
	"context"
	"testing"

	"github.com/chef/automate/api/interservice/infra_proxy/request"
	"github.com/chef/automate/api/interservice/infra_proxy/response"
	"github.com/chef/automate/components/automate-gateway/handler/infra_proxy"

	"github.com/stretchr/testify/assert"
)


func TestNodesReturnsEmptyList(t *testing.T) {
	// rpc GetNodes (request.Nodes) returns (google.protobuf.ListValue)
	ctx := context.Background()
	req := request.Clients{}

	expected := new(response.Clients)
	res, err := infra_proxy.

	assert.Nil(t, err)
	assert.Equal(t, expected, res)
}