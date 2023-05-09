package models

type StartMockServerRequestBody struct {
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
	Cert     string `json:"cert"`
	Key      string `json:"key"`
}
