package pgdb

import (
	"context"
	"time"

	"github.com/chef/automate/components/applications-service/pkg/storage"
)

type Telemetry struct {
	ID                      string    `db:"id" json:"id"`
	LastTelemetryReportedAt time.Time `db:"last_telemetry_reported_at" json:"last_telemetry_reported_at"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
}

const (
	selectServicesTotalCount = `
SELECT count (DISTINCT supervisior_id) from service_full where health_updated_at between now()::date - 15 AND now()::date 
`
)

// Get last services telemetry reported timestamp
func (db *DB) GetTelemetry(ctx context.Context) (Telemetry, error) {
	var t Telemetry
	rows, err := db.Query(`SELECT id,last_telemetry_reported_at, created_at from telemetry`)
	if err != nil {
		return Telemetry{}, err
	}
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.LastTelemetryReportedAt, &t.CreatedAt)
		if err != nil {
			return Telemetry{}, err
		}
	}
	return t, nil
}

// Get last 15 days services telemetry reported timestamp
func (db *DB) GetUniqueServicesFromPostgres(ctx context.Context) (int64, error) {
	rows, err := db.Query(`SELECT count (DISTINCT supervisior_id) from service_full where health_updated_at between now()::date - 15 AND now()::date`)
	if err != nil {
		return *storage.GetServicesCount{}, err
	}
	return rows.count, nil
}
