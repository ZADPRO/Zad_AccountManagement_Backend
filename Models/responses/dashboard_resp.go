package responses

import "invoice-backend/Models/dto"

type DashboardSummary struct {
	TotalClients  int     `json:"totalClients"`
	ActiveUsers   int     `json:"activeUsers"`
	TotalRevenue  float64 `json:"totalRevenue"`
	PendingAmount float64 `json:"pendingAmount"`
	OverdueCount  int     `json:"overdueCount"`
}

type DashboardResponse struct {
	dto.BaseResponse
	Data DashboardSummary `json:"data"`
}