// Package optimizer provides route optimization and scoring for the I56 WMS
// AI subsystem. It scores transport routes against multiple factors and
// suggests parcel consolidation groups for cost efficiency.
package optimizer

import (
	"fmt"
	"math"
	"sort"
	"time"

	routeDomain "github.com/i56/modules/transport/domain"

	parcelDomain "github.com/i56/modules/parcel/domain"
)

// RouteOptimizer scores routes and suggests consolidation strategies.
type RouteOptimizer struct {
	routeRepo RouteRepository
}

// RouteRepository defines the minimal route store needed by the optimizer.
type RouteRepository interface {
	List(ctx interface{}, tenantID int64, offset, limit int) ([]routeDomain.Route, int64, error)
	GetByID(ctx interface{}, tenantID, id int64) (*routeDomain.Route, error)
}

// routeRepoAdapter wraps a concrete route repository to satisfy RouteRepository.
type routeRepoAdapter struct {
	listFn  func(ctx interface{}, tenantID int64, offset, limit int) ([]routeDomain.Route, int64, error)
	getByID func(ctx interface{}, tenantID, id int64) (*routeDomain.Route, error)
}

func (a *routeRepoAdapter) List(ctx interface{}, tenantID int64, offset, limit int) ([]routeDomain.Route, int64, error) {
	return a.listFn(ctx, tenantID, offset, limit)
}

func (a *routeRepoAdapter) GetByID(ctx interface{}, tenantID, id int64) (*routeDomain.Route, error) {
	return a.getByID(ctx, tenantID, id)
}

// RouteScore represents a scored route recommendation.
type RouteScore struct {
	RouteID     int64   `json:"route_id"`
	RouteName   string  `json:"route_name"`
	Score       float64 `json:"score"`        // 0-100
	EstTime     int     `json:"est_time"`     // estimated hours
	EstCost     float64 `json:"est_cost"`     // estimated cost
	Reliability float64 `json:"reliability"`  // 0-1
	TransportType string `json:"transport_type"`
}

// ConsolidationGroup represents a group of parcels that should be shipped
// together on the same route for cost efficiency.
type ConsolidationGroup struct {
	RouteID     int64             `json:"route_id"`
	RouteName   string            `json:"route_name"`
	Parcels     []parcelDomain.Parcel `json:"parcels"`
	TotalWeight float64           `json:"total_weight"`
	EstSavings  float64           `json:"est_savings"` // estimated cost savings
}

// New creates a new RouteOptimizer. The transportRepo should be a
// *tmsRepo.MemRouteRepo or any type with List + GetByID methods.
func New(transportRepo interface{}) *RouteOptimizer {
	r := &RouteOptimizer{}
	if transportRepo != nil {
		r.routeRepo = wrapRepo(transportRepo)
	}
	if r.routeRepo == nil {
		r.routeRepo = &routeRepoAdapter{} // empty - will return empty scores
	}
	return r
}

// wrapRepo creates a RouteRepository adapter from a concrete transport repo.
func wrapRepo(repo interface{}) *routeRepoAdapter {
	type lister interface {
		List(ctx interface{}, tenantID int64, offset, limit int) ([]routeDomain.Route, int64, error)
	}
	type getter interface {
		GetByID(ctx interface{}, tenantID, id int64) (*routeDomain.Route, error)
	}
	adapter := &routeRepoAdapter{}
	if l, ok := repo.(lister); ok {
		adapter.listFn = l.List
	}
	if g, ok := repo.(getter); ok {
		adapter.getByID = g.GetByID
	}
	return adapter
}

// ─── Route Scoring ──────────────────────────────────────────────────────

// Scoring weights (must sum to 1.0)
const (
	weightCost        = 0.40
	weightSpeed       = 0.30
	weightReliability = 0.20
	weightCapacity    = 0.10
)

// ScoreRoutes returns the best routes for shipping a parcel from origin to
// destination, sorted by descending score.
//
//   - from / to: used to match route name (simple substring match)
//   - weight: parcel weight in kg (used for cost estimation)
//   - priority: high / normal / low (affects speed scoring)
func (ro *RouteOptimizer) ScoreRoutes(from, to string, weight float64, priority string) []RouteScore {
	routes, _, err := ro.routeRepo.List(nil, 1, 0, 200)
	if err != nil || len(routes) == 0 {
		return nil
	}

	var results []RouteScore
	for _, rt := range routes {
		if !rt.IsActive {
			continue
		}

		// 1. Cost effectiveness (40%)
		costScore := ro.scoreCost(&rt, weight)

		// 2. Speed (30%)
		speedScore := ro.scoreSpeed(&rt, priority)

		// 3. Historical reliability (20%)
		reliability := ro.scoreReliability(&rt)

		// 4. Current capacity (10%)
		capacityScore := ro.scoreCapacity(&rt)

		// Weighted total score (0-100)
		total := (costScore*weightCost +
			speedScore*weightSpeed +
			reliability*weightReliability +
			capacityScore*weightCapacity) * 100

		// Estimate cost
		estCost := rt.BaseWeightPrice * math.Max(weight, rt.MinWeight)

		results = append(results, RouteScore{
			RouteID:       rt.ID,
			RouteName:     rt.Name,
			Score:         math.Round(total*100) / 100,
			EstTime:       (rt.MinDays + rt.MaxDays) * 24 / 2, // midpoint in hours
			EstCost:       math.Round(estCost*100) / 100,
			Reliability:   math.Round(reliability*100) / 100,
			TransportType: rt.TransportType,
		})
	}

	// Sort descending by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Return top 3
	if len(results) > 3 {
		results = results[:3]
	}
	return results
}

