// Package scale provides weighbridge (地磅) driver for RS-232/485 scales.
package scale

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/i56/device-gateway/internal/protocol"
	"github.com/i56/device-gateway/internal/session"
)

// deadlineReadWriter extends io.ReadWriteCloser with SetReadDeadline capability
// (typically available on net.Conn implementations like serial ports).
type deadlineReadWriter interface {
	io.ReadWriteCloser
	SetReadDeadline(t time.Time) error
}

// Scale represents a weighbridge device (RS-232/485 serial scale).
type Scale struct {
	ID       string
	Port     string // e.g., /dev/ttyUSB0
	BaudRate int    // 9600, 19200, etc.
	Protocol string // "MODBUS_RTU" | "CONTINUOUS" | "TOLEDO" | "CUSTOM"

	conn    deadlineReadWriter
	adapter protocol.ProtocolAdapter

	mu            sync.RWMutex
	lastWeight    float64
	lastUnit      string
	stableWeight  float64
	isStable      bool
	stableCallback func(weight float64, unit string)
	weightCallback func(weight float64, unit string)
	stopCh        chan struct{}
	running       bool
	sessionMgr    *session.SessionManager
}

// New creates a new Scale driver.
func New(id, port string, baudRate int, protoName string, sm *session.SessionManager) *Scale {
	return &Scale{
		ID:         id,
		Port:       port,
		BaudRate:   baudRate,
		Protocol:   protoName,
		stopCh:     make(chan struct{}),
		sessionMgr: sm,
	}
}

// Connect initializes the protocol adapter. The caller must set the connection
// via SetConnection after opening the serial port.
func (s *Scale) Connect() error {
	var adapter protocol.ProtocolAdapter
	switch s.Protocol {
	case "MODBUS_RTU":
		adapter = protocol.NewModbusRTUAdapter(0x01)
	case "CONTINUOUS":
		adapter = protocol.NewContinuousAdapter()
	case "TOLEDO":
		adapter = protocol.NewToledoAdapter()
	case "CUSTOM":
		adapter = protocol.NewCustomAdapter("CUSTOM", []byte{0x02}, []byte{0x03}, 1, 8)
	default:
		return fmt.Errorf("scale: unknown protocol %q", s.Protocol)
	}

	s.mu.Lock()
	s.adapter = adapter
	s.mu.Unlock()

	if s.sessionMgr != nil {
		s.sessionMgr.Register(s.ID, "scale")
	}

	log.Printf("[scale:%s] initialized with protocol %s on %s@%d",
		s.ID, adapter.Name(), s.Port, s.BaudRate)
	return nil
}

// SetConnection sets the I/O connection after serial port opening.
func (s *Scale) SetConnection(conn deadlineReadWriter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conn = conn
}

// OnStable registers a callback that fires when weight stabilizes.
func (s *Scale) OnStable(callback func(weight float64, unit string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stableCallback = callback
}

// OnWeight registers a callback that fires on every weight reading.
func (s *Scale) OnWeight(callback func(weight float64, unit string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weightCallback = callback
}

// ReadWeight sends a read command and returns the current weight.
func (s *Scale) ReadWeight() (float64, error) {
	s.mu.RLock()
	adapter := s.adapter
	conn := s.conn
	s.mu.RUnlock()

	if adapter == nil {
		return 0, fmt.Errorf("scale %s: not connected", s.ID)
	}
	if conn == nil {
		return 0, fmt.Errorf("scale %s: no serial connection", s.ID)
	}

	// Send read command
	cmd, err := adapter.Encode(protocol.Command{Code: 'R'})
	if err != nil {
		// Some adapters (Modbus) use specific read commands
		cmd, err = adapter.Encode(protocol.Command{Code: 0x00, Data: []byte{0x00, 0x00, 0x01}})
		if err != nil {
			return 0, fmt.Errorf("scale %s: encode error: %w", s.ID, err)
		}
	}

	if _, err := conn.Write(cmd); err != nil {
		return 0, fmt.Errorf("scale %s: write error: %w", s.ID, err)
	}

	// Read response
	buf := make([]byte, 256)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	n, err := conn.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("scale %s: read error: %w", s.ID, err)
	}

	val, err := adapter.Parse(buf[:n])
	if err != nil {
		return 0, fmt.Errorf("scale %s: parse error: %w", s.ID, err)
	}

	weight, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("scale %s: unexpected parse result type %T", s.ID, val)
	}

	s.mu.Lock()
	s.lastWeight = weight
	s.lastUnit = "kg"
	s.mu.Unlock()

	if s.sessionMgr != nil {
		s.sessionMgr.UpdateWeight(s.ID, weight, "kg")
	}

	return weight, nil
}

