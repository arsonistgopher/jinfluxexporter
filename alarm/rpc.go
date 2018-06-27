package alarm

// AlarmRPC struct
type AlarmRPC struct {
	Details []AlarmDetails `xml:"alarm-detail"`
	Summary []AlarmSummary `xml:"alarm-summary"`
}

// AlarmDetails struct
type AlarmDetails struct {
	Class       string `xml:"alarm-class"`
	Description string `xml:"alarm-description"`
	Type        string `xml:"alarm-type"`
}

// AlarmSummary struct
type AlarmSummary struct {
	Count int `xml:"active-alarm-count"`
}
