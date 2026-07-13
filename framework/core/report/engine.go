// Package report provides built-in business reports with real data.
package report

import (
	"context"
	"fmt"
	"sync"
)

// Built-in report names
const (
	ReportDailyOrderSummary   = "daily-order-summary"
	ReportMonthlyRevenue      = "monthly-revenue"
	ReportClientBalance       = "client-balance"
	ReportWarehouseThroughput = "warehouse-throughput"
)

// DataProvider abstracts data access for reports.
type DataProvider interface {
	Orders(ctx context.Context) ([]OrderRow, error)
	Clients(ctx context.Context) ([]ClientRow, error)
	Parcels(ctx context.Context) ([]ParcelRow, error)
}

type OrderRow struct {
	OrderNo, Status, ClientName string
	TotalPrice                  float64
	ParcelCount                 int
	CreatedAt                   string
}

type ClientRow struct {
	ID                          int64
	Name, Code                  string
	Balance, CreditLimit        float64
}

type ParcelRow struct {
	TrackingNo, ProductName, Status, WarehouseName string
	ActualWeight                                   float64
}

// ReportParamDef defines a parameter for a built-in report.
type ReportParamDef struct {
	Name, Label, Type string
	Required          bool
}

// ReportDef defines a built-in report.
type ReportDef struct {
	Name, Title, Query, Display string
	Params                      []ReportParamDef
}

// EngineResult contains the result of executing a report.
type EngineResult struct {
	Name, Title, Display string
	Columns              []string
	Rows                 [][]interface{}
	Total                int
}

// BuiltinEngine executes predefined reports using real data.
type BuiltinEngine struct {
	mu       sync.RWMutex
	reports  map[string]*ReportDef
	provider DataProvider
}

// NewBuiltinEngine creates a new report engine with a data provider.
func NewBuiltinEngine(provider DataProvider) *BuiltinEngine {
	e := &BuiltinEngine{
		reports:  make(map[string]*ReportDef),
		provider: provider,
	}
	e.registerBuiltinReports()
	return e
}

func (e *BuiltinEngine) registerBuiltinReports() {
	e.reports[ReportDailyOrderSummary] = &ReportDef{
		Name: ReportDailyOrderSummary, Title: "每日订单汇总",
		Query: "real", Display: "table",
	}
	e.reports[ReportMonthlyRevenue] = &ReportDef{
		Name: ReportMonthlyRevenue, Title: "月度营收报表",
		Query: "real", Display: "chart",
	}
	e.reports[ReportClientBalance] = &ReportDef{
		Name: ReportClientBalance, Title: "客户余额报表",
		Query: "real", Display: "table",
	}
	e.reports[ReportWarehouseThroughput] = &ReportDef{
		Name: ReportWarehouseThroughput, Title: "仓库吞吐量报表",
		Query: "real", Display: "chart",
	}
}

// Execute runs a report with real data from the provider.
func (e *BuiltinEngine) Execute(ctx context.Context, reportName string, params map[string]interface{}) (*EngineResult, error) {
	e.mu.RLock()
	report, ok := e.reports[reportName]
	e.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("report: unknown report %q", reportName)
	}

	if e.provider == nil {
		return e.fallbackData(report), nil
	}

	columns, rows := e.queryRealData(ctx, report)
	return &EngineResult{
		Name: report.Name, Title: report.Title, Columns: columns,
		Rows: rows, Total: len(rows), Display: report.Display,
	}, nil
}

func (e *BuiltinEngine) queryRealData(ctx context.Context, report *ReportDef) ([]string, [][]interface{}) {
	switch report.Name {
	case ReportDailyOrderSummary:
		orders, _ := e.provider.Orders(ctx)
		dayCounts := map[string]int{}
		dayRevenue := map[string]float64{}
		for _, o := range orders {
			day := o.CreatedAt[:10]
			dayCounts[day]++
			dayRevenue[day] += o.TotalPrice
		}
		rows := [][]interface{}{}
		for d, c := range dayCounts {
			rows = append(rows, []interface{}{d, c, dayRevenue[d]})
		}
		return []string{"日期", "订单数", "总营收"}, rows

	case ReportMonthlyRevenue:
		orders, _ := e.provider.Orders(ctx)
		monthCounts := map[string]int{}
		monthRevenue := map[string]float64{}
		for _, o := range orders {
			month := o.CreatedAt[:7]
			monthCounts[month]++
			monthRevenue[month] += o.TotalPrice
		}
		rows := [][]interface{}{}
		for m, c := range monthCounts {
			avg := 0.0
			if c > 0 {
				avg = monthRevenue[m] / float64(c)
			}
			rows = append(rows, []interface{}{m, c, monthRevenue[m], fmt.Sprintf("%.2f", avg)})
		}
		return []string{"月份", "订单数", "总营收", "客单价"}, rows

	case ReportClientBalance:
		clients, _ := e.provider.Clients(ctx)
		rows := [][]interface{}{}
		for _, c := range clients {
			available := c.CreditLimit - c.Balance
			rows = append(rows, []interface{}{c.ID, c.Name, c.Code, c.Balance, c.CreditLimit, available})
		}
		return []string{"ID", "客户名称", "代码", "余额", "信用额度", "可用额度"}, rows

	case ReportWarehouseThroughput:
		parcels, _ := e.provider.Parcels(ctx)
		wh := map[string]struct{count int; weight float64}{}
		for _, p := range parcels {
			w := wh[p.WarehouseName]
			w.count++
			w.weight += p.ActualWeight
			wh[p.WarehouseName] = w
		}
		rows := [][]interface{}{}
		for name, w := range wh {
			rows = append(rows, []interface{}{name, w.count, w.weight})
		}
		return []string{"仓库", "包裹数", "总重量(kg)"}, rows
	}
	return nil, nil
}

func (e *BuiltinEngine) fallbackData(report *ReportDef) *EngineResult {
	return &EngineResult{
		Name: report.Name, Title: report.Title,
		Columns: []string{"提示"},
		Rows:    [][]interface{}{{"暂无数据 — 需要配置数据源"}},
		Total: 1, Display: report.Display,
	}
}

// ListReports returns all registered built-in reports.
func (e *BuiltinEngine) ListReports() []*ReportDef {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*ReportDef, 0, len(e.reports))
	for _, r := range e.reports {
		result = append(result, r)
	}
	return result
}
