package route

type RouteRpc struct {
	Tables []RouteTable `xml:"route-table"`
}

type RouteTable struct {
	Name         string               `xml:"table-name"`
	MaxRoutes    int64                `xml:"prefix-max"`
	TotalRoutes  int64                `xml:"total-route-count"`
	ActiveRoutes int64                `xml:"active-route-count"`
	Protocols    []RouteTableProtocol `xml:"protocols"`
}

type RouteTableProtocol struct {
	Name         string `xml:"protocol-name"`
	Routes       int64  `xml:"protocol-route-count"`
	ActiveRoutes int64  `xml:"active-route-count"`
}
