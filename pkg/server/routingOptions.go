package server

import "fmt"

const forexRouteGroupPrefix = "/forex"
const projectsRouteGroupPrefix = "/projects"

func (s *Server) withDefaultRoutes() {
	projects := s.Group(projectsRouteGroupPrefix)
	projects.POST("", s.CreateProject)

	forex := s.Group(forexRouteGroupPrefix)
	forex.GET(fmt.Sprintf("/rates/:%s", fromCurParam), s.ExchangeRatesForCurrency)
}
