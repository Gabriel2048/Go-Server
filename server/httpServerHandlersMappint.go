package server

import (
	"errors"
	"http-server/server/routing"
)

func (s *Server) MapGet(path string, handler HttpHandler) error {

	route, err := routing.NewRouteFromString(path)

	if err != nil {
		return errors.Join(errors.New("unable to register the route"), err)
	}

	s.routes = append(s.routes, routeWithHandler{
		verb:    GET,
		route:   route,
		handler: handler,
	})

	return nil
}

func (s *Server) MapPost(path string, handler HttpHandler) error {
	route, err := routing.NewRouteFromString(path)

	if err != nil {
		return errors.Join(errors.New("unable to register the route"), err)
	}

	s.routes = append(s.routes, routeWithHandler{
		verb:    POST,
		route:   route,
		handler: handler,
	})

	return nil
}
