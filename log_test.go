package simpleq

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestDefaultLogger_Error(t *testing.T) {
	t.Run("it_should_output_formatted_error", func(t *testing.T) {
		lg := DefaultLogger{}

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		lg.Error(fmt.Errorf("error"))

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout

		if !strings.Contains(string(out), "[ERROR] error") {
			t.Errorf("Expected Error() to contain output [ERROR] error, got %v", out)
		}
	})
}

func TestDefaultLogger_Info(t *testing.T) {
	t.Run("it_should_output_formatted_info", func(t *testing.T) {
		lg := DefaultLogger{}

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		lg.Info("info")

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout

		if !strings.Contains(string(out), "[INFO] info") {
			t.Errorf("Expected Error() to contain output [INFO] info, got %v", out)
		}
	})
}

func TestDefaultLogger_Warn(t *testing.T) {
	t.Run("it_should_output_formatted_warning", func(t *testing.T) {
		lg := DefaultLogger{}

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		lg.Warn("warn")

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout

		if !strings.Contains(string(out), "[WARNING] warn") {
			t.Errorf("Expected Error() to contain output [WARNING] warn, got %v", out)
		}
	})
}
