package simpleq

// Driver is queue driver interface, can be implemented externally
type Driver interface {
	Write(queue string, d []byte) error
	Read(queue string) ([]byte, error)
	SetProcessed(queue string) error
	Register(queue string) error
	SetFailed(queue string, taskID string) error
	GetStats() (*Stats, error)
}
