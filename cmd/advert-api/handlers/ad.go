package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/Reywaltz/avito_advertising/pkg/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AdRepository interface {
	GetAll() ([]models.Ad, error)
	Create(models.Ad) (uuid.UUID, error)
}

type AdHandlers struct {
	log    log.Logger
	adRepo AdRepository
}

func NewHandlers(log log.Logger, adRepo AdRepository) (*AdHandlers, error) {
	return &AdHandlers{
		log:    log,
		adRepo: adRepo,
	}, nil
}

func (h *AdHandlers) GetAds(w http.ResponseWriter, r *http.Request) {
	res, err := h.adRepo.GetAll()
	if err != nil {
		h.log.Errorf("Can't get data from DB: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	out, err := json.Marshal(res)
	if err != nil {
		h.log.Errorf("Can't marshall result: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(out)
}

func (h *AdHandlers) CreateAd(w http.ResponseWriter, r *http.Request) {
	var newAd models.Ad
	if err := newAd.Bind(r); err != nil {
		h.log.Errorf("Can't bind ad: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newAd.ID = uuid.New().String()
	out, err := h.adRepo.Create(newAd)
	if err != nil {
		h.log.Errorf("Can't insert new ad: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	RawJSON := []byte(`{"id": "` + out.String() + `"}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(RawJSON)

}

func (h *AdHandlers) Routes(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/ads", h.GetAds).Methods("GET")
	subrouter.HandleFunc("/ads", h.CreateAd).Methods("POST")

}
