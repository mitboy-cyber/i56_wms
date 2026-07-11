// Package scanner provides barcode scanner device driver.
package scanner

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"
	"time"
)

// Scanner represents a barcode scanner device (RS-232 or USB-Serial).
type Scanner struct {
	deviceID string
	port     string
	conn     io.ReadWriteCloser
	reader   *bufio.Reader

	mu       sync.RWMutex
	callback func(barcode string)
	stopCh   chan struct{}
	running  bool
}

// New creates a new barcode scanner driver.
func New(deviceID, port string) *Scanner {
	return &Scanner{
		deviceID: deviceID,
		port:     port,
		stopCh:   make(chan struct{}),
	}
}

// DeviceID returns the scanner's identifier.
func (s *Scanner) DeviceID() string { return s.deviceID }

// Port returns the scanner's serial port.
func (s *Scanner) Port() string { return s.port }

// SetConnection sets the I/O connection (used after serial port open).
func (s *Scanner) SetConnection(conn io.ReadWriteCloser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.conn = conn
	s.reader = bufio.NewReader(conn)
}

// OnScan registers a callback that fires when a barcode is scanned.
// Barcode scanners typically append CR/LF after the barcode.
func (s *Scanner) OnScan(callback func(barcode string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callback = callback
}

// Start begins listening for barcode scans on the serial port.
func (s *Scanner) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	go s.readLoop()
	return nil
}

// Stop stops the barcode scanner read loop.
func (s *Scanner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.running {
		return
	}
	s.running = false
	close(s.stopCh)
}

// readLoop continuously reads from the serial port and fires callbacks.
func (s *Scanner) readLoop() {
	log.Printf("[scanner:%s] read loop started", s.deviceID)

	for {
		select {
		case <-s.stopCh:
			log.Printf("[scanner:%s] read loop stopped", s.deviceID)
			return
		default:
		}

		s.mu.RLock()
		reader := s.reader
		s.mu.RUnlock()

		if reader == nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			log.Printf("[scanner:%s] read error: %v", s.deviceID, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		barcode := strings.TrimSpace(line)
		if barcode == "" {
			continue
		}

		// Validate barcode: typical barcodes are 8-30 chars of alphanumeric
		if len(barcode) < 4 || len(barcode) > 50 {
			log.Printf("[scanner:%s] ignoring invalid barcode length: %d", s.deviceID, len(barcode))
			continue
		}

		log.Printf("[scanner:%s] scanned: %s", s.deviceID, barcode)

		s.mu.RLock()
		cb := s.callback
		s.mu.RUnlock()

		if cb != nil {
			cb(barcode)
		}
	}
}

// Close closes the serial connection.
func (s *Scanner) Close() error {
	s.Stop()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
