package simpleq

import (
	"encoding/json"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestContent_BindJSON(t *testing.T) {
	t.Run("it_should_unmarshal_json", func(t *testing.T) {
		c := Content(`{"a": "b"}`)

		expect := map[string]interface{}{"a": "b"}

		var got map[string]interface{}

		_ = c.BindJSON(&got)

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expected BindJSON() to bind %v, got %v", expect, got)
		}
	})
}

func TestNewMessage(t *testing.T) {
	t.Run("it_should_return_new_message", func(t *testing.T) {
		expect := Message{Content: Content("{}")}

		if got := NewMessage(Content("{}")); !reflect.DeepEqual(got, &expect) {
			t.Errorf("Expected NewMessage() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_GetContent(t *testing.T) {
	t.Run("it_should_return_task_content", func(t *testing.T) {
		m := Message{Content: Content("{}")}

		expect := Content("{}")

		if got := m.GetContent(); !reflect.DeepEqual(expect, got) {
			t.Errorf("Expected GetContent() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_GetID(t *testing.T) {
	t.Run("it_should_return_task_id", func(t *testing.T) {
		id := uuid.New().String()

		msg := Message{ID: id}

		expect := id

		if got := msg.GetID(); expect != got {
			t.Errorf("Expected GetID() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_GetAttempts(t *testing.T) {
	t.Run("it_should_return_number_of_attempts", func(t *testing.T) {
		m := Message{Attempts: 4}
		expect := 4

		if got := m.GetAttempts(); expect != got {
			t.Errorf("Expected GetAttempts() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_GetMaxAttempts(t *testing.T) {
	t.Run("it_should_return_number_of_attempts", func(t *testing.T) {
		m := Message{MaxAttempts: 10}
		expect := 10

		if got := m.GetMaxAttempts(); expect != got {
			t.Errorf("Expected GetMaxAttempts() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_NewAttempt(t *testing.T) {
	t.Run("it_should_increment_number_of_attempts", func(t *testing.T) {
		m := Message{Attempts: 5}
		m.NewAttempt()

		expect := 6

		if got := m.Attempts; got != expect {
			t.Errorf("Expected NewAttempt() to set %v, got %v", expect, got)
		}
	})
}

func TestMessage_Marshal(t *testing.T) {
	t.Run("it_should_marshal_task", func(t *testing.T) {
		m := Message{Attempts: 1}

		got, _ := m.Marshal()
		expect, _ := json.Marshal(m)

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expected Marshal() to return %v, got %v", expect, got)
		}
	})
}

func TestMessage_Bind(t *testing.T) {
	t.Run("it_should_bind_data", func(t *testing.T) {
		m := Message{Content: Content(`{"a": "b"}`)}

		expect := map[string]interface{}{"a": "b"}

		var got map[string]interface{}

		_ = m.Bind(&got)

		if !reflect.DeepEqual(got, expect) {
			t.Errorf("Expected Bind() to bind %v, got %v", expect, got)
		}
	})
}

func TestMessage_SetID(t *testing.T) {
	t.Run("it_should_set_task_id", func(t *testing.T) {
		m := Message{}
		m.SetID()

		if m.ID == "" {
			t.Errorf("Expected GetID() not to be empty")
		}
	})
}

func TestMessage_SetContent(t *testing.T) {
	t.Run("it_should_set_content", func(t *testing.T) {
		m := Message{}
		m.SetContent(Content("test-content"))

		expect := Content("test-content")

		if !reflect.DeepEqual(m.Content, expect) {
			t.Errorf("Expected SetContent() to set %v, got %v", expect, m.Content)
		}
	})
}

func TestMessage_SetMaxAttempts(t *testing.T) {
	t.Run("it_should_set_message_max_attempts", func(t *testing.T) {
		m := Message{}
		m.SetMaxAttempts(5)

		expect := 5

		if m.MaxAttempts != expect {
			t.Errorf("Expected SetMaxAttempts() to set %v, got %v", expect, m.MaxAttempts)
		}
	})
}
