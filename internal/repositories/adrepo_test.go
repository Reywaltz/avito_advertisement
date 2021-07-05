package repositories

import (
	"regexp"
	"testing"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func initMockItemHandler() (pgxmock.PgxConnIface, *AdRepo, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}

	repo := NewRepo(mock)
	return mock, repo, nil
}

var (
	testCost = decimal.RequireFromString("231.1")
)

func TestAdRepositoryGetAll(t *testing.T) {
	mock, repo, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}
	rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"}).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, testCost).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, testCost)

	mock.ExpectQuery(GetAds).WillReturnRows(rows)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}

}

func TestCreate(t *testing.T) {
	mock, adRepo, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}

	tmpUUID := uuid.New()
	tmp := models.Ad{
		ID:          tmpUUID,
		Name:        "New Ad",
		Description: "Описание",
		Photos:      []string{"1.jpg"},
		Cost:        decimal.Zero,
	}

	rows := pgxmock.NewRows([]string{"id"}).AddRow(tmp.ID.String())

	mock.ExpectQuery(regexp.QuoteMeta(CreateAdQuery)).WithArgs(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Photos, &tmp.Cost).
		WillReturnRows(rows)

	q, err := adRepo.Create(tmp)
	if err != nil {
		assert.NotNil(t, q)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}
}
