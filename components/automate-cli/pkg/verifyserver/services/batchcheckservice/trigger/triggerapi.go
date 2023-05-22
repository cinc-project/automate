package trigger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chef/automate/lib/logger"
	"github.com/gofiber/fiber"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/constants"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
)

// RunCheck triggers API on all the nodes present in the config like automate,chef server, opensearch, postgres and bastion
func RunCheck(config models.Config, log logger.Logger, port string, path string, depState string) []models.CheckTriggerResponse {
	var result []models.CheckTriggerResponse
	count := config.Hardware.AutomateNodeCount +
		config.Hardware.ChefInfraServerNodeCount +
		config.Hardware.PostgresqlNodeCount +
		config.Hardware.OpenSearchNodeCount

	outputCh := make(chan models.CheckTriggerResponse)

	// added one for bastion node
	if path == constants.SOFTWARE_VERSION_CHECK_API_PATH || path == constants.SYSTEM_RESOURCE_CHECK_API_PATH {
		count = count + 1
		endpoint := prepareEndpoint(path, "127.0.0.1", port, constants.BASTION, depState)
		go triggerCheckAPI(endpoint, "127.0.0.1", constants.BASTION, outputCh, http.MethodGet, nil)
	}

	for _, ip := range config.Hardware.AutomateNodeIps {
		endpoint := prepareEndpoint(path, ip, port, constants.AUTOMATE, depState)
		go triggerCheckAPI(endpoint, ip, constants.AUTOMATE, outputCh, http.MethodGet, nil)
	}
	for _, ip := range config.Hardware.ChefInfraServerNodeIps {
		endpoint := prepareEndpoint(path, ip, port, constants.CHEF_INFRA_SERVER, depState)
		go triggerCheckAPI(endpoint, ip, constants.CHEF_INFRA_SERVER, outputCh, http.MethodGet, nil)
	}
	for _, ip := range config.Hardware.OpenSearchNodeIps {
		endpoint := prepareEndpoint(path, ip, port, constants.OPENSEARCH, depState)
		go triggerCheckAPI(endpoint, ip, constants.OPENSEARCH, outputCh, http.MethodGet, nil)
	}
	for _, ip := range config.Hardware.PostgresqlNodeIps {
		endpoint := prepareEndpoint(path, ip, port, constants.POSTGRESQL, depState)
		go triggerCheckAPI(endpoint, ip, constants.POSTGRESQL, outputCh, http.MethodGet, nil)
	}

	for i := 0; i < count; i++ {
		select {
		case res := <-outputCh:
			result = append(result, res)
		}
	}

	return result
}

// RunCheckOnSpecifiedNodes triggers the API on gives node ips only, requires for various API's like S3/Minio backup config
func RunCheckOnSpecifiedNode(nodeIps []string, log logger.Logger, port string, path string, nodeType string, method string, reqBody interface{}) []models.CheckTriggerResponse {
	log.Debugf("Triggering the api call for specified nodes only")
	outputCh := make(chan models.CheckTriggerResponse)
	for _, ip := range nodeIps {
		log.Debugf("Triggering api %s for the node %s", path, ip)
		endpoint := prepareEndpoint(path, ip, port, nodeType, "")
		go triggerCheckAPI(endpoint, ip, nodeType, outputCh, method, reqBody)
	}

	return getResultFromOutputChan(len(nodeIps), outputCh)
}

// RunCheckWithEndPointSpecified triggers the API on given endpoint with node ip, node type
func RunCheckWithEndPointSpecified(endPoint string, log logger.Logger, reqList []models.NodeIpRequest, method string) []models.CheckTriggerResponse {
	log.Debugf("Triggering the api call for specified nodes only")

	outputCh := make(chan models.CheckTriggerResponse)

	for _, req := range reqList {
		log.Debugf("Triggering api on enpoint %s for the node %s", endPoint, req.NodeIP)
		go triggerCheckAPI(endPoint, req.NodeIP, req.NodeType, outputCh, method, req.Request)

	}

	return getResultFromOutputChan(len(reqList), outputCh)
}

func getResultFromOutputChan(reqList int, outputCh chan models.CheckTriggerResponse) []models.CheckTriggerResponse {
	var result []models.CheckTriggerResponse

	for i := 0; i < reqList; i++ {
		select {
		case res := <-outputCh:
			result = append(result, res)
		}
	}

	return result
}

func prepareEndpoint(path, ip, port, nodeType, depState string) string {
	endPoint := ""
	if path == constants.SOFTWARE_VERSION_CHECK_API_PATH {
		endPoint = fmt.Sprintf("http://%s:%s%s?node_type=%s", ip, port, path, nodeType)

	} else if path == constants.SYSTEM_RESOURCE_CHECK_API_PATH {
		endPoint = fmt.Sprintf("http://%s:%s%s?node_type=%s&deployment_state=%s", ip, port, path, nodeType, depState)

	} else {
		endPoint = fmt.Sprintf("http://%s:%s%s", ip, port, path)
	}

	return endPoint
}

func triggerCheckAPI(endPoint, host, nodeType string, output chan<- models.CheckTriggerResponse, method string, reqBody interface{}) {
	var ctr models.CheckTriggerResponse

	reader, err := interfaceToIOReader(reqBody)
	if err != nil {
		output <- models.CheckTriggerResponse{
			Error: &fiber.Error{
				Code:    http.StatusBadRequest,
				Message: fmt.Sprintf("error while reading the request body: %s", err.Error()),
			},
			Host:     host,
			NodeType: nodeType,
		}
		return
	}

	req, err := http.NewRequest(method, endPoint, reader)
	if err != nil {
		output <- models.CheckTriggerResponse{
			Host: host,
			Error: &fiber.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("error while creating the request:%s", err.Error()),
			},
			NodeType: nodeType,
		}
		return
	}

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		output <- models.CheckTriggerResponse{
			Host: host,
			Error: &fiber.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("error while connecting to the endpoint:%s", err.Error()),
			},
			NodeType: nodeType,
		}
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		output <- models.CheckTriggerResponse{
			Host: host,
			Error: &fiber.Error{
				Code:    resp.StatusCode,
				Message: "error while connecting to the endpoint, received invalid status code",
			},
			NodeType: nodeType,
		}
		return
	}

	if err := json.NewDecoder(resp.Body).Decode(&ctr); err != nil {
		output <- models.CheckTriggerResponse{
			Host: host,
			Error: &fiber.Error{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("error while parsing the response data:%s", err.Error()),
			},
			NodeType: nodeType,
		}
		return
	}

	ctr.Host = host
	ctr.NodeType = nodeType
	output <- ctr
}

func interfaceToIOReader(body interface{}) (io.Reader, error) {
	var reader io.Reader
	if body != nil {
		bx, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reader = bytes.NewBuffer(bx)

	}
	return reader, nil
}
