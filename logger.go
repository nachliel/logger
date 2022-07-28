package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type any = interface{}

// Writer
var logWriter *bufio.Writer

// logStruct struct is the main structure of the logger settings
type logStruct struct {
	level        logLevel
	msgCounter   int64
	es           *elasticsearch.Client
	proccessName string
	indexName    string
	timeFormat   string
	fatalTime    int
}

// logLevel of error
type logLevel int

// logDoc is the Json structure of log
type logDoc struct {
	ProccessName string    `json:"proccess"`
	StatusTime   time.Time `json:"status_time"`
	Level        string    `json:"level"`
	Message      string    `json:"message"`
}

// Const Name of Log Levels
const (
	LevelDebug logLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Settings of the logger
var settings logStruct

// NewLogger create a new Logger
func SetLevel(level logLevel) {
	settings.level = level
}

// SetFatalTimer sets time to pause after Fatal Error.
// Used to pause the proccess from restart by any deamon.
func SetFatalTimer(ms int) {
	settings.fatalTime = ms
}

// SetupWriter Initiate the logger, must be declared on begging
func SetupWriter(lvl logLevel) {
	logWriter = bufio.NewWriter(os.Stdout)
	settings.timeFormat = time.RFC822
	settings.level = lvl
	settings.proccessName = "logger"
}

// SetTimeFormat settings the logger time format
func SetTimeFormat(format string) {
	settings.timeFormat = format
}

//	Add ElasticSearch Client to log messages there
func SetElasticClient(esClient *elasticsearch.Client, index string) {
	settings.es = esClient
	settings.indexName = index
}

func SetProccessName(name string) {
	settings.proccessName = name
}

//	Debug Logging
func Debug(format string, args ...any) {
	if LevelDebug < settings.level {
		return
	}
	write(LevelDebug, format, args...)
}

// Info Logging
func Info(format string, args ...any) {
	if LevelInfo < settings.level {
		return
	}
	write(LevelInfo, format, args...)
}

// Warn Logging
func Warn(format string, args ...any) {
	if LevelWarn < settings.level {
		return
	}
	write(LevelWarn, format, args...)
}

// Error Logging
func Error(format string, args ...any) {
	if LevelError < settings.level {
		return
	}
	write(LevelError, format, args...)
}

// Fatal Logging
func Fatal(format string, args ...any) {
	write(LevelFatal, format, args...)
	// Exit Failure
	time.Sleep(time.Duration(settings.fatalTime) * time.Millisecond)
	os.Exit(1)
}

// write Log to console, if ES is setup then also to ES.
func write(level logLevel, format string, args ...any) {
	// Log Message counter
	settings.msgCounter++

	logWriter.WriteString(time.Now().Format(settings.timeFormat))
	logWriter.WriteString(" [")
	logWriter.WriteString(level.toString())
	logWriter.WriteString("]\t")
	logWriter.WriteString(fmt.Sprintf(format, args...))
	logWriter.WriteString("\n")
	// Output to Terminal Buffered Write
	logWriter.Flush()

	// Send to ElasticSearch if any
	if settings.es != nil {
		// Write Log to ES
		writeESDoc(logDoc{
			ProccessName: settings.proccessName,
			StatusTime:   time.Now().UTC(),
			Level:        level.toString(),
			Message:      fmt.Sprintf(format, args...),
		})
	}
}

func writeESDoc(doc logDoc) {
	docReader := esutil.NewJSONReader(&doc)
	res, err := settings.es.Index(
		settings.indexName, // Index name
		docReader,          // Document body
	)
	if err != nil {
		settings.msgCounter++

		logWriter.WriteString(time.Now().Format(settings.timeFormat))
		logWriter.WriteString(" [")
		logWriter.WriteString(LevelFatal.toString())
		logWriter.WriteString("]\tElasticSearch Fatal Error: ")
		logWriter.WriteString(err.Error())
		logWriter.WriteString("\n")
		// Output to Terminal Buffered Write
		logWriter.Flush()
		// Exit!
		time.Sleep(time.Duration(settings.fatalTime) * time.Millisecond)
		os.Exit(1)
	}
	res.Body.Close()
}

// Converts Level to string.
func (level logLevel) toString() string {
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
