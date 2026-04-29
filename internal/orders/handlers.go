package orders

import (
	"errors"
	"log"
	"net/http"

	"github.com/AmmanSajid1/go-ecom-api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req PlaceOrderRequest
	if err := json.Read(r, &req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.PlaceOrder(r.Context(), req.CustomerID, req.Items)
	if err != nil {
		log.Println(err)

		if errors.Is(err, ErrInvalidOrder) ||
			errors.Is(err, ErrInvalidCustomer) ||
			errors.Is(err, ErrInvalidItem) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if errors.Is(err, ErrProductNoStock) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, order)
}

func (h *handler) GetOrderByID(w http.ResponseWriter, r *http.Request, orderId int) {
	order, err := h.service.GetOrderByID(
		r.Context(),
		orderId,
	)

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, order)
}
