package logger

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestCreating(t *testing.T) {
	SetupWriter(LevelWarn)
	Debug("This is Log Debug Message")
	Info("This is a Log Information Message")
	Warn("This is a Log Warning Message")
	Error("This is a Log Error Message")
	//Fatal("This is a Log Info Message")
}

func TestLoopMessages(t *testing.T) {
	Info("Start Sequence Messeging")
	for i := 0; i < 10; i++ {
		Debug("Message Number: %d.", i)
	}
	//Fatal("This is a Log Info Message")
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
/*
func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")
	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}
*/
