package handlers_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/Reywaltz/avito_advertising/cmd/advert-api/handlers"
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

func initMockItemHandler() (pgxmock.PgxConnIface, *handlers.AdHandlers, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}

	log, err := log.NewLogger("DEV")
	if err != nil {
		return nil, nil, err
	}

	repo := repositories.NewRepo(mock)
	adHandlers, err := handlers.NewHandlers(log, repo)
	if err != nil {
		return nil, nil, err
	}

	return mock, adHandlers, nil
}

func TestAdGetAllHandler(t *testing.T) {
	t.Parallel()
	mock, adHandler, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}

	type TestCase struct {
		Name         string
		In           []byte
		ExpectedCode int
	}

	testcases := []TestCase{
		{Name: "OK get", In: []byte{}, ExpectedCode: http.StatusOK},
		{Name: "Internal error", In: []byte{}, ExpectedCode: http.StatusInternalServerError},
	}

	handler := http.HandlerFunc(adHandler.GetAds)

	for _, tc := range testcases {
		tc := tc
		switch tc.Name {
		case "OK get":
			t.Run(tc.Name, func(t *testing.T) {
				rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"}).
					AddRow(uuid.New(), "Индекс", "Очень интересное объявление№1", []string{"1.jpg, 2.jpg, 3.jpg"}, decimal.Zero).
					AddRow(uuid.New(), "Индекс", "Очень интересное объявление№2", []string{"1.jpg, 2.jpg, 3.jpg"}, decimal.Zero)

				mock.ExpectQuery(GetAds).WillReturnRows(rows)

				request := httptest.NewRequest(http.MethodGet, "/api/v1/ads", nil)
				recorder := httptest.NewRecorder()
				handler.ServeHTTP(recorder, request)

				assert.Equal(t, tc.ExpectedCode, recorder.Code, "Got: %s. Expected %s", recorder.Code, http.StatusOK)
			})

		case "Internal error":
			t.Run(tc.Name, func(t *testing.T) {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT *`)).
					WillReturnError(errors.New("Can't get db conn"))

				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodGet, "/api/v1/ads", bytes.NewBuffer(tc.In))
				handler.ServeHTTP(recorder, request)

				if err := mock.ExpectationsWereMet(); err != nil {
					assert.Equal(t, tc.ExpectedCode, recorder.Code, "Got: %s. Expected %s. Err: %v", recorder.Code, http.StatusOK, err.Error())
				}
			})
		}
	}
}

func TestAdCreateHandler(t *testing.T) {
	t.Parallel()
	mock, adHandler, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}

	type TestCase struct {
		Name         string
		In           []byte
		ExpectedCode int
	}

	OKJSON := []byte(`{
		"name": "New Ad",
		"description": "Описание",
		"photos": ["1.jpg"],
		"cost": "33.3"
	}`)

	BadJSON := []byte(`{
		"description": "Описание",
		"photos": ["1.jpg"],
		"cost": "33.3"
	}`)

	handler := http.HandlerFunc(adHandler.CreateAd)
	testcases := []TestCase{
		{Name: "Normal input", In: OKJSON, ExpectedCode: http.StatusOK},
		{Name: "Bad input", In: BadJSON, ExpectedCode: http.StatusBadRequest},
		{Name: "Internal error", In: OKJSON, ExpectedCode: http.StatusInternalServerError},
	}

	for _, tc := range testcases {
		tc := tc
		switch tc.Name {
		case "Normal input":
			t.Run(tc.Name, func(t *testing.T) {
				tmpUUID := uuid.New()

				rows := pgxmock.NewRows([]string{"id"}).AddRow(tmpUUID.String())

				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO`)).
					WillReturnRows(rows)

				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewBuffer(tc.In))
				handler.ServeHTTP(recorder, request)

				if err := mock.ExpectationsWereMet(); err != nil {
					assert.Equal(t, http.StatusOK, recorder.Code, "Got: %s. Expected %s. Err: %v", recorder.Code, http.StatusOK, err.Error())
				}
			})

		case "Bad input":
			t.Run(tc.Name, func(t *testing.T) {
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewBuffer(tc.In))
				handler.ServeHTTP(recorder, request)

				assert.Equal(t, tc.ExpectedCode, recorder.Code, "Got: %s. Actual: %s")
			})

		case "Internal error":
			t.Run(tc.Name, func(t *testing.T) {
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO`)).
					WillReturnError(errors.New("Duplicate"))

				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/api/v1/ads", bytes.NewBuffer(tc.In))
				handler.ServeHTTP(recorder, request)

				if err := mock.ExpectationsWereMet(); err != nil {
					assert.Equal(t, http.StatusOK, recorder.Code, "Got: %s. Expected %s. Err: %v", recorder.Code, http.StatusOK, err.Error())
				}
			})
		}
	}
}
