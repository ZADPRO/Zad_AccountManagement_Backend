package Service

import (
	"database/sql"
	"fmt"
	"invoice-backend/Model"
	"invoice-backend/Query"
)

func GetDashboardSummary(db *sql.DB) (Model.DashboardSummary, error) {
	var summary Model.DashboardSummary

	// We use QueryRow because we are getting a single row of aggregated totals
	err := db.QueryRow(Query.GetDashboardStatsQuery).Scan(
		&summary.TotalClients,
		&summary.ActiveUsers,
		&summary.TotalRevenue,
		&summary.PendingAmount,
		&summary.OverdueCount,
	)

	if err != nil {
		return summary, fmt.Errorf("failed to fetch dashboard stats: %w", err)
	}

	return summary, nil
} 

