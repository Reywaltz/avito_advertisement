package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/Reywaltz/avito_advertising/pkg/log"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

type AdRepository interface {
	GetAll() ([]models.Ad, error)
	Create(models.Ad) (uuid.UUID, error)
	GetOne(reqUUID uuid.UUID) (models.Ad, error)
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (h *AdHandlers) GetAd(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	adID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	reqUUID, err := uuid.Parse(adID)
	if err != nil {
		h.log.Errorf("Can't parse string:%s to UUID: %v", adID, err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := h.adRepo.GetOne(reqUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			h.log.Errorf("Ad with id:%s doesn't exists. %s", reqUUID, err)
			w.WriteHeader(http.StatusNotFound)

			return
		} else {
			h.log.Errorf("Can't get data from database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	out, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (h *AdHandlers) CreateAd(w http.ResponseWriter, r *http.Request) {
	var newAd models.Ad
	if err := newAd.Bind(r); err != nil {
		h.log.Errorf("Can't bind ad: %s", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	newAd.ID = uuid.New()
	out, err := h.adRepo.Create(newAd)
	if err != nil {
		h.log.Errorf("Can't insert new ad: %s", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	RawJSON := []byte(`{"id": "` + out.String() + `"}`)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(RawJSON)
}

func (h *AdHandlers) Routes(router *mux.Router) {
	subrouter := router.PathPrefix("/api/v1").Subrouter()
	subrouter.HandleFunc("/ads", h.GetAds).Methods(http.MethodGet)
	subrouter.HandleFunc("/ads", h.CreateAd).Methods(http.MethodPost)
	subrouter.HandleFunc("/ads/{id}", h.GetAd).Methods(http.MethodGet)
}
