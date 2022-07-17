package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type any = interface{}

//
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

// Levels of Logs
const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// NewLogger create a new Logger
func (logger *Logger) SetLevel(level Level) {
	logger.level = level
	logger.msgCounter = 0
	logger.es = nil
}

func (logger *Logger) addElasticClient(esClient *elasticsearch.Client, index string) {
	logger.es = esClient
	logger.indexName = index
}

//	Debug Logging
func (log *Logger) Debug(format string, args ...any) {
	log.write(LevelDebug, format, args...)
}

// Info Logging
func (log *Logger) Info(format string, args ...any) {
	log.write(LevelInfo, format, args...)
}

// Warn Logging
func (log *Logger) Warn(format string, args ...any) {
	log.write(LevelWarn, format, args...)
}

// Error Logging
func (log *Logger) Error(format string, args ...any) {
	log.write(LevelError, format, args...)
}

// Fatal Logging
func (log *Logger) Fatal(format string, args ...any) {
	log.write(LevelFatal, format, args...)
	// Exit Failure
	os.Exit(1)
}

// write Log to ES and console.
func (log *Logger) write(level Level, format string, args ...any) {
	// Log Message counter
	log.msgCounter++

	// Check level.. ask Levi... i dunno
	if level < log.level {
		return
	}
	// Output to Terminal
	fmt.Printf("%s [%s] ", time.Now().Format(time.RFC822), level)
	fmt.Printf(format, args...)
	fmt.Printf("\n")
	if log.es != nil {
		// Write Log to ES
		jsonElementStruct := JLog{
			StatusTime: time.Now().UTC(),
			Level:      level.String(),
			Message:    fmt.Sprintf(format, args...),
		}
		docReader := esutil.NewJSONReader(&jsonElementStruct)

		res, err := log.es.Index(
			log.indexName, // Index name
			docReader,     // Document body
		)
		if err != nil {
			log.msgCounter++
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
