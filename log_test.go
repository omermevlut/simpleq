package simpleq

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
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

		expect := []byte(`[1;31merror
[0m`)

		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Expected Error() to output %v, got %v", out, expect)
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

		expect := []byte(`[1;34minfo
[0m`)

		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Expected Info() to output %v, got %v", out, expect)
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

		expect := []byte(`[1;33mwarn
[0m`)

		if !reflect.DeepEqual(out, expect) {
			t.Errorf("Expected Warn() to output %v, got %v", out, expect)
		}
	})
}