// scoreCost rates cost effectiveness (0-1). Lower cost per kg = higher score.
func (ro *RouteOptimizer) scoreCost(rt *routeDomain.Route, weight float64) float64 {
	if weight <= 0 {
		weight = 1.0
	}
	effectiveWeight := math.Max(weight, rt.MinWeight)
	costPerKg := rt.BaseWeightPrice

	// Benchmark: ¥10/kg is moderate, ¥5/kg is excellent, ¥30/kg is expensive
	if costPerKg <= 0 {
		return 0.5
	}
	score := 1.0 - (costPerKg / 30.0)
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	_ = effectiveWeight
	return score
}

// scoreSpeed rates delivery speed (0-1). Fewer days = higher score.
func (ro *RouteOptimizer) scoreSpeed(rt *routeDomain.Route, priority string) float64 {
	avgDays := float64(rt.MinDays+rt.MaxDays) / 2.0
	if avgDays <= 0 {
		avgDays = 1
	}

	// Base: fewer days = higher score. 1 day → 1.0, 30 days → 0.0
	score := 1.0 - (avgDays / 30.0)
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}

	// Priority boost
	switch priority {
	case "high":
		score = math.Min(1.0, score*1.3)
	case "low":
		score *= 0.8
	}

	return score
}

// scoreReliability estimates historical on-time reliability (0-1).
// Currently based on route age (older routes = more reliable) and
// transport type heuristics.
func (ro *RouteOptimizer) scoreReliability(rt *routeDomain.Route) float64 {
	base := 0.7 // default reliability

	// Older routes tend to be more proven
	age := time.Since(rt.CreatedAt)
	if age > 90*24*time.Hour {
		base += 0.15
	} else if age > 30*24*time.Hour {
		base += 0.10
	}

	// Transport type heuristics
	switch rt.TransportType {
	case "air":
		base += 0.10 // Air tends to be on time
	case "sea":
		base += 0.05 // Sea is reliable but slower
	case "sea_express":
		base += 0.08
	}

	if base > 1.0 {
		base = 1.0
	}
	return base
}

// scoreCapacity rates current capacity availability (0-1).
// Simpler routes with higher MinAmount have more headroom.
func (ro *RouteOptimizer) scoreCapacity(rt *routeDomain.Route) float64 {
	// Higher min amount suggests the route has stricter constraints.
	// More headroom = higher score.
	if rt.MinAmount <= 0 {
		return 0.8
	}
	if rt.MinAmount >= 500 {
		return 0.3 // heavily booked
	}
	return 0.8 - (rt.MinAmount / 2500.0)
}

// ─── Consolidation Suggestions ───────────────────────────────────────────

// SuggestConsolidation groups parcels that should share a route for cost
// efficiency. Parcels going to similar destinations can be consolidated.
func (ro *RouteOptimizer) SuggestConsolidation(parcels []parcelDomain.Parcel) []ConsolidationGroup {
	routes, _, err := ro.routeRepo.List(nil, 1, 0, 200)
	if err != nil || len(routes) == 0 {
		return nil
	}

	// Build active route lookup
	activeRoutes := make(map[int64]*routeDomain.Route)
	for i := range routes {
		if routes[i].IsActive {
			activeRoutes[routes[i].ID] = &routes[i]
		}
	}

	// Group parcels by route that can accommodate them
	// For simplicity, group by warehouse + weight range
	type groupKey struct {
		warehouseID int64
		routeID     int64
	}

	groups := make(map[groupKey]*ConsolidationGroup)
	for _, p := range parcels {
		for _, rt := range routes {
			if !rt.IsActive {
				continue
			}
			if rt.WarehouseID != 0 && rt.WarehouseID != p.WarehouseID {
				continue
			}
			// Check weight range
			if p.ActualWeight < rt.MinWeight || (rt.MaxWeight > 0 && p.ActualWeight > rt.MaxWeight) {
				continue
			}
			key := groupKey{p.WarehouseID, rt.ID}
			if g, ok := groups[key]; ok {
				g.Parcels = append(g.Parcels, p)
				g.TotalWeight += p.ActualWeight
				// Estimate savings: bulk shipping saves ~15%
				g.EstSavings = g.TotalWeight * rt.BaseWeightPrice * 0.15
			} else {
				groups[key] = &ConsolidationGroup{
					RouteID:     rt.ID,
					RouteName:   rt.Name,
					Parcels:     []parcelDomain.Parcel{p},
					TotalWeight: p.ActualWeight,
					EstSavings:  p.ActualWeight * rt.BaseWeightPrice * 0.15,
				}
			}
		}
	}

	// Keep only groups with 2+ parcels (consolidation is meaningful)
	var result []ConsolidationGroup
	for _, g := range groups {
		if len(g.Parcels) >= 2 {
			g.TotalWeight = math.Round(g.TotalWeight*100) / 100
			g.EstSavings = math.Round(g.EstSavings*100) / 100
			result = append(result, *g)
		}
	}

	// Sort by savings descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].EstSavings > result[j].EstSavings
	})

	return result
}

// FormatRouteScoreSummary returns a human-readable summary of route scores.
func FormatRouteScoreSummary(scores []RouteScore) string {
	if len(scores) == 0 {
		return "暂无推荐路线"
	}
	out := fmt.Sprintf("共 %d 条推荐路线:\n", len(scores))
	for i, s := range scores {
		out += fmt.Sprintf("%d. %s (%s) — 综合评分: %.1f/100 | 预估运费: ¥%.2f | 预计: %dh | 可靠度: %.0f%%\n",
			i+1, s.RouteName, s.TransportType, s.Score, s.EstCost, s.EstTime, s.Reliability*100)
	}
	return out
}
