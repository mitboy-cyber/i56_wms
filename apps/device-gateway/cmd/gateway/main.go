// Device Gateway — I56 WMS hardware integration service.
// Manages scales (地磅), conveyors (入库机), and barcode scanners via RS-232/485.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/i56/device-gateway/internal/client"
	"github.com/i56/device-gateway/internal/conveyor"
	"github.com/i56/device-gateway/internal/dispatcher"
	"github.com/i56/device-gateway/internal/scale"
	"github.com/i56/device-gateway/internal/scanner"
	"github.com/i56/device-gateway/internal/session"

	"gopkg.in/yaml.v3"
)

// Config represents the device gateway configuration.
type Config struct {
	Server struct {
		Port   int `yaml:"port"`
		WSPort int `yaml:"ws_port"`
	} `yaml:"server"`

	WMS struct {
		APIURL string `yaml:"api_url"`
		APIKey string `yaml:"api_key"`
	} `yaml:"wms"`

	Devices struct {
		Scales    []DeviceConfig `yaml:"scales"`
		Conveyors []DeviceConfig `yaml:"conveyors"`
		Scanners  []DeviceConfig `yaml:"scanners"`
	} `yaml:"devices"`
}

// DeviceConfig defines a hardware device configuration.
type DeviceConfig struct {
	ID        string `yaml:"id"`
	Port      string `yaml:"port"`
	Baud      int    `yaml:"baud"`
	Protocol  string `yaml:"protocol"` // MODBUS_RTU, CONTINUOUS, TOLEDO
	Warehouse string `yaml:"warehouse"`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	// Load configuration
	configPath := "configs/device-gateway.yaml"
	if envPath := os.Getenv("DEVICE_GATEWAY_CONFIG"); envPath != "" {
		configPath = envPath
	}

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("=== I56 Device Gateway Starting ===")
	log.Printf("Config: %s", configPath)
	log.Printf("Devices: %d scales, %d conveyors, %d scanners",
		len(cfg.Devices.Scales), len(cfg.Devices.Conveyors), len(cfg.Devices.Scanners))

	// Initialize core components
	sessionMgr := session.NewSessionManager()
	wmsClient := client.NewWMSClient(cfg.WMS.APIURL, cfg.WMS.APIKey)

	// Initialize dispatcher
	disp := dispatcher.New(wmsClient, sessionMgr)

	// Initialize devices
	var scales []*scale.Scale
	var conveyors []*conveyor.Conveyor
	var scanners []*scanner.Scanner

	// ── Scales ──
	for _, dc := range cfg.Devices.Scales {
		s := scale.New(dc.ID, dc.Port, dc.Baud, dc.Protocol, sessionMgr)
		if err := s.Connect(); err != nil {
			log.Printf("[scale:%s] connect error: %v", dc.ID, err)
			continue
		}
		s.OnStable(func(weight float64, unit string) {
			disp.OnWeightStable(dc.ID, weight)
		})
		s.OnWeight(func(weight float64, unit string) {
			log.Printf("[scale:%s] reading: %.3f %s", dc.ID, weight, unit)
		})

		// For continuous-mode scales, start read loop
		if dc.Protocol == "CONTINUOUS" || dc.Protocol == "TOLEDO" {
			s.Start()
		}

		scales = append(scales, s)
		log.Printf("[scale:%s] registered: %s@%d protocol=%s warehouse=%s",
			dc.ID, dc.Port, dc.Baud, dc.Protocol, dc.Warehouse)
	}

	// ── Conveyors ──
	for _, dc := range cfg.Devices.Conveyors {
		c := conveyor.New(dc.ID, dc.Port, sessionMgr)
		c.OnBarcodeScan(func(barcode string) {
			disp.OnBarcodeScan(dc.ID, barcode)
			if err := disp.DispatchToConveyor(c, dc.ID); err != nil {
				log.Printf("[conveyor:%s] dispatch error: %v", dc.ID, err)
			}
		})
		c.OnWeightCapture(func(weight float64) {
			disp.OnWeightStable(dc.ID, weight)
		})
		c.OnArrival(func(location string) {
			disp.OnConveyorArrival(dc.ID, location)
		})

		if err := c.Start(); err != nil {
			log.Printf("[conveyor:%s] start error: %v", dc.ID, err)
			continue
		}

		conveyors = append(conveyors, c)
		log.Printf("[conveyor:%s] registered: %s@%d warehouse=%s",
			dc.ID, dc.Port, dc.Baud, dc.Warehouse)
	}

	// ── Scanners ──
	for _, dc := range cfg.Devices.Scanners {
		s := scanner.New(dc.ID, dc.Port)
		s.OnScan(func(barcode string) {
			disp.OnBarcodeScan(dc.ID, barcode)
		})

		if err := s.Start(); err != nil {
			log.Printf("[scanner:%s] start error: %v", dc.ID, err)
			continue
		}

		scanners = append(scanners, s)
		log.Printf("[scanner:%s] registered: %s", dc.ID, dc.Port)
	}

	// ── Heartbeat goroutine ──
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			for _, s := range scales {
				if err := wmsClient.Heartbeat(s.ID); err != nil {
					log.Printf("[heartbeat] scale %s: %v", s.ID, err)
				} else {
					sessionMgr.Ping(s.ID)
				}
			}
			for _, c := range conveyors {
				if err := wmsClient.Heartbeat(c.ID); err != nil {
					log.Printf("[heartbeat] conveyor %s: %v", c.ID, err)
				} else {
					sessionMgr.Ping(c.ID)
				}
			}
			for _, s := range scanners {
				if err := wmsClient.Heartbeat(s.DeviceID()); err != nil {
					log.Printf("[heartbeat] scanner %s: %v", s.DeviceID(), err)
				} else {
					sessionMgr.Ping(s.DeviceID())
				}
			}
		}
	}()

	// ── HTTP Server ──
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":     "ok",
			"service":    "device-gateway",
			"version":    "1.0.0",
			"devices":    len(scales) + len(conveyors) + len(scanners),
			"scales":     len(scales),
			"conveyors":  len(conveyors),
			"scanners":   len(scanners),
		})
	})

	// Device sessions endpoint
	mux.HandleFunc("GET /api/sessions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sessionMgr.List())
	})

	// Manual scale read
	mux.HandleFunc("GET /api/scale/{id}/read", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		for _, s := range scales {
			if s.ID == id {
				weight, err := s.ReadWeight()
				w.Header().Set("Content-Type", "application/json")
				if err != nil {
					json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
					return
				}
				json.NewEncoder(w).Encode(map[string]any{
					"scale_id": id,
					"weight":   weight,
					"unit":     "kg",
					"stable":   s.IsStable(),
				})
				return
			}
		}
		http.Error(w, "scale not found", http.StatusNotFound)
	})

	// Scale tare
	mux.HandleFunc("POST /api/scale/{id}/tare", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		for _, s := range scales {
			if s.ID == id {
				if err := s.Tare(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"status": "tared", "scale_id": id})
				return
			}
		}
		http.Error(w, "scale not found", http.StatusNotFound)
	})

	// Scale zero
	mux.HandleFunc("POST /api/scale/{id}/zero", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		for _, s := range scales {
			if s.ID == id {
				if err := s.Zero(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{"status": "zeroed", "scale_id": id})
				return
			}
		}
		http.Error(w, "scale not found", http.StatusNotFound)
	})

	serverPort := cfg.Server.Port
	if serverPort == 0 {
		serverPort = 9100
	}
	addr := fmt.Sprintf(":%d", serverPort)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("Device Gateway listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// ── Graceful shutdown ──
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal %v, shutting down...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	// Stop all devices
	for _, s := range scales {
		s.Close()
	}
	for _, c := range conveyors {
		c.Close()
	}
	for _, s := range scanners {
		s.Close()
	}

	log.Println("Device Gateway shutdown complete")
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return &cfg, nil
}
