package types

type Metadata struct {
	Version     string `json:"version"`
	Route       string `json:"route"`
	Name        string `json:"name"`
	ServiceGuid string `json:"serviceguid"`
	RouteGuid   string `json:"routeguid"`
}
