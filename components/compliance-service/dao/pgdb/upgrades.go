package pgdb

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strings"
)

type UpgradesDB struct {
	DB *DB
}

func NewDB(db *DB) *UpgradesDB {
	return &UpgradesDB{db}
}

//UpdateControlFlagToFalse updates the control index flags to false
func (u *UpgradesDB) UpdateControlFlagToFalse() error {
	_, err := u.DB.Exec(getUpdateQuery(ControlIndexFlag))
	if err != nil {
		return errors.Wrapf(err, "Unable to Control Index Flag to db")
	}
	return nil
}

//UpdateControlFlagTimestamp updates the control index flags to false
func (u *UpgradesDB) UpdateControlFlagTimeStamp() error {
	_, err := u.DB.Exec(getUpdateTimeStampQuery(ControlIndexFlag))
	if err != nil {
		return errors.Wrapf(err, "Unable to Control Index Flag to db")
	}
	return nil
}

//GetUpgradeFlags Gets the all the upgrade flags and Status from the pg database
func (u *UpgradesDB) GetUpgradeFlags() (map[string]bool, error) {
	flagMap := make(map[string]bool)

	logrus.Info("Inside the comp run info Flag")
	flags := []string{ControlIndexFlag}
	rows, err := u.DB.Query(getQueryForFlag(flags))
	if err != nil {
		return flagMap, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Errorf("failed to close db rows: %s", err.Error())
		}
	}()

	for rows.Next() {
		flag := Flag{}
		if err := rows.Scan(&flag.Flag, &flag.Status); err != nil {
			logrus.Errorf("Unable to get the flags with error %v", err)
			return nil, err
		}
		flagMap[flag.Flag] = flag.Status
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error retrieving result rows")
	}
	return flagMap, err
}

//GetUpgradeFlagsTimestamp Gets the all the upgrade flags and upgrade timestamp from the pg database
func (u *UpgradesDB) GetUpgradeFlagsTimestamp() (map[string]Flag, error) {
	flagMap := make(map[string]Flag)

	logrus.Info("Inside the comp run info Flag")
	flags := []string{ControlIndexFlag}
	rows, err := u.DB.Query(getQueryForFlagTimestamp(flags))
	if err != nil {
		return flagMap, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logrus.Errorf("failed to close db rows: %s", err.Error())
		}
	}()

	for rows.Next() {
		flag := Flag{}
		if err := rows.Scan(&flag.Flag, &flag.Status, &flag.UpgradeTimestamp); err != nil {
			logrus.Errorf("Unable to get the flags with error %v", err)
			return nil, err
		}
		flagMap[flag.Flag] = flag
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error retrieving result rows")
	}
	return flagMap, err
}

//getQueryForFlag gets the query for Flag
func getQueryForFlag(flag []string) string {
	flags := `'` + strings.Join(flag, `','`) + `'`
	return fmt.Sprintf("Select upgrade_flag,upgrade_value from upgrade_flags where upgrade_flag in (%s)", flags)
}

//getQueryForFlagTimestamp gets the query for Flag and updatetimestamp
func getQueryForFlagTimestamp(flag []string) string {
	flags := `'` + strings.Join(flag, `','`) + `'`
	return fmt.Sprintf("Select upgrade_flag, upgrade_value, upgrade_timestamp from upgrade_flags where upgrade_flag in (%s)", flags)
}

//getUpdateQuery gets the update query for Flag
func getUpdateQuery(flag string) string {
	return fmt.Sprintf("Update upgrade_flags set upgrade_value=false where upgrade_flag='%s'", flag)
}

//getUpdateTimeStampQuery gets the update query for Flag
func getUpdateTimeStampQuery(flag string) string {
	return fmt.Sprintf("Update upgrade_flags set upgrade_timestamp=current_timestamp where upgrade_flag='%s'", flag)
}
