# logger
Simple logger for Go Application.

## Usage ##
Declare:
```go

package main

import (
	"github.com/nachliel/logger"
)

var logme logger.Logger

func main() {
  // Sets from what level to output messages
  logme.SetLevel(logger.LevelInfo)
  logme.Info("Hello World")
  logme.Debug("This is Debugging error number: %d",15)
  // Assign ElasticSearch logs ..
}
```

### Logger Levels ###
1. Debug - logger.LevelDebug  - For debug purposes
2. Info - logger.LevelInfo	- for Information
3. Warning - logger.LevelWarning	- Warning log notes
4. Error - logger.LevelError	- Error on log notes
5. Fatal - logger.LevelFatal	- Fatal Error leads to EXIT 0