// Tare sends a tare (去皮) command to the scale.
func (s *Scale) Tare() error {
	return s.sendCommand('T')
}

// Zero sends a zero (归零) command to the scale.
func (s *Scale) Zero() error {
	return s.sendCommand('Z')
}

func (s *Scale) sendCommand(code byte) error {
	s.mu.RLock()
	adapter := s.adapter
	conn := s.conn
	s.mu.RUnlock()

	if adapter == nil {
		return fmt.Errorf("scale %s: not connected", s.ID)
	}
	if conn == nil {
		return fmt.Errorf("scale %s: no serial connection", s.ID)
	}

	cmd, err := adapter.Encode(protocol.Command{Code: code})
	if err != nil {
		return fmt.Errorf("scale %s: command encode error: %w", s.ID, err)
	}

	if _, err := conn.Write(cmd); err != nil {
		return fmt.Errorf("scale %s: write error: %w", s.ID, err)
	}

	log.Printf("[scale:%s] command 0x%02X sent", s.ID, code)
	return nil
}

// Start begins the continuous read loop for scales that stream data.
func (s *Scale) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	go s.readLoop()
}

// Stop stops the read loop.
func (s *Scale) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	close(s.stopCh)
}

// readLoop continuously reads weight data from the serial port.
func (s *Scale) readLoop() {
	log.Printf("[scale:%s] read loop started (protocol=%s)", s.ID, s.Protocol)

	const stabilityWindow = 3 // number of consistent readings to consider "stable"
	const stabilityTolerance = 0.1 // kg tolerance
	readings := make([]float64, 0, stabilityWindow)

	for {
		select {
		case <-s.stopCh:
			log.Printf("[scale:%s] read loop stopped", s.ID)
			return
		default:
		}

		s.mu.RLock()
		adapter := s.adapter
		conn := s.conn
		s.mu.RUnlock()

		if adapter == nil || conn == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		buf := make([]byte, 256)
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Printf("[scale:%s] read error: %v", s.ID, err)
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if n == 0 {
			continue
		}

		val, err := adapter.Parse(buf[:n])
		if err != nil {
			log.Printf("[scale:%s] parse error: %v (raw=%q)", s.ID, err, buf[:n])
			continue
		}

		weight, ok := val.(float64)
		if !ok {
			continue
		}

		// Update last known weight
		s.mu.Lock()
		s.lastWeight = weight
		s.lastUnit = "kg"
		wc := s.weightCallback
		s.mu.Unlock()

		// Update session
		if s.sessionMgr != nil {
			s.sessionMgr.UpdateWeight(s.ID, weight, "kg")
		}

		// Fire weight callback
		if wc != nil {
			wc(weight, "kg")
		}

		// Stability detection via moving window
		readings = append(readings, weight)
		if len(readings) > stabilityWindow {
			readings = readings[1:]
		}

		if len(readings) == stabilityWindow {
			min, max := readings[0], readings[0]
			for _, r := range readings[1:] {
				if r < min {
					min = r
				}
				if r > max {
					max = r
				}
			}

			if max-min <= stabilityTolerance {
				stableWeight := (min + max) / 2.0
				s.mu.Lock()
				wasStable := s.isStable
				s.stableWeight = stableWeight
				s.isStable = true
				sc := s.stableCallback
				s.mu.Unlock()

				if !wasStable && sc != nil {
					log.Printf("[scale:%s] weight stabilized at %.3f kg", s.ID, stableWeight)
					sc(stableWeight, "kg")
				}
			} else {
				s.mu.Lock()
				s.isStable = false
				s.mu.Unlock()
			}
		}
	}
}

// LastWeight returns the most recent weight reading.
func (s *Scale) LastWeight() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastWeight
}

// IsStable returns whether the weight is currently stable.
func (s *Scale) IsStable() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isStable
}

// Close closes the serial connection and stops the read loop.
func (s *Scale) Close() error {
	s.Stop()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
