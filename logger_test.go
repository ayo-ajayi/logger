package logger

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
)

// Helper function to capture stdout output
func captureOutput(f func()) string {
	var buf bytes.Buffer
	stdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = stdout
	buf.ReadFrom(r)
	return buf.String()
}

func TestLogger_PrintEvent(t *testing.T) {
	l := NewLogger(DEBUG, false)

	output := captureOutput(func() {
		l.Info("This is an info message")
	})

	if !strings.Contains(output, "INFO") || !strings.Contains(output, "This is an info message") {
		t.Errorf("Expected log message not found in output: %s", output)
	}
}

func TestLogger_Callback(t *testing.T) {
	l := NewLogger(DEBUG, false)

	var callbackOutput string
	l.AddCallback(func(e EventInfo) {
		callbackOutput = fmt.Sprintf("%s: %s", levelStrings[e.Level()], e.Message())
	})

	l.Error("This is an error message")

	expected := "ERROR: This is an error message"
	if callbackOutput != expected {
		t.Errorf("Expected %s, got %s", expected, callbackOutput)
	}
}

func TestLogger_ThreadSafety(t *testing.T) {
	l := NewLogger(DEBUG, false)
	wg := sync.WaitGroup{}

	// Simulate concurrent log messages
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			l.Info(fmt.Sprintf("This is message %d", i))
		}(i)
	}
	wg.Wait()

	output := captureOutput(func() {
		l.Info("Final message")
	})

	if !strings.Contains(output, "Final message") {
		t.Errorf("Expected final log message not found in output: %s", output)
	}
}

func TestLogger_LevelFiltering(t *testing.T) {
	l := NewLogger(WARN, false)

	output := captureOutput(func() {
		l.Info("This info message should not appear")
		l.Warn("This warning message should appear")
	})

	if strings.Contains(output, "INFO") {
		t.Errorf("INFO log message should have been filtered out")
	}
	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected WARN log message not found in output: %s", output)
	}
}

func TestLogger_Fatal(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		l := NewLogger(DEBUG, false)
		l.Fatal("This is a fatal message")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogger_Fatal")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 1", err)
}
