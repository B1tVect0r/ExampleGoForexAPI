package server

import "fmt"

const forexRouteGroupPrefix = "/forex"
const projectsRouteGroupPrefix = "/projects"

func (s *Server) withDefaultRoutes() {
	projects := s.Group(projectsRouteGroupPrefix)
	projects.POST("", s.CreateProject)

	forex := s.Group(forexRouteGroupPrefix, VerifyAPIKeyMiddleware(s))
	forex.GET(fmt.Sprintf("/rates/:%s", fromCurParam), s.ExchangeRatesForCurrency)
}
