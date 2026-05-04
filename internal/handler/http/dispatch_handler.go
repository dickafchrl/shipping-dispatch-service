package http

import (
	"net/http"

	"github.com/archera/shipping-service-ver2/internal/domain"
)

type DispatchHandler struct {
	service domain.DispatchService
}

func NewDispatchHandler(svc domain.DispatchService) *DispatchHandler {
	return &DispatchHandler{service: svc}
}

// AutoDispatch adalah endpoint untuk HTTP POST /dispatch
func (h *DispatchHandler) AutoDispatch(w http.ResponseWriter, r *http.Request) {
	// Red Phase: Belum ada logika parsing JSON atau pemanggilan service.
	// Kita sengaja langsung mengembalikan error 501.
	http.Error(w, "not implemented", http.StatusNotImplemented)
}