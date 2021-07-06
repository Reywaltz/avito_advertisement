package postgres_test

import (
	"testing"

	"github.com/Reywaltz/avito_advertising/pkg/postgres"
	"github.com/stretchr/testify/assert"
)

const (
	CONN_DB = `postgres://advert_user:pass@localhost:5433/advert`
)

func TestDBConnection(t *testing.T) {
	t.Parallel()

	type testCase struct {
		Name        string
		In          string
		ExpectedErr bool
	}

	testcases := []testCase{
		{Name: "Connection string is not set", In: "", ExpectedErr: true},
		{Name: "Can't connect to database", In: "SomeConn", ExpectedErr: true},
		{Name: "Connection established", In: CONN_DB, ExpectedErr: false},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			_, err := postgres.NewDB(tc.In)
			assert.Equal(t, tc.ExpectedErr, err != nil)
		})
	}
}
