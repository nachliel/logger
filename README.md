# logger

A Simple ElasticSearch logger for Go Application.

## Quice Start

The recommended way to use log is to create your own logger ;)
The simplest way to do this is using this logger.

Start new project and follow the steps below:
Get the package:
```
go get github.com/nachliel/logger@latest
```
logger must be initialized before use:
``` go
logger.SetupWriter(logger.LevelInfo)
```

## Usage ##

```go
package main

import (
	"github.com/nachliel/logger"
)

func main() {
  // Sets from what level to output messages
  logger.SetupWriter(logger.LevelInfo)
  logger.Info("Hello World")
  logger.Debug("This is Debugging error number: %d",15) // Will not show, due to choosen level Info.
}
```

## Why?
`logger` is designed to be a simple and universal go logging library with
support for indexing the logs to elasticsearch server. You are able to choose from which 
Level the logger will output. anyhow the loger will output to the console, and ES if 
choosen to do so.

### Logger Levels ###
1. Debug - logger.LevelDebug  - For debug purposes
2. Info - logger.LevelInfo	- for Information
3. Warning - logger.LevelWarning	- Warning log notes
4. Error - logger.LevelError	- Error on log notes
5. Fatal - logger.LevelFatal	- Fatal Error leads to EXIT 0
