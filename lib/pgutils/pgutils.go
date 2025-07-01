package pgutils

import (
	"fmt"
	"regexp"
	"strings"
)

func EscapeLiteralForPG(value string) string {
	value = strings.Replace(value, "'", "''", -1)
	value = strings.Replace(value, `\`, `\\`, -1)
	return value
}

func EscapeLiteralForPGPatternMatch(value string) string {
	value = EscapeLiteralForPG(value)
	value = strings.Replace(value, `_`, `\_`, -1)
	value = strings.Replace(value, `%`, `\%`, -1)
	return value
}

// Ensuring str only has the characters we support at the moment
// TODO: unit tests
func IsSqlSafe(str string) bool {
	r := regexp.MustCompile("^[a-zA-Z0-9\\._-]*$")
	return r.MatchString(str)
}

var blacklistedPatterns =[]string{
	"--",";","/*","*/","'","\"","\\","`",
	" OR ", " AND ", " DROP ", " SELECT ", " INSERT ", " UPDATE ", " DELETE ",
}

func isSuspiciousValue(value string) bool {
	for _, pattern := range blacklistedPatterns {
		if strings.Contains(strings.ToUpper(value), pattern) {
			return true
		}
	}
	return false
}

func ValidateColumnNames(tableName string, columnNames []string, allowed map[string]map[string]bool) error {
	allowedColumns, ok := allowed[tableName]
	if !ok {
		return fmt.Errorf("table %s is not recognized", tableName)
	}
	for _, columnName := range columnNames {
		if !allowedColumns[columnName] {
			return fmt.Errorf("invalid column %s for table %s", columnName, tableName)
		}
		if isSuspiciousValue(columnName) {
			return fmt.Errorf("suspicious column name detected: %s", columnName)
		}
	}
	return nil
}

