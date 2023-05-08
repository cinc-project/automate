package nfsmountservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"
	"github.com/chef/automate/components/automate-cli/pkg/verifyserver/response"
	"github.com/gofiber/fiber"
)

var compareWith models.NFSMountLocResponse

const (
	MOUNT_SUCCESS_MSG    = "NFS mount location found"
	MOUNT_ERROR_MSG      = "NFS mount location not found"
	MOUNT_RESOLUTION_MSG = "NFS volume should be mounted on %s"

	SHARE_SUCCESS_MSG    = "NFS mount location is shared across given nodes"
	SHARE_ERROR_MSG      = "NFS mount location %s is not shared across all given nodes"
	SHARE_RESOLUTION_MSG = "NFS volume should be common across all given nodes at mount location: %s"
)

type INFSService interface {
	GetNFSMountDetails(models.NFSMountRequest, bool) *[]*models.NFSMountResponse
}

type NFSMountService struct {
}

func NewNFSMountService() INFSService {
	return &NFSMountService{}
}

func (nm *NFSMountService) GetNFSMountDetails(reqBody models.NFSMountRequest, test bool) *[]*models.NFSMountResponse {
	respBody := new([]*models.NFSMountResponse)

	for _, ip := range reqBody.AutomateNodeIPs {
		res := doAPICall(ip, test, "automate", reqBody.MountLocation)
		*respBody = append(*respBody, res)
	}

	for _, ip := range reqBody.ChefInfraServerNodeIPs {
		res := doAPICall(ip, test, "chef-infra-server", reqBody.MountLocation)
		*respBody = append(*respBody, res)
	}

	for _, ip := range reqBody.PostgresqlNodeIPs {
		res := doAPICall(ip, test, "postgresql", reqBody.MountLocation)
		*respBody = append(*respBody, res)
	}

	for _, ip := range reqBody.OpensearchNodeIPs {
		res := doAPICall(ip, test, "opensearch", reqBody.MountLocation)
		*respBody = append(*respBody, res)
	}
	return respBody
}

func doAPICall(ip string, test bool, node_type string, mountLocation string) *models.NFSMountResponse {
	node := new(models.NFSMountResponse)
	node.IP = ip
	node.NodeType = node_type
	var url string
	if test {
		url = ip
	} else {
		url = fmt.Sprintf("http://%s:7799", ip)
	}

	// It will trigger /nfs-mount-loc API
	resp, err := triggerAPI(url, mountLocation)
	if err != nil {
		node.Error = fiber.NewError(http.StatusBadRequest, err.Error())
		return node
	}

	result, err := getResultStructFromRespBody(resp.Body)
	if err != nil {
		node.Error = fiber.NewError(http.StatusBadRequest, err.Error())
		return node
	}
	// fmt.Println(result)

	// Test1 - Check for NFS Volume is mounted at correct Mount Location or not
	nfsMounted := checkMount(mountLocation, node, result)

	// Test2 - Check for NFS Volume is Shared among all nodes or not
	checkShare(result, node, nfsMounted)

	return node
}

func triggerAPI(url, mountLocation string) (*http.Response, error) {
	// Request Body for /nfs-mount-loc API
	reqBody := models.NFSMountLocRequest{
		MountLocation: mountLocation,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, errors.New("Failed to Marshal: " + err.Error())
	}

	reqURL := url + "/api/v1/fetch/nfs-mount-loc"
	// fmt.Println(reqURL)
	httpReq, err := http.NewRequest(http.MethodPost, reqURL, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return nil, errors.New("Failed to Create HTTP request: " + err.Error())
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, errors.New("Failed to send the HTTP request: " + err.Error())
	}

	return resp, nil
}

func getResultStructFromRespBody(respBody io.Reader) (*models.NFSMountLocResponse, error) {
	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return nil, errors.New("Cannot able to read data from response body: " + err.Error())
	}

	// Converting API Response Body into Generic Response Struct.
	APIRespStruct := response.ResponseBody{}
	err = json.Unmarshal(body, &APIRespStruct)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal: " + err.Error())
	}

	// If API(/nfs-mount-loc) is itself failing.
	if APIRespStruct.Error != nil {
		return nil, APIRespStruct.Error
	}

	// Converting interface into JSON encoding. APIResp.Result is a interface and for accessing the values we are converting that into json.
	resultByte, err := json.Marshal(APIRespStruct.Result)
	if err != nil {
		return nil, errors.New("Failed to Marshal: " + err.Error())
	}

	resultField := new(models.NFSMountLocResponse)
	// converting JSON into struct.
	err = json.Unmarshal(resultByte, &resultField)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal: " + err.Error())
	}
	// fmt.Println(resultField)

	return resultField, nil
}

func checkMount(mountLocation string, node *models.NFSMountResponse, data *models.NFSMountLocResponse) bool {
	if data.Address != "" {
		check := createCheck("NFS Mount", true, MOUNT_SUCCESS_MSG, "", "")
		node.CheckList = append(node.CheckList, check)
		return true
	} else {
		check := createCheck("NFS Mount", false, "", MOUNT_ERROR_MSG, fmt.Sprintf(MOUNT_RESOLUTION_MSG, mountLocation))
		node.CheckList = append(node.CheckList, check)
		return false
	}
}

func checkShare(data *models.NFSMountLocResponse, node *models.NFSMountResponse, nfsMounted bool) {
	// For First IP we need to store the data for comparing
	if compareWith.Address == "" {
		compareWith = *data
	}
	// nfsMounted is holding volume is mounted or not. If volume is not mounted then how we can check it's shareability
	if nfsMounted && data.Address == compareWith.Address && data.Nfs == compareWith.Nfs && data.MountLocation == compareWith.MountLocation {
		check := createCheck("NFS Mount", true, SHARE_SUCCESS_MSG, "", "")
		node.CheckList = append(node.CheckList, check)
	} else {
		check := createCheck("NFS Mount", false, "", fmt.Sprintf(SHARE_ERROR_MSG, compareWith.MountLocation), fmt.Sprintf(SHARE_RESOLUTION_MSG, compareWith.MountLocation))
		node.CheckList = append(node.CheckList, check)
	}
}

func createCheck(title string, passed bool, successMsg, errorMsg, resolutionMsg string) models.Checks {
	return models.Checks{
		Title:         title,
		Passed:        passed,
		SuccessMsg:    successMsg,
		ErrorMsg:      errorMsg,
		ResolutionMsg: resolutionMsg,
	}
}
