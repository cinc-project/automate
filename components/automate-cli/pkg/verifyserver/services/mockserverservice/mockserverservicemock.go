package mockserverservice

import "github.com/chef/automate/components/automate-cli/pkg/verifyserver/models"

type MockServersServiceMock struct {
	StartFunc                                      func(cfg *models.StartMockServerRequestBody) error
	StopFunc                                       func(cfg *models.StopMockServerRequestBody) error
	IsMockServerRunningOnGivenPortFunc             func(port int) bool
	IsMockServerRunningOnGivenPortAndProctocolFunc func(port int, protocol string) bool
}

func (mss *MockServersServiceMock) Start(cfg *models.StartMockServerRequestBody) error {
	return mss.StartFunc(cfg)
}

func (mss *MockServersServiceMock) Stop(cfg *models.StopMockServerRequestBody) error {
	return mss.StopFunc(cfg)
}

func (mss *MockServersServiceMock) IsMockServerRunningOnGivenPort(port int) bool {
	return mss.IsMockServerRunningOnGivenPortFunc(port)
}

func (mss *MockServersServiceMock) IsMockServerRunningOnGivenPortAndProctocol(port int, protocol string) bool {
	return mss.IsMockServerRunningOnGivenPortAndProctocolFunc(port, protocol)
}
