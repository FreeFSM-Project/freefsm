package services

import (
	"context"
	"time"

	"github.com/MartialM1nd/freefsm/internal/ent"
	"github.com/MartialM1nd/freefsm/internal/ent/invoice"
)

type DashboardStats struct {
	TotalCustomers int
	TotalJobs      int
	TotalEstimates int
	TotalInvoices  int
	RevenueMonth   float64
}

type DashboardService struct {
	client *ent.Client
}

func NewDashboardService(client *ent.Client) *DashboardService {
	return &DashboardService{client: client}
}

func (s *DashboardService) Stats(ctx context.Context) (DashboardStats, error) {
	customers, _ := s.client.Customer.Query().Count(ctx)
	jobs, _ := s.client.Job.Query().Count(ctx)
	estimates, _ := s.client.Estimate.Query().Count(ctx)
	invoices, _ := s.client.Invoice.Query().Count(ctx)

	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	monthInvoices, _ := s.client.Invoice.Query().
		Where(invoice.InvoiceDateGTE(startOfMonth), invoice.InvoiceDateLTE(endOfMonth)).
		All(ctx)

	var revenue float64
	for _, i := range monthInvoices {
		payments, _ := ParsePayments(i.Payments)
		for _, p := range payments {
			revenue += p.Amount
		}
	}

	return DashboardStats{
		TotalCustomers: customers,
		TotalJobs:      jobs,
		TotalEstimates: estimates,
		TotalInvoices:  invoices,
		RevenueMonth:   revenue,
	}, nil
}
