package Services

import (
	"database/sql"
	"fmt"
	"invoice-backend/Models/responses"
	"invoice-backend/Query"
)

func GetDashboardSummary(db *sql.DB) (responses.DashboardSummary, error) {
	var summary responses.DashboardSummary

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

