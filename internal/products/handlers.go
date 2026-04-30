package products

import (
	"log"
	"net/http"

	repo "github.com/AmmanSajid1/go-ecom-api/internal/adapters/postgresql/sqlc"
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

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if products == nil {
		products = []repo.Product{}
	}

	json.Write(w, http.StatusOK, products)

}

func (h *handler) FindProductByID(w http.ResponseWriter, r *http.Request, productId int) {
	product, err := h.service.FindProductByID(r.Context(), int(productId))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}
