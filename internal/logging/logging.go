package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel tanımları
type LogLevel int

const (
	ERROR LogLevel = iota // 0
	WARN                  // 1
	INFO                  // 2
	DEBUG                 // 3
)

var LogLevelNames = map[LogLevel]string{
	ERROR: "ERROR",
	WARN:  "WARN",
	INFO:  "INFO",
	DEBUG: "DEBUG",
}

var currentLogLevel LogLevel

// SetLogLevel fonksiyonu, ortam değişkeninden log seviyesini ayarlar
func SetLogLevel() {
	lvl := os.Getenv("LOG_LEVEL")
	switch lvl {
	case "ERROR":
		currentLogLevel = ERROR
	case "WARN":
		currentLogLevel = WARN
	case "DEBUG":
		currentLogLevel = DEBUG
	default:
		currentLogLevel = INFO
	}
}

// Genel uygulama logları için logger
var AppLogger *log.Logger

// Audit logları için logger
var AuditLogger *log.Logger

func init() {
	SetLogLevel()

	// Container-friendly: Sadece stdout'a yaz
	// Docker/Kubernetes logları otomatik olarak toplar
	AppLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	AuditLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

// LogApp genel uygulama logları için fonksiyon
func LogApp(level LogLevel, format string, v ...interface{}) {
	if level <= currentLogLevel {
		AppLogger.Printf("[%s] %s", LogLevelNames[level], fmt.Sprintf(format, v...))
	}
}

// AuditLog audit log için struct
type AuditLog struct {
	Timestamp  string      `json:"timestamp"`
	User       string      `json:"user"`
	Endpoint   string      `json:"endpoint"`
	Method     string      `json:"method"`
	StatusCode int         `json:"status_code"`
	ClientIP   string      `json:"client_ip"`
	Params     interface{} `json:"params,omitempty"`
	Message    string      `json:"message,omitempty"`
}

// LogAudit audit logları için fonksiyon (JSON formatında)
func LogAudit(user, endpoint, method string, statusCode int, clientIP string, params interface{}, message string) {
	auditLog := AuditLog{
		Timestamp:  time.Now().Format(time.RFC3339),
		User:       user,
		Endpoint:   endpoint,
		Method:     method,
		StatusCode: statusCode,
		ClientIP:   clientIP,
		Params:     params,
		Message:    message,
	}

	jsonEntry, err := json.Marshal(auditLog)
	if err != nil {
		LogApp(ERROR, "Audit log oluşturulamadı: %v", err)
		return
	}

	AuditLogger.Printf("[AUDIT] %s", string(jsonEntry))
}
