package models

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdBind(t *testing.T) {
	type testCase struct {
		Name        string
		In          []byte
		ExpectedErr bool
	}

	const (
		URL = `/ads`
	)

	AdWithoutName := `{"description": "Ахуэнный товар", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img"], "cost": "12"}`

	AdWithLongName := `{"name": ` + strings.Repeat("LoooongName", 500) + `}`

	AdWithoutDescription := `{"name": "Name", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img"], "cost": "12"}`

	AdWithLongDescription := `{"name": "Item Name", "desciption": ` + strings.Repeat("LoooongName", 500) + `}`

	AdWithoutPhotos := `{"name": "Name", "description": "Ахуэнный товар", "cost": "12"}`

	AdWithManyPhotos := `{"name": "Item Name", "desciption": "Ахуэнный товар", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img", "tmp4/4.img"]}`

	AdWithoutCost := `{"name": "Name", "description": "Ахуэнный товар", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img"]}`

	AdWithNegativeCost := `{"name": "Item Name", "desciption": "Ахуэнный товар", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img"], "cost": "-12"}`

	AdJSON := `{"name": "New ad","description": "Ахуэнный товар", "photos": ["tmp1/1.img", "tmp2/2.img", "tmp3/3.img"], "cost": "12"}`

	InvalidJson := `"name": "New ad","description": "Ахуэнный товар"}`

	Testcases := []testCase{
		{Name: "Ad without name", In: []byte(AdWithoutName), ExpectedErr: true},
		{Name: "Ad with long namev", In: []byte(AdWithLongName), ExpectedErr: true},
		{Name: "Ad without description", In: []byte(AdWithoutDescription), ExpectedErr: true},
		{Name: "Ad with long description", In: []byte(AdWithLongDescription), ExpectedErr: true},
		{Name: "Ad without photots", In: []byte(AdWithoutPhotos), ExpectedErr: true},
		{Name: "Ad with many photots", In: []byte(AdWithManyPhotos), ExpectedErr: true},
		{Name: "Ad without photots", In: []byte(AdWithoutPhotos), ExpectedErr: true},
		{Name: "Ad with many photots", In: []byte(AdWithManyPhotos), ExpectedErr: true},
		{Name: "Ad without cost", In: []byte(AdWithoutCost), ExpectedErr: true},
		{Name: "Ad with negative cost", In: []byte(AdWithNegativeCost), ExpectedErr: true},
		{Name: "Invalid JSON", In: []byte(InvalidJson), ExpectedErr: true},
		{Name: "OK ad", In: []byte(AdJSON), ExpectedErr: false},
	}
	for _, tc := range Testcases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			reader := bytes.NewBuffer(tc.In)
			r := httptest.NewRequest("GET", URL, reader)
			var tmp Ad

			err := tmp.Bind(r)
			assert.Equal(t, tc.ExpectedErr, err != nil, "Binded item: %v", tmp)
		})
	}
}
