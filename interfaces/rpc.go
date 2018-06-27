package interfaces

type InterfaceRpc struct {
	Interfaces []PhyInterface `xml:"physical-interface"`
}

type PhyInterface struct {
	Name              string         `xml:"name"`
	AdminStatus       string         `xml:"admin-status"`
	OperStatus        string         `xml:"oper-status"`
	Description       string         `xml:"description"`
	MacAddress        string         `xml:"current-physical-address"`
	Stats             TrafficStat    `xml:"traffic-statistics"`
	LogicalInterfaces []LogInterface `xml:"logical-interface"`
	InputErrors       struct {
		Drops  int64 `xml:"input-drops"`
		Errors int64 `xml:"input-errors"`
	} `xml:"input-error-list"`
	OutputErrors struct {
		Drops  int64 `xml:"output-drops"`
		Errors int64 `xml:"output-errors"`
	} `xml:"output-error-list"`
}

type LogInterface struct {
	Name        string      `xml:"name"`
	Description string      `xml:"description"`
	Stats       TrafficStat `xml:"traffic-statistics"`
}

type TrafficStat struct {
	InputBytes  int64 `xml:"input-bytes"`
	OutputBytes int64 `xml:"output-bytes"`
}
