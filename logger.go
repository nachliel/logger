package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type any = interface{}

// Logger struct is the main structure of CSV-ES Logger
type Logger struct {
	level        Level
	msgCounter   int64
	es           *elasticsearch.Client
	proccessName string
	indexName    string
}

// Level of error
type Level int

// JLog is the Json structure of log
type JLog struct {
	ProccessName string    `json:"proccess"`
	StatusTime   time.Time `json:"status_time"`
	Level        string    `json:"level"`
	Message      string    `json:"message"`
}

// Const Name of Log Levels
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

var Settings Logger

// NewLogger create a new Logger
func SetLevel(level Level) {
	Settings.level = level
	Settings.msgCounter = 0
	Settings.es = nil
}

//	Add ElasticSearch Client to log messages there
func addElasticClient(esClient *elasticsearch.Client, index string) {
	Settings.es = esClient
	Settings.indexName = index
}

//	Debug Logging
func Debug(format string, args ...any) {
	write(LevelDebug, format, args...)
}

// Info Logging
func Info(format string, args ...any) {
	write(LevelInfo, format, args...)
}

// Warn Logging
func Warn(format string, args ...any) {
	write(LevelWarn, format, args...)
}

// Error Logging
func Error(format string, args ...any) {
	write(LevelError, format, args...)
}

// Fatal Logging
func Fatal(format string, args ...any) {
	write(LevelFatal, format, args...)
	// Exit Failure
	os.Exit(1)
}

// write Log to console, if ES is setup then also to ES.
func write(level Level, format string, args ...any) {
	// Log Message counter
	Settings.msgCounter++

	// Check level.. ask Levi... i dunno
	if level < Settings.level {
		return
	}
	// Output to Terminal
	fmt.Printf("%s [%s] ", time.Now().Format(time.RFC822), level)
	fmt.Printf(format, args...)
	fmt.Printf("\n")
	if Settings.es != nil {
		// Write Log to ES
		jsonElementStruct := JLog{
			StatusTime: time.Now().UTC(),
			Level:      level.String(),
			Message:    fmt.Sprintf(format, args...),
		}
		docReader := esutil.NewJSONReader(&jsonElementStruct)

		res, err := Settings.es.Index(
			Settings.indexName, // Index name
			docReader,          // Document body
		)
		if err != nil {
			Settings.msgCounter++
			fmt.Printf("%s [%s] ", time.Now().Format(time.RFC822), LevelFatal)
			fmt.Printf("ElasticSearch Index Update Error: %s\n", err)
		}
		res.Body.Close()
	}
}

// Converts Level to string.
func (level Level) String() string {
	switch level {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKN"
	}
}
