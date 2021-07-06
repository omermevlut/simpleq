package simpleq

// Stats is a list of stats of registered queues
type Stats map[string]Stat

// Stat is a single stat item
type Stat struct {
	Failed    int
	Processed int64
	FailedIDs []string
}
