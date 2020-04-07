package requestfiltration

import "encoding/json"

type Route struct {
	Location string
	Pass     RoutePass
}

type RoutePass struct {
	Host string
	Port string
	Path string
}

func ParseRoutes(b []byte) ([]Route, error) {
	var (
		routes []Route
		err    error
	)

	err = json.Unmarshal(b, &routes)
	if err != nil {
		return routes, err
	}

	return routes, nil
}
