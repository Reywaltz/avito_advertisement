package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Reywaltz/avito_advertising/internal/repositories"
	"github.com/Reywaltz/avito_advertising/pkg/log"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	AdsFields = `name, description, photos, cost`
	AllFileds = `id, ` + AdsFields
	GetAds    = `SELECT ` + AllFileds + ` FROM advertisement`
)

var (
	testCost = decimal.RequireFromString("231.1")
)

func initMockItemHandler() (pgxmock.PgxConnIface, *AdHandlers, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}

	log, err := log.NewLogger("DEV")
	if err != nil {
		return nil, nil, err
	}

	repo := repositories.NewRepo(mock)
	adHandlers, err := NewHandlers(log, repo)
	if err != nil {
		return nil, nil, err
	}
	return mock, adHandlers, nil
}

func TestAdGetAllHandler(t *testing.T) {
	mock, adHandler, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}
	rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"}).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, testCost).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, testCost)

	mock.ExpectQuery(GetAds).WillReturnRows(rows)

	handler := http.HandlerFunc(adHandler.GetAds)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/ads", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	assert.Equal(t, http.StatusOK, recorder.Code, "Got: %s. Expected %s", recorder.Code, http.StatusOK)

}

const (
	CreateAdQuery = `INSERT INTO advertisement (` + AllFileds + `) VALUES ($1, $2, $3, $4, $5) returning id`
)

func TestAdCreateAllHandler(t *testing.T) {
	mock, adHandler, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}

	tmpUUID := uuid.New()

	rows := pgxmock.NewRows([]string{"id"}).
		AddRow(tmpUUID)

	mock.ExpectQuery("INSERT INTO advertisement").WithArgs(&tmpUUID, "New Ad", "Описание", []string{"1.jpg"}, decimal.Zero).
		WillReturnRows(rows)

	incomingJSON := []byte(`{
		"name": "New Ad",
		"description": "Описание",
		"photos": ["1.jpg"],
		"cost": "33.3"
	}`)

	handler := http.HandlerFunc(adHandler.CreateAd)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewBuffer(incomingJSON))
	handler.ServeHTTP(recorder, request)

	if err := mock.ExpectationsWereMet(); err != nil {
		assert.Equal(t, http.StatusOK, recorder.Code, "Got: %s. Expected %s. Err: %s", recorder.Code, http.StatusOK, err)
	}

}
