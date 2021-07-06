package simpleq

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Task is a task interface to be implemented a users
// Once implemented it should be passed into executor as queue.OnExec(new(impl))
type Task interface {
	Run(c Context) error
	Fail(err error)
}

// Context is a task interface
type Context interface {
	SetID()
	SetContent(c Content)
	Marshal() ([]byte, error)
	Bind(v interface{}) error
	GetContent() Content
	GetID() string
	SetMaxAttempts(a int)
	GetAttempts() int
	GetMaxAttempts() int
	NewAttempt()
}

// Content is a task content helper construct
type Content []byte

// BindJSON binds JSON data to a given construct
func (c Content) BindJSON(v interface{}) error {
	return json.Unmarshal(c, &v)
}

// NewMessage returns a pointer to a new message instance
func NewMessage(c Content) *Message {
	return &Message{Content: c}
}

// Message is a single message instance
type Message struct {
	Attempts    int     `json:"attempts"`
	MaxAttempts int     `json:"max_attempts"`
	ID          string  `json:"id"`
	Content     Content `json:"content"`
}

// GetContent returns message content
func (m *Message) GetContent() Content {
	return m.Content
}

// GetID returns message ID
func (m *Message) GetID() string {
	return m.ID
}

// GetAttempts returns number of attempts
func (m *Message) GetAttempts() int {
	return m.Attempts
}

// GetMaxAttempts returns the max attempt count for given task
func (m *Message) GetMaxAttempts() int {
	return m.MaxAttempts
}

// SetID sets an arbitrary unique ID to a given task
func (m *Message) SetID() {
	m.ID = uuid.New().String()
}

// SetContent sets data to task
func (m *Message) SetContent(c Content) {
	m.Content = c
}

// SetMaxAttempts sets the number of max attempts for current task
func (m *Message) SetMaxAttempts(a int) {
	m.MaxAttempts = a
}

// NewAttempt increments the attempt number
func (m *Message) NewAttempt() {
	m.Attempts++
}

// Marshal json marshals base task
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Bind task data to a given construct
func (m *Message) Bind(v interface{}) error {
	return m.Content.BindJSON(v)
}
