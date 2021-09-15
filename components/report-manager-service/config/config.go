package config

// Configuration for the Report Manager Service
type ReportManager struct {
	Service Service
	Log     Log
}

type Service struct {
	Port    int
	Message string
}

type Log struct {
	Level  string
	Format string
}
