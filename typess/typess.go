package typess

// Object is anything with a key string and a value string
type ObjectCh struct {
	Key   string            `json:"key"`
	Value map[string]string `json:"value"`
}
