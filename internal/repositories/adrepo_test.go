package repositories_test

import (
	"regexp"
	"testing"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/Reywaltz/avito_advertising/internal/repositories"
	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
	"github.com/shopspring/decimal"
)

func initMockAdRepo() (pgxmock.PgxConnIface, *repositories.AdRepo, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}

	repo := repositories.NewRepo(mock)

	return mock, repo, nil
}

func TestAdRepositoryGetAll(t *testing.T) {
	t.Parallel()
	mock, repo, err := initMockAdRepo()
	if err != nil {
		t.Fatal(err)
	}
	rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"}).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, decimal.Zero).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, decimal.Zero)

	mock.ExpectQuery(repositories.GetAds).WillReturnRows(rows)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()
	mock, adRepo, err := initMockAdRepo()
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

	mock.ExpectQuery(regexp.QuoteMeta(repositories.CreateAdQuery)).WithArgs(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Photos, &tmp.Cost).
		WillReturnRows(rows)

	_, err = adRepo.Create(tmp)
	if err != nil {
		t.Errorf("Can't insert obj: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}
}
