package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	Env             string
	DBDriver        string // "postgres" or "sqlite"
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	DBSQLitePath    string
	JWTSecret       string
	JWTExpiryHours  int
	LiveKitURL      string
	LiveKitAPIKey   string
	LiveKitSecret   string
	UploadDir       string
	CORSAllowedUrls []string
}

func Load() *Config {
	_ = godotenv.Load(".env", "../.env")

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		Env:             getEnv("APP_ENV", "development"),
		DBDriver:        getEnv("DB_DRIVER", "postgres"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "lms_db"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		DBSQLitePath:    getEnv("DB_SQLITE_PATH", "./lms.db"),
		JWTSecret:       getEnv("JWT_SECRET", "lms-super-secure-production-jwt-key-32chars"),
		JWTExpiryHours:  getEnvAsInt("JWT_EXPIRY_HOURS", 24),
		LiveKitURL:      getEnv("LIVEKIT_URL", "ws://localhost:7880"),
		LiveKitAPIKey:   getEnv("LIVEKIT_API_KEY", "devkey"),
		LiveKitSecret:   getEnv("LIVEKIT_API_SECRET", "secret"),
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads"),
		CORSAllowedUrls: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:5173"}),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	parts := strings.Split(valStr, ",")
	var res []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
