package server

import "fmt"

const forexRouteGroupPrefix = "/forex"
const projectsRouteGroupPrefix = "/projects"

func (s *Server) withDefaultRoutes() {
	projects := s.Group(projectsRouteGroupPrefix)
	projects.POST("", s.CreateProject)

	// Guard forex routes with API key verification
	forex := s.Group(forexRouteGroupPrefix, verifyAPIKeyMiddleware(s))
	forex.GET(fmt.Sprintf("/rates/:%s", fromCurParam), s.ExchangeRatesForCurrency)
}
