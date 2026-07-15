// Package service provides business logic services for the WMS backend.
package service

import (
	custRepo "github.com/i56/modules/customer/repository"
	"github.com/i56/i56-apps/i56-wms/internal/domain"
)

// FinanceService provides finance-related business logic using domain stores and repos.
type FinanceService struct {
	ledgerRepo *custRepo.MemLedgerRepo
}

// NewFinanceService creates a new FinanceService.
func NewFinanceService(lr *custRepo.MemLedgerRepo) *FinanceService {
	return &FinanceService{ledgerRepo: lr}
}

// RevenueReport generates a revenue report from monthly statement data.
func (s *FinanceService) RevenueReport() map[string]interface{} {
	statements := domain.MonthlyStatementStore.List()
	var totalRevenue, totalPaid float64
	for _, stmt := range statements {
		totalRevenue += stmt.Total
		totalPaid += stmt.PaidAmount
	}
	return map[string]interface{}{
		"report":       "revenue",
		"total_revenue": totalRevenue,
		"total_paid":   totalPaid,
		"outstanding":  totalRevenue - totalPaid,
		"count":        len(statements),
	}
}

// CostReport generates a cost report from recharge data.
func (s *FinanceService) CostReport() map[string]interface{} {
	recharges := domain.ClientRechargeStore.List()
	var totalCost float64
	for _, r := range recharges {
		totalCost += r.Amount
	}
	return map[string]interface{}{
		"report":     "cost",
		"total_cost":  totalCost,
		"count":       len(recharges),
	}
}
