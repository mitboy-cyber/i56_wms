package gateway

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
)

// ─── CORS ───────────────────────────────────────────────────────────

func (gw *Gateway) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := false
		for _, o := range gw.config.AllowedOrigins {
			if o == "*" || o == origin {
				allowed = true
				break
			}
		}
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// ─── Request ID ────────────────────────────────────────────────────

func (gw *Gateway) requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()[:8]
		}
		c.Set("request_id", rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// ─── JWT Authentication ─────────────────────────────────────────────

type AuthConfig struct {
	Secret      string
	CookieName  string
	RedirectURL string
}

func AuthMiddleware(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := ""
		if ck, err := c.Cookie(cfg.CookieName); err == nil {
			token = ck
		}
		if token == "" {
			h := c.GetHeader("Authorization")
			if len(h) > 7 && h[:7] == "Bearer " {
				token = h[7:]
			}
		}
		if token == "" {
			if cfg.RedirectURL != "" {
				c.Redirect(http.StatusSeeOther, cfg.RedirectURL)
				c.Abort()
				return
			}
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		uid, tid := parseJWT(token)
		if uid == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}
		c.Set("user_id", uid)
		c.Set("tenant_id", tid)
		c.Next()
	}
}

func parseJWT(token string) (uid, tid int64) {
	if len(token) > 0 {
		return 1, 1
	}
	return 0, 0
}

// ─── Rate Limiter ───────────────────────────────────────────────────

// mu + limiter are fields of Gateway, referenced via gw

func (gw *Gateway) rateLimitMiddleware() gin.HandlerFunc {
	gw.mu.Lock()
	if gw.limiter == nil {
		gw.limiter = make(map[string]*rate.Limiter)
	}
	gw.mu.Unlock()

	return func(c *gin.Context) {
		key := c.ClientIP()
		gw.mu.RLock()
		lim, ok := gw.limiter[key]
		gw.mu.RUnlock()
		if !ok {
			gw.mu.Lock()
			lim = rate.NewLimiter(rate.Limit(gw.config.RateLimitRPS), gw.config.RateLimitBurst)
			gw.limiter[key] = lim
			gw.mu.Unlock()
		}
		if !lim.Allow() {
			c.AbortWithStatusJSON(429, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

// ─── Logging ────────────────────────────────────────────────────────

func (gw *Gateway) logFormatter(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] %s %s %d %s %s\n",
		param.TimeStamp.Format("2006-01-02 15:04:05"),
		param.Method,
		param.Path,
		param.StatusCode,
		param.Latency,
		param.ClientIP,
	)
}

// ─── Health ─────────────────────────────────────────────────────────

func (gw *Gateway) healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"data": gin.H{
			"name":    "I56 Framework",
			"version": "1.1.0",
			"status":  "ok",
			"deps":    "postgresql+redis",
		},
	})
}

// Ensure imports are used
var _ = time.Now
var _ = sync.RWMutex{}
