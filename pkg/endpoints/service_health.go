package endpoints

import "net/http"

func (s *Service) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}
