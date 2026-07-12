package core

import (
	"context"
	"fmt"
	"strings"

	orderRepo "github.com/i56/modules/order/repository"
	orderSvc "github.com/i56/modules/order/service"
	parcelRepo "github.com/i56/modules/parcel/repository"
	parcelSvc "github.com/i56/modules/parcel/service"
	whRepo "github.com/i56/modules/warehouse/repository"
	whSvc "github.com/i56/modules/warehouse/service"
	custRepo "github.com/i56/modules/customer/repository"
)

// BusinessContext provides domain-aware context injection for AI queries.
type BusinessContext struct {
	orderRepo   *orderRepo.MemOrderRepo
	parcelRepo  *parcelRepo.MemParcelRepo
	whRepo      *whRepo.MemWarehouseRepo
	clientRepo  *custRepo.MemClientRepo
	orderSvc    *orderSvc.OrderService
	parcelSvc   *parcelSvc.ParcelService
	whSvc       *whSvc.WarehouseService
}

// NewBusinessContext creates a new BusinessContext with the given repositories and services.
func NewBusinessContext(
	or *orderRepo.MemOrderRepo,
	pr *parcelRepo.MemParcelRepo,
	wr *whRepo.MemWarehouseRepo,
	cr *custRepo.MemClientRepo,
	osvc *orderSvc.OrderService,
	ps *parcelSvc.ParcelService,
	ws *whSvc.WarehouseService,
) *BusinessContext {
	return &BusinessContext{
		orderRepo:  or,
		parcelRepo: pr,
		whRepo:     wr,
		clientRepo: cr,
		orderSvc:   osvc,
		parcelSvc:  ps,
		whSvc:      ws,
	}
}

// GetBusinessContext analyzes the user's query and injects relevant business data
// as context for the AI model.
func (bc *BusinessContext) GetBusinessContext(tenantID int64, query string) string {
	ctx := context.Background()
	var parts []string

	// Check for order-related queries
	if containsAny(query, []string{"订单", "order"}) {
		orders, _, err := bc.orderSvc.List(ctx, tenantID, 0, 20)
		if err == nil && len(orders) > 0 {
			var sb strings.Builder
			sb.WriteString("最近订单:\n")
			for i, o := range orders {
				if i >= 10 {
					break
				}
				sb.WriteString(fmt.Sprintf("- %s: 收件人:%s 状态:%s 金额:¥%.2f 件数:%d\n",
					o.OrderNo, o.RecipientName, string(o.Status), o.TotalPrice, o.ParcelCount))
			}
			parts = append(parts, sb.String())
		}
	}

	// Check for parcel-related queries
	if containsAny(query, []string{"包裹", "parcel", "快递", "tracking"}) {
		parcels, _, err := bc.parcelSvc.List(ctx, tenantID, 0, 30)
		if err == nil && len(parcels) > 0 {
			var sb strings.Builder
			sb.WriteString("包裹列表:\n")
			for i, p := range parcels {
				if i >= 10 {
					break
				}
				sb.WriteString(fmt.Sprintf("- %s: %s 状态:%s 实重:%.2fkg\n",
					p.TrackingNumber, p.ProductName, string(p.Status), p.ActualWeight))
			}
			parts = append(parts, sb.String())
		}
	}

	// Check for warehouse-related queries
	if containsAny(query, []string{"仓库", "warehouse", "库存"}) {
		whs, _, err := bc.whSvc.List(ctx, tenantID, 0, 50)
		if err == nil && len(whs) > 0 {
			var sb strings.Builder
			sb.WriteString("仓库列表:\n")
			for _, w := range whs {
				status := "运营中"
				if !w.IsActive {
					status = "停用"
				}
				sb.WriteString(fmt.Sprintf("- %s (%s): %s 联系人:%s 电话:%s %s\n",
					w.Name, w.Code, w.Address, w.Contact, w.Phone, status))
			}
			parts = append(parts, sb.String())
		}
	}

	// Check for client-related queries
	if containsAny(query, []string{"客户", "client", "余额", "balance"}) {
		clients, _, err := bc.clientRepo.List(ctx, tenantID, 0, 20)
		if err == nil && len(clients) > 0 {
			var sb strings.Builder
			sb.WriteString("客户信息:\n")
			for _, c := range clients {
				sb.WriteString(fmt.Sprintf("- %s (%s): 余额:¥%.2f 邮箱:%s\n",
					c.Name, c.Code, c.Balance, c.ContactEmail))
			}
			parts = append(parts, sb.String())
		}
	}

	if len(parts) == 0 {
		return ""
	}

	return "=== WMS 业务上下文 ===\n" + strings.Join(parts, "\n")
}

// containsAny checks if the query contains any of the given keywords.
func containsAny(query string, keywords []string) bool {
	lower := strings.ToLower(query)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
