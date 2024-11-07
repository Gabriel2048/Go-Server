package routing

import (
	"errors"
	"slices"
	"strings"
)

type HttpRouteTemplate struct {
	nodes []routeNodeTemplate
}

type routeNodeTemplate struct {
	value       string
	isParameter bool
}

func NewRouteFromString(rawRoute string) (HttpRouteTemplate, error) {
	if rawRoute != "/" {
		rawRoute = strings.Trim(rawRoute, "/")
	}

	parts := strings.Split(rawRoute, "/")

	nodes := make([]routeNodeTemplate, 0, len(parts))

	for index, part := range parts {

		if part == "/" && parts[index+1] == "/" && index < len(parts)-1 { // two '/' next to each other make for an invalid url
			return HttpRouteTemplate{}, errors.New("invalid route due to 2 '/'")
		}
		value := part

		isParameter := len(part) > 2 && strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}")
		if isParameter {
			value = part[1 : len(part)-1] //trims the {} surrounding the value
			if strings.Contains(value, "}") || strings.Contains(value, "{") {
				return HttpRouteTemplate{}, errors.New("malformed parameter")
			}
		}

		node := routeNodeTemplate{
			value:       value,
			isParameter: isParameter,
		}

		if isParameter && slices.Contains(nodes, node) {
			return HttpRouteTemplate{}, errors.New("duplicated path parameter")
		}

		nodes = append(nodes, node)
	}

	return HttpRouteTemplate{
		nodes: nodes,
	}, nil
}

func (h HttpRouteTemplate) CanHandlerPath(path string) (RouteParameters, bool) {

	if path != "/" {
		path = strings.Trim(path, "/")
	}

	parts := strings.Split(path, "/")

	if len(parts) != len(h.nodes) {
		return nil, false
	}

	routeParams := RouteParameters{}

	for index, nodeTemplate := range h.nodes {
		if !nodeTemplate.isParameter {
			if nodeTemplate.value != parts[index] {
				return nil, false
			}
		} else {
			routeParams[nodeTemplate.value] = parts[index]
		}
	}

	return routeParams, true
}

type RouteParameters map[string]string
