package api

func (s *Server) MountHandlers() {

	// Mount all handlers here
	s.Router.Get("/", Test)

}
