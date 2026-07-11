// Package dispatcher orchestrates device events into WMS business logic.
package dispatcher

import (
	"fmt"
	"log"
	"sync"

	"github.com/i56/device-gateway/internal/client"
	"github.com/i56/device-gateway/internal/conveyor"
	"github.com/i56/device-gateway/internal/session"
)

// Dispatcher coordinates barcode scans, weight captures, and conveyor arrivals
// with the WMS API to execute inbound task workflows.
type Dispatcher struct {
	wmsClient   *client.WMSClient
	sessionMgr  *session.SessionManager

	mu          sync.RWMutex
	currentTask map[string]*client.InboundTask // keyed by deviceID
	weightTolerance float64 // allowed weight difference ratio (default 0.05 = 5%)
}

// New creates a new Dispatcher.
func New(wmsClient *client.WMSClient, sm *session.SessionManager) *Dispatcher {
	return &Dispatcher{
		wmsClient:       wmsClient,
		sessionMgr:      sm,
		currentTask:     make(map[string]*client.InboundTask),
		weightTolerance: 0.05,
	}
}

// OnBarcodeScan handles a barcode scan event from any device (scanner or conveyor).
func (d *Dispatcher) OnBarcodeScan(deviceID, barcode string) {
	log.Printf("[dispatcher] barcode scan: device=%s barcode=%s", deviceID, barcode)

	// 1. Query WMS for inbound task
	task, err := d.wmsClient.GetInboundTaskByBarcode(barcode)
	if err != nil {
		log.Printf("[dispatcher] WMS query failed for barcode %q: %v", barcode, err)
		if d.sessionMgr != nil {
			d.sessionMgr.SetError(deviceID, "WMS query failed: "+err.Error())
		}
		return
	}

	d.mu.Lock()
	d.currentTask[deviceID] = task
	d.mu.Unlock()

	if d.sessionMgr != nil {
		d.sessionMgr.UpdateBarcode(deviceID, barcode)
		d.sessionMgr.Register(deviceID, d.sessionMgr.Get(deviceID).Type)
	}

	log.Printf("[dispatcher] task found: waybill=%s product=%s declaredWeight=%.3f kg location=%s",
		task.WaybillNo, task.ProductName, task.DeclaredWeight, task.LocationCode)
}

// OnWeightStable handles a weight stabilization event from a scale.
func (d *Dispatcher) OnWeightStable(deviceID string, weight float64) {
	log.Printf("[dispatcher] weight stable: device=%s weight=%.3f kg", deviceID, weight)

	d.mu.RLock()
	task := d.currentTask[deviceID]
	d.mu.RUnlock()

	if task == nil {
		log.Printf("[dispatcher] no current task for device %s — standalone weigh", deviceID)
		// Record as standalone weight
		_, err := d.wmsClient.RecordWeight("", weight, deviceID)
		if err != nil {
			log.Printf("[dispatcher] record weight error: %v", err)
		}
		return
	}

	// 2. Record weight to WMS
	record, err := d.wmsClient.RecordWeight(task.WaybillNo, weight, deviceID)
	if err != nil {
		log.Printf("[dispatcher] record weight error: %v", err)
		if d.sessionMgr != nil {
			d.sessionMgr.SetError(deviceID, "Weight record failed: "+err.Error())
		}
		return
	}

	// 3. Validate against declared weight
	if task.DeclaredWeight > 0 {
		diff := weight - task.DeclaredWeight
		diffRatio := diff / task.DeclaredWeight
		if diffRatio < 0 {
			diffRatio = -diffRatio
		}

		if diffRatio > d.weightTolerance {
			log.Printf("[dispatcher] WEIGHT MISMATCH: waybill=%s declared=%.3f actual=%.3f diff=%.1f%%",
				task.WaybillNo, task.DeclaredWeight, weight, diffRatio*100)
			if d.sessionMgr != nil {
				d.sessionMgr.SetError(deviceID,
					fmt.Sprintf("Weight mismatch: expected %.3f, actual %.3f (diff %.1f%%)",
						task.DeclaredWeight, weight, diffRatio*100))
			}
			// Continue — flag as abnormal but don't block
		}
	}

	log.Printf("[dispatcher] weight recorded: waybill=%s weight=%.3f kg (recordID=%d)",
		task.WaybillNo, weight, record.ID)
}

// OnConveyorArrival handles a package arrival at a destination location.
func (d *Dispatcher) OnConveyorArrival(deviceID, location string) {
	log.Printf("[dispatcher] arrival: device=%s location=%s", deviceID, location)

	d.mu.RLock()
	task := d.currentTask[deviceID]
	d.mu.RUnlock()

	if task == nil {
		log.Printf("[dispatcher] no current task for conveyor %s at location %s", deviceID, location)
		return
	}

	// Confirm inbound with location
	if err := d.wmsClient.ConfirmInbound(task.WaybillNo, location); err != nil {
		log.Printf("[dispatcher] confirm inbound error: %v", err)
		if d.sessionMgr != nil {
			d.sessionMgr.SetError(deviceID, "Confirm inbound failed: "+err.Error())
		}
		return
	}

	log.Printf("[dispatcher] inbound confirmed: waybill=%s location=%s", task.WaybillNo, location)

	// Clear current task
	d.mu.Lock()
	delete(d.currentTask, deviceID)
	d.mu.Unlock()
}

// DispatchToConveyor sends a task to a conveyor for routing.
func (d *Dispatcher) DispatchToConveyor(conv *conveyor.Conveyor, deviceID string) error {
	d.mu.RLock()
	task := d.currentTask[deviceID]
	d.mu.RUnlock()

	if task == nil {
		return nil // no task to dispatch
	}

	convTask := conveyor.InboundTask{
		TaskID:         task.WaybillNo,
		WaybillNo:      task.WaybillNo,
		TrackingNumber: task.TrackingNumber,
		Barcode:        task.TrackingNumber,
		SKUCode:        task.SKUCode,
		ProductName:    task.ProductName,
		TargetLocation: task.LocationCode,
		DeclaredWeight: task.DeclaredWeight,
	}

	return conv.Dispatch(convTask)
}
