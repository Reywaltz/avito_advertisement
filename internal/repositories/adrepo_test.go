package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock"
)

func initMockItemHandler() (pgxmock.PgxConnIface, *AdRepo, error) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		return nil, nil, err
	}

	repo := NewRepo(mock)
	return mock, repo, nil
}

func TestAdRepository(t *testing.T) {
	mock, repo, err := initMockItemHandler()
	if err != nil {
		t.Fatal(err)
	}
	rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"}).
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, "500").
		AddRow(uuid.New(), "Индекс", "Очень интересное объявление", []string{"1.jpg, 2.jpg, 3.jpg"}, "500")

	mock.ExpectQuery(GetAds).WillReturnRows(rows)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}

}
