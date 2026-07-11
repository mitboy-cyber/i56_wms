// Package conveyor provides inbound conveyor (入库机) driver with task dispatch.
package conveyor

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/i56/device-gateway/internal/session"
)

// InboundTask represents an inbound task assigned to the conveyor.
type InboundTask struct {
	TaskID         string  `json:"task_id"`
	WaybillNo      string  `json:"waybill_no"`
	TrackingNumber string  `json:"tracking_number"`
	Barcode        string  `json:"barcode"`
	SKUCode        string  `json:"sku_code"`
	ProductName    string  `json:"product_name"`
	TargetLocation string  `json:"target_location"` // chute/lane ID
	DeclaredWeight float64 `json:"declared_weight"`
}

// Conveyor represents an inbound conveyor system with integrated barcode scanner and scale.
type Conveyor struct {
	ID   string
	Port string

	conn   io.ReadWriteCloser
	reader *bufio.Reader

	mu              sync.RWMutex
	currentBarcode  string
	currentWeight   float64
	barcodeCallback func(barcode string)
	weightCallback  func(weight float64)
	arrivalCallback func(location string)
	stopCh          chan struct{}
	running         bool
	sessionMgr      *session.SessionManager
}

// New creates a new Conveyor driver.
func New(id, port string, sm *session.SessionManager) *Conveyor {
	return &Conveyor{
		ID:         id,
		Port:       port,
		stopCh:     make(chan struct{}),
		sessionMgr: sm,
	}
}

// SetConnection sets the I/O connection after serial port opening.
func (c *Conveyor) SetConnection(conn io.ReadWriteCloser) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.conn = conn
	c.reader = bufio.NewReader(conn)
}

// OnBarcodeScan registers a callback for barcode scan events on the conveyor.
func (c *Conveyor) OnBarcodeScan(callback func(barcode string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.barcodeCallback = callback
}

// OnWeightCapture registers a callback for weight capture events on the conveyor.
func (c *Conveyor) OnWeightCapture(callback func(weight float64)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.weightCallback = callback
}

// OnArrival registers a callback for package arrival at a destination.
func (c *Conveyor) OnArrival(callback func(location string)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.arrivalCallback = callback
}

// Dispatch sends a dispatch command for an inbound task (分拨到滑道/货架).
func (c *Conveyor) Dispatch(task InboundTask) error {
	c.mu.RLock()
	conn := c.conn
	c.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("conveyor %s: not connected", c.ID)
	}

	// Protocol: "DISPATCH:<location>\r\n"
	cmd := fmt.Sprintf("DISPATCH:%s\r\n", task.TargetLocation)
	if _, err := conn.Write([]byte(cmd)); err != nil {
		return fmt.Errorf("conveyor %s: dispatch write error: %w", c.ID, err)
	}

	c.mu.Lock()
	c.currentBarcode = task.Barcode
	c.mu.Unlock()

	if c.sessionMgr != nil {
		c.sessionMgr.UpdateBarcode(c.ID, task.Barcode)
	}

	log.Printf("[conveyor:%s] dispatched task %s to %s", c.ID, task.TaskID, task.TargetLocation)
	return nil
}

// Start begins the read loop for conveyor events (barcode scans, weight captures, arrivals).
func (c *Conveyor) Start() error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = true
	c.mu.Unlock()

	if c.sessionMgr != nil {
		c.sessionMgr.Register(c.ID, "conveyor")
	}

	go c.readLoop()
	return nil
}

// Stop stops the conveyor read loop.
func (c *Conveyor) Stop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.running {
		return
	}
	c.running = false
	close(c.stopCh)
}

// readLoop continuously reads events from the conveyor serial port.
func (c *Conveyor) readLoop() {
	log.Printf("[conveyor:%s] read loop started", c.ID)

	const (
		prefixBarcode = "BARCODE:"
		prefixWeight  = "WEIGHT:"
		prefixArrival = "ARRIVED:"
		prefixStatus  = "STATUS:"
	)

	for {
		select {
		case <-c.stopCh:
			log.Printf("[conveyor:%s] read loop stopped", c.ID)
			return
		default:
		}

		c.mu.RLock()
		reader := c.reader
		c.mu.RUnlock()

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
			log.Printf("[conveyor:%s] read error: %v", c.ID, err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse event type
		switch {
		case strings.HasPrefix(line, prefixBarcode):
			barcode := strings.TrimPrefix(line, prefixBarcode)
			barcode = strings.TrimSpace(barcode)
			log.Printf("[conveyor:%s] barcode scan: %s", c.ID, barcode)

			c.mu.Lock()
			c.currentBarcode = barcode
			cb := c.barcodeCallback
			c.mu.Unlock()

			if c.sessionMgr != nil {
				c.sessionMgr.UpdateBarcode(c.ID, barcode)
			}
			if cb != nil {
				cb(barcode)
			}

		case strings.HasPrefix(line, prefixWeight):
			weightStr := strings.TrimPrefix(line, prefixWeight)
			weightStr = strings.TrimSpace(weightStr)
			var weight float64
			if _, err := fmt.Sscanf(weightStr, "%f", &weight); err != nil {
				log.Printf("[conveyor:%s] invalid weight: %q", c.ID, weightStr)
				continue
			}

			c.mu.Lock()
			c.currentWeight = weight
			wc := c.weightCallback
			c.mu.Unlock()

			if c.sessionMgr != nil {
				c.sessionMgr.UpdateWeight(c.ID, weight, "kg")
			}
			if wc != nil {
				wc(weight)
			}

		case strings.HasPrefix(line, prefixArrival):
			location := strings.TrimPrefix(line, prefixArrival)
			location = strings.TrimSpace(location)
			log.Printf("[conveyor:%s] arrival at: %s", c.ID, location)

			c.mu.RLock()
			ac := c.arrivalCallback
			c.mu.RUnlock()

			if ac != nil {
				ac(location)
			}

		case strings.HasPrefix(line, prefixStatus):
			status := strings.TrimPrefix(line, prefixStatus)
			log.Printf("[conveyor:%s] status: %s", c.ID, strings.TrimSpace(status))

		default:
			log.Printf("[conveyor:%s] unknown event: %q", c.ID, line)
		}
	}
}

// CurrentBarcode returns the last scanned barcode.
func (c *Conveyor) CurrentBarcode() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentBarcode
}

// CurrentWeight returns the last captured weight.
func (c *Conveyor) CurrentWeight() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentWeight
}

// Close closes the serial connection and stops the read loop.
func (c *Conveyor) Close() error {
	c.Stop()
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
