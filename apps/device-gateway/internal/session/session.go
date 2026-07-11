// Package session manages device connection sessions and tracks real-time state.
package session

import (
	"sync"
	"time"
)

// DeviceSession tracks the real-time state of a connected hardware device.
type DeviceSession struct {
	DeviceID       string    `json:"device_id"`
	Type           string    `json:"type"` // scale, conveyor, scanner
	Status         string    `json:"status"` // connected, disconnected, error
	ConnectedAt    time.Time `json:"connected_at"`
	CurrentBarcode string    `json:"current_barcode,omitempty"`
	CurrentWeight  float64   `json:"current_weight"`
	CurrentUnit    string    `json:"current_unit,omitempty"`
	LastActivity   time.Time `json:"last_activity"`
	ErrorMsg       string    `json:"error_msg,omitempty"`
}

// SessionManager manages all device sessions.
type SessionManager struct {
	mu       sync.RWMutex
	sessions map[string]*DeviceSession
}

// NewSessionManager creates a new session manager.
func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*DeviceSession),
	}
}

// Register adds or updates a device session.
func (sm *SessionManager) Register(deviceID, deviceType string) *DeviceSession {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if existing, ok := sm.sessions[deviceID]; ok {
		existing.Status = "connected"
		existing.LastActivity = time.Now()
		existing.ErrorMsg = ""
		return existing
	}

	session := &DeviceSession{
		DeviceID:    deviceID,
		Type:        deviceType,
		Status:      "connected",
		ConnectedAt: time.Now(),
	}
	sm.sessions[deviceID] = session
	return session
}

// Disconnect marks a device as disconnected.
func (sm *SessionManager) Disconnect(deviceID string, reason string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[deviceID]; ok {
		s.Status = "disconnected"
		s.ErrorMsg = reason
		s.LastActivity = time.Now()
	}
}

// SetError marks a device with an error.
func (sm *SessionManager) SetError(deviceID, errMsg string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[deviceID]; ok {
		s.Status = "error"
		s.ErrorMsg = errMsg
		s.LastActivity = time.Now()
	}
}

// UpdateWeight updates the current weight for a device session.
func (sm *SessionManager) UpdateWeight(deviceID string, weight float64, unit string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[deviceID]; ok {
		s.CurrentWeight = weight
		s.CurrentUnit = unit
		s.LastActivity = time.Now()
	}
}

// UpdateBarcode updates the current barcode for a device session.
func (sm *SessionManager) UpdateBarcode(deviceID, barcode string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[deviceID]; ok {
		s.CurrentBarcode = barcode
		s.LastActivity = time.Now()
	}
}

// Get returns a device session by ID.
func (sm *SessionManager) Get(deviceID string) *DeviceSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.sessions[deviceID]
}

// List returns all active sessions.
func (sm *SessionManager) List() []*DeviceSession {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	result := make([]*DeviceSession, 0, len(sm.sessions))
	for _, s := range sm.sessions {
		result = append(result, s)
	}
	return result
}

// PruneStale removes sessions inactive for longer than the given duration.
func (sm *SessionManager) PruneStale(maxAge time.Duration) int {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	cutoff := time.Now().Add(-maxAge)
	removed := 0
	for id, s := range sm.sessions {
		if s.LastActivity.Before(cutoff) && s.Status == "disconnected" {
			delete(sm.sessions, id)
			removed++
		}
	}
	return removed
}

// Ping heartbeat updates the last activity timestamp.
func (sm *SessionManager) Ping(deviceID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.sessions[deviceID]; ok {
		s.LastActivity = time.Now()
		if s.Status == "disconnected" {
			s.Status = "connected"
		}
	}
}
