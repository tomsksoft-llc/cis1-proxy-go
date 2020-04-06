package requestfiltration

import "encoding/json"

type Route struct {
	Location string
	Pass     string
}

func ParseRoutes(data []byte) ([]Route, error) {
	var (
		routes []Route
		err    error
	)

	err = json.Unmarshal(data, &routes)
	if err != nil {
		return make([]Route, 0), err
	}

	return routes, nil
}
