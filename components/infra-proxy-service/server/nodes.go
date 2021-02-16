package server

import (
	"context"
	"encoding/json"

	chef "github.com/go-chef/chef"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/chef/automate/api/interservice/infra_proxy/request"
	"github.com/chef/automate/api/interservice/infra_proxy/response"
	"github.com/chef/automate/components/infra-proxy-service/validation"
)

// GetNodes fetches the nodes from chef infra server
func (s *Server) GetNodes(ctx context.Context, req *request.Nodes) (*response.Nodes, error) {
	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"name":             []string{"name"},
		"fqdn":             []string{"fqdn"},
		"ipaddress":        []string{"ipaddress"},
		"platform":         []string{"platform"},
		"chef_environment": []string{"chef_environment"},
		"policy_group":     []string{"policy_group"},
		"chef_guid":        []string{"chef_guid"},
		"uptime":           []string{"uptime"},
		"ohai_time":        []string{"ohai_time"},
	}

	res, err := c.SearchObjectsWithDefaults("node", req.SearchQuery, params)
	if err != nil {
		return nil, err
	}

	return &response.Nodes{
		Nodes: fromSearchAPIToNodes(res),
		Page:  res.Page,
		Total: res.Total,
	}, nil
}

// GetNode fetches the node from chef infra server
func (s *Server) GetNode(ctx context.Context, req *request.Node) (*response.Node, error) {
	err := validation.New(validation.Options{
		Target:          "node",
		Request:         *req,
		RequiredDefault: true,
	}).Validate()

	if err != nil {
		return nil, err
	}

	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Nodes.Get(req.Name)
	if err != nil {
		return nil, err
	}

	defaultAttributes, err := json.Marshal(res.DefaultAttributes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	automaticAttributes, err := json.Marshal(res.AutomaticAttributes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	normalAttributes, err := json.Marshal(res.NormalAttributes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	overrideAttributes, err := json.Marshal(res.OverrideAttributes)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// TODO: Chef Infra Server node detail API is not returning the node ID https://docs.chef.io/api_chef_server/#get-38
	return &response.Node{
		NodeId:              res.Name,
		Name:                res.Name,
		Environment:         res.Environment,
		PolicyName:          res.PolicyName,
		PolicyGroup:         res.PolicyGroup,
		RunList:             res.RunList,
		Tags:                SafeSliceFromMap(res.NormalAttributes, "tags"),
		DefaultAttributes:   string(defaultAttributes),
		AutomaticAttributes: string(automaticAttributes),
		NormalAttributes:    string(normalAttributes),
		OverrideAttributes:  string(overrideAttributes),
	}, nil
}

// DeleteNode deletes the node by name
func (s *Server) DeleteNode(ctx context.Context, req *request.DeleteNode) (*response.DeleteNode, error) {
	err := validation.New(validation.Options{
		Target:  "node",
		Request: *req,
		Rules: validation.Rules{
			"OrgId":    []string{"required"},
			"ServerId": []string{"required"},
			"Name":     []string{"required"},
		},
	}).Validate()

	if err != nil {
		return nil, err
	}

	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}

	err = c.client.Nodes.Delete(req.Name)
	if err != nil {
		return nil, err
	}

	return &response.DeleteNode{
		Name: req.Name,
	}, nil
}

// UpdateNode update the node attributes
func (s *Server) UpdateNode(ctx context.Context, req *request.UpdateNode) (*response.UpdateNode, error) {
	c, err := s.createClient(ctx, req.OrgId, req.ServerId)
	if err != nil {
		return nil, err
	}

	automatic, err := StructToJSON(req.AutomaticAttributes)
	if err != nil {
		return nil, err
	}

	normal, err := StructToJSON(req.NormalAttributes)
	if err != nil {
		return nil, err
	}

	defaults, err := StructToJSON(req.DefaultAttributes)
	if err != nil {
		return nil, err
	}

	override, err := StructToJSON(req.OverrideAttributes)
	if err != nil {
		return nil, err
	}

	_, err = c.client.Nodes.Put(chef.Node{
		Name:                req.Name,
		Environment:         req.Environment,
		RunList:             req.RunList,
		AutomaticAttributes: automatic.(map[string]interface{}),
		NormalAttributes:    normal.(map[string]interface{}),
		DefaultAttributes:   defaults.(map[string]interface{}),
		OverrideAttributes:  override.(map[string]interface{}),
		PolicyName:          req.PolicyName,
		PolicyGroup:         req.PolicyGroup,
	})
	if err != nil {
		return nil, err
	}

	return &response.UpdateNode{
		Name: req.Name,
	}, nil
}

// This return the response node attributes array
// parse by the getting the attributes from partial search query.
func fromSearchAPIToNodes(sr chef.SearchResult) []*response.NodeAttribute {
	results := make([]*response.NodeAttribute, len(sr.Rows))
	index := 0
	for _, element := range sr.Rows {
		m := element.(map[string]interface{})["data"].(map[string]interface{})
		results[index] = &response.NodeAttribute{
			Id:          SafeStringFromMap(m, "chef_guid"),
			Name:        SafeStringFromMap(m, "name"),
			Fqdn:        SafeStringFromMap(m, "fqdn"),
			IpAddress:   SafeStringFromMap(m, "ipaddress"),
			CheckIn:     SafeStringFromMapFloat(m, "ohai_time"),
			Environment: SafeStringFromMap(m, "chef_environment"),
			Platform:    SafeStringFromMap(m, "platform"),
			PolicyGroup: SafeStringFromMap(m, "policy_group"),
			Uptime:      SafeStringFromMap(m, "uptime"),
		}
		index++
	}

	return results
}
