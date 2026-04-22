package shortly

import "net/http"

type HealthService struct{}

func (hs *HealthService) HandlerHealth(w http.ResponseWriter, r *http.Request) {

	_, _ = w.Write([]byte("Server is up and running"))

}

func newHealthHandler() *HealthService {
	return &HealthService{}
}
