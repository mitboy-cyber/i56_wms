package gateway

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type Config struct {
	Port            string
	Mode            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	MaxHeaderBytes  int
	TrustedProxies  []string
	AllowedOrigins  []string
	RateLimitRPS    float64
	RateLimitBurst  int
	JWTSecret       string
	RedisAddr       string
	RedisPassword   string
	EnableWebSocket bool
	EnableSSE       bool
}

type Gateway struct {
	engine  *gin.Engine
	config  Config
	redis   *redis.Client
	wsHub   *WSHub
	sseHub  *SSEHub
	mu      sync.RWMutex
	limiter map[string]*rate.Limiter
}

func New(cfg Config) *Gateway {
	if cfg.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.RateLimitRPS == 0 {
		cfg.RateLimitRPS = 100
	}
	if cfg.RateLimitBurst == 0 {
		cfg.RateLimitBurst = 200
	}
	gw := &Gateway{engine: gin.New(), config: cfg, limiter: make(map[string]*rate.Limiter)}
	gw.engine.Use(gin.LoggerWithFormatter(gw.logFormatter))
	gw.engine.Use(gin.Recovery())
	gw.engine.Use(gw.corsMiddleware())
	gw.engine.Use(gw.requestIDMiddleware())
	gw.engine.GET("/api/v1/health", gw.healthHandler)

	if cfg.RedisAddr != "" {
		gw.redis = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword})
		if err := gw.redis.Ping(context.Background()).Err(); err != nil {
			log.Printf("[GATEWAY] Redis ERR: %v", err)
		} else {
			log.Println("[GATEWAY] Redis OK")
		}
	}
	if cfg.EnableWebSocket {
		gw.wsHub = NewWSHub()
		go gw.wsHub.Run()
		gw.engine.GET("/ws", gw.wsHandler)
	}
	if cfg.EnableSSE {
		gw.sseHub = NewSSEHub()
		gw.engine.GET("/sse", gw.sseHandler)
	}
	return gw
}

func (gw *Gateway) Engine() *gin.Engine        { return gw.engine }
func (gw *Gateway) Redis() *redis.Client        { return gw.redis }
func (gw *Gateway) WSHub() *WSHub               { return gw.wsHub }
func (gw *Gateway) SSEHub() *SSEHub             { return gw.sseHub }

func (gw *Gateway) Serve() error {
	srv := &http.Server{Addr: ":" + gw.config.Port, Handler: gw.engine}
	if gw.config.ReadTimeout > 0 {
		srv.ReadTimeout = gw.config.ReadTimeout
	}
	if gw.config.WriteTimeout > 0 {
		srv.WriteTimeout = gw.config.WriteTimeout
	}
	log.Printf("[GATEWAY] :%s mode=%s", gw.config.Port, gw.config.Mode)
	return srv.ListenAndServe()
}

func (gw *Gateway) Shutdown(ctx context.Context) error {
	if gw.redis != nil {
		gw.redis.Close()
	}
	return nil
}
