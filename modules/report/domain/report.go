package domain
type ProfitReport struct {
	Dimension    string  `json:"dimension"`
	DimensionKey string  `json:"dimension_key"`
	Revenue      float64 `json:"revenue"`
	Cost         float64 `json:"cost"`
	Profit       float64 `json:"profit"`
	Margin       float64 `json:"margin"`
	OrderCount   int64   `json:"order_count"`
}

type ReportService struct {
	orderRevenue   map[string]float64
	orderCost      map[string]float64
	clientRevenue  map[string]float64
	clientCost     map[string]float64
	routeRevenue   map[string]float64
	routeCost      map[string]float64
}

func NewReportService() *ReportService {
	return &ReportService{
		orderRevenue: make(map[string]float64),
		orderCost:    make(map[string]float64),
		clientRevenue: make(map[string]float64),
		clientCost:   make(map[string]float64),
		routeRevenue: make(map[string]float64),
		routeCost:    make(map[string]float64),
	}
}

func (s *ReportService) RecordOrder(orderNo, clientName, routeName string, revenue, cost float64) {
	s.orderRevenue[orderNo] += revenue; s.orderCost[orderNo] += cost
	s.clientRevenue[clientName] += revenue; s.clientCost[clientName] += cost
	s.routeRevenue[routeName] += revenue; s.routeCost[routeName] += cost
}

func (s *ReportService) OrderProfit() []ProfitReport {
	var r []ProfitReport
	for k,v := range s.orderRevenue { c:=s.orderCost[k]; r=append(r,ProfitReport{Dimension:"order",DimensionKey:k,Revenue:v,Cost:c,Profit:v-c,Margin:margin(v,c),OrderCount:1}) }
	return r
}
func (s *ReportService) ClientProfit() []ProfitReport {
	var r []ProfitReport
	for k,v := range s.clientRevenue { c:=s.clientCost[k]; r=append(r,ProfitReport{Dimension:"client",DimensionKey:k,Revenue:v,Cost:c,Profit:v-c,Margin:margin(v,c)}) }
	return r
}
func (s *ReportService) RouteProfit() []ProfitReport {
	var r []ProfitReport
	for k,v := range s.routeRevenue { c:=s.routeCost[k]; r=append(r,ProfitReport{Dimension:"route",DimensionKey:k,Revenue:v,Cost:c,Profit:v-c,Margin:margin(v,c)}) }
	return r
}
func margin(r,c float64) float64 { if r==0{return 0}; return ((r-c)/r)*100 }
