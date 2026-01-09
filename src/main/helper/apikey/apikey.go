package apikey

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/cache"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/internal/shared/response"
	"github.com/Sistem-Informasi-Akademik/academic-system-information-service/src/main/model/entity/school"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	// CacheKeyPrefix prefix untuk API key cache
	CacheKeyPrefix = "apikey"
	// CacheTTL durasi cache (1 jam)
	CacheTTL = time.Hour
)

// ApiKeyCache struktur untuk cache API key validation
type ApiKeyCache struct {
	SchoolCode string `json:"school_code"`
	ClientID   string `json:"client_id"`
	Origin     string `json:"origin"`
}

// global DB reference for middleware (set during app initialization)
var (
	globalDB   *gorm.DB
	globalOnce sync.Once
)

// SetDB sets the database connection for the apikey middleware
// This should be called once during application initialization
func SetDB(db *gorm.DB) {
	globalOnce.Do(func() {
		globalDB = db
	})
}

// buildCacheKey generates cache key from client_id and origin
func buildCacheKey(clientID, origin string) string {
	return fmt.Sprintf("%s:%s:%s", CacheKeyPrefix, clientID, origin)
}

// ApiKeyMiddleware validates client_id from header against origin
// This middleware should be used on public endpoints that require API key validation
// such as login endpoints
//
// Required Headers:
//   - X-Client-ID: The public client identifier
//   - Origin: The origin domain of the request
//
// On success, sets "api_school_code" in context for use in handlers
//
// Uses Redis cache to avoid hitting a database on every request
func ApiKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := c.GetHeader("X-Client-ID")
		origin := c.Request.Header.Get("Origin")

		// Strip protocol from origin (http:// or https://)
		re := regexp.MustCompile(`^https?://`)
		origin = re.ReplaceAllString(origin, "")

		// Validate client ID exists
		if clientID == "" {
			response.UnauthorizedError(c, "Missing X-Client-ID header")
			c.Abort()
			return
		}

		// Validate origin exists
		if origin == "" {
			response.UnauthorizedError(c, "Missing Origin header")
			c.Abort()
			return
		}

		// Try to get from Redis cache first
		cacheKey := buildCacheKey(clientID, origin)
		var cachedData ApiKeyCache
		err := cache.GetJSON(cacheKey, &cachedData)

		if err == nil && cachedData.SchoolCode != "" {
			// Cache hit - use cached data
			c.Set("api_school_code", cachedData.SchoolCode)
			c.Next()
			return
		}

		// Cache miss - query database
		if globalDB == nil {
			response.InternalError(c, "Database not initialized", nil)
			c.Abort()
			return
		}

		var config school.SystemConfig
		globalDB.Where("client_id = ? AND origin = ?", clientID, origin).First(&config)

		if config.ID == 0 {
			response.UnauthorizedError(c, "Invalid API credentials")
			c.Abort()
			return
		}

		// Store in Redis cache
		cacheData := ApiKeyCache{
			SchoolCode: config.SchoolCode,
			ClientID:   config.ClientId,
			Origin:     config.Origin,
		}
		_ = cache.SetJSON(cacheKey, cacheData, CacheTTL)

		// Set school_code ke context untuk dipakai di handler
		c.Set("api_school_code", config.SchoolCode)
		c.Next()
	}
}

// GetApiSchoolCode returns the school_code from API key validation
// Use this in handlers after ApiKeyMiddleware
func GetApiSchoolCode(c *gin.Context) string {
	schoolCode, exists := c.Get("api_school_code")
	if !exists {
		return ""
	}
	return schoolCode.(string)
}

// InvalidateApiKeyCache removes the cache for a specific client_id and origin
// Call this when SystemConfig is updated
func InvalidateApiKeyCache(clientID, origin string) error {
	cacheKey := buildCacheKey(clientID, origin)
	return cache.RemoveByKey(cacheKey)
}

// InvalidateAllApiKeyCache removes all API key caches
// Call this when doing bulk updates to SystemConfig
func InvalidateAllApiKeyCache() error {
	return cache.DeleteByPrefix(CacheKeyPrefix)
}
