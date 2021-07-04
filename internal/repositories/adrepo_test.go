package repositories

import (
	"testing"

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
	rows := pgxmock.NewRows([]string{"id", "name", "desciption", "photos", "cost"})

	mock.ExpectQuery(GetAds).WillReturnRows(rows)

	_, err = repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Not met: %s", err)
	}

}
