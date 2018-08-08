package environment

// EnvironmentRpc type with XML tags for easy marhsaling/unmarshaling
type EnvironmentRpc struct {
	Items []EnvironmentItemRpc `xml:"environment-item"`
}

// EnvironmentItemRpc type with XML tags for easy marhsaling/unmarshaling
type EnvironmentItemRpc struct {
	Name        string `xml:"name"`
	Class       string `xml:"class"`
	Status      string `xml:"status"`
	Temperature *struct {
		Value float64 `xml:"celsius,attr"`
	} `xml:"temperature,omitempty"`
}
