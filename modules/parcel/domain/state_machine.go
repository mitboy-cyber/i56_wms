package domain

// ParcelStateMachine enforces valid status transitions
type ParcelStateMachine struct{}

// ValidTransitions defines allowed status transitions
var ValidParcelTransitions = map[ParcelStatus][]ParcelStatus{
	StatusPreDeclared: {StatusReceived},
	StatusReceived:    {StatusWeighed, StatusStored, StatusAbnormal},
	StatusWeighed:     {StatusStored, StatusAbnormal},
	StatusStored:      {StatusPicked, StatusAbnormal},
	StatusPicked:      {StatusPacked},
	StatusPacked:      {StatusOutbound},
	StatusOutbound:    {StatusContainerArea},
	StatusContainerArea: {StatusLoaded},
	StatusLoaded:      {StatusShipped},
	StatusShipped:     {StatusCustoms},
	StatusCustoms:     {StatusDelivering},
	StatusDelivering:  {StatusDelivered},
	StatusDelivered:   {}, // terminal
	StatusAbnormal:    {StatusReceived, StatusStored, StatusReturned},
	StatusReturned:    {}, // terminal
}

// CanTransition checks if a status transition is valid
func (sm *ParcelStateMachine) CanTransition(from, to ParcelStatus) bool {
	allowed, ok := ValidParcelTransitions[from]
	if !ok { return false }
	for _, s := range allowed {
		if s == to { return true }
	}
	return false
}

// Transition attempts to transition; returns error if invalid
func (sm *ParcelStateMachine) Transition(p *Parcel, to ParcelStatus) error {
	if !sm.CanTransition(p.Status, to) {
		return &InvalidTransitionError{From: p.Status, To: to}
	}
	p.Status = to
	return nil
}

// MustTransition panics if invalid (use only when caller already validated)
func (sm *ParcelStateMachine) MustTransition(p *Parcel, to ParcelStatus) {
	if err := sm.Transition(p, to); err != nil {
		panic(err)
	}
}

// RequiredFields returns fields required for each status transition
func RequiredFieldsForTransition(from, to ParcelStatus) map[string]string {
	switch {
	case from == StatusPreDeclared && to == StatusReceived:
		return map[string]string{"weight": "required", "length": "required", "width": "required", "height": "required"}
	case (from == StatusReceived || from == StatusWeighed) && to == StatusStored:
		return map[string]string{"location_barcode": "required"}
	case from == StatusStored && to == StatusPicked:
		return map[string]string{"order_id": "required"}
	default:
		return nil
	}
}

// StatusDisplay returns user-friendly Chinese display name
func StatusDisplay(s ParcelStatus) string {
	switch s {
	case StatusPreDeclared: return "已预报"
	case StatusReceived: return "已入库"
	case StatusWeighed: return "已核重"
	case StatusStored: return "已上架"
	case StatusPicked: return "已拣货"
	case StatusPacked: return "已打包"
	case StatusOutbound: return "已出库"
	case StatusContainerArea: return "已送装柜区"
	case StatusLoaded: return "已装柜"
	case StatusShipped: return "运输中"
	case StatusCustoms: return "清关中"
	case StatusDelivering: return "派送中"
	case StatusDelivered: return "已签收"
	case StatusAbnormal: return "异常"
	case StatusReturned: return "已退件"
	default: return string(s)
	}
}

// InvalidTransitionError is returned when a status transition is blocked
type InvalidTransitionError struct {
	From, To ParcelStatus
}

func (e *InvalidTransitionError) Error() string {
	return "invalid parcel status transition: " + string(e.From) + " → " + string(e.To)
}

// NewParcelStateMachine creates a new state machine
func NewParcelStateMachine() *ParcelStateMachine {
	return &ParcelStateMachine{}
}
