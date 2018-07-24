package environment

type EnvironmentRpc struct {
	Items []EnvironmentItemRpc `xml:"environment-item"`
}

type EnvironmentItemRpc struct {
	Name        string `xml:"name"`
	Class       string `xml:"class"`
	Status      string `xml:"status"`
	Temperature *struct {
		Value float64 `xml:"celsius,attr"`
	} `xml:"temperature,omitempty"`
}
