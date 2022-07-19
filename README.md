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
  logme.Info("Hello World")
  logme.Debug("This is Debugging error number: %d",15)
 
}
```
