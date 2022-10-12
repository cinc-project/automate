package pgdb

import "time"

type Storage interface {
	GetUpgradeFlags() (map[string]bool, error)
	GetUpgradeFlagsTimestamp() (map[string]Flag, error)
	UpdateControlFlagToFalse() error
	UpdateControlFlagTimeStamp() error
}

const DayLatestFlag = "day_latest"

const ControlIndexFlag = "control_index"

const CompRunInfoFlag = "comp_run_info"

type Flag struct {
	Flag             string
	Status           bool
	UpgradeTimestamp time.Time
}
