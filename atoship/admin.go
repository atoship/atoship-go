package atoship

import "context"

// AdminService handles admin-related operations
type AdminService struct {
	client *Client
}

// Stats represents system statistics
type Stats struct {
	TotalOrders    int     `json:"totalOrders"`
	TotalShipments int     `json:"totalShipments"`
	TotalRevenue   float64 `json:"totalRevenue"`
	ActiveUsers    int     `json:"activeUsers"`
}

// GetStats gets system statistics
func (s *AdminService) GetStats(ctx context.Context) (*Stats, error) {
	var stats Stats
	err := s.client.get(ctx, "/api/admin/stats", &stats)
	return &stats, err
}