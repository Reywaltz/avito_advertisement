package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Ad struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Photos      []string        `json:"photos"`
	Cost        decimal.Decimal `json:"cost"`
	Created     time.Time       `json:"created"`
}

type AdMainPhoto struct {
	ID          uuid.UUID       `json:"-"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Photo       string          `json:"photo"`
	Cost        decimal.Decimal `json:"cost"`
	Created     time.Time       `json:"-"`
}

func (a *Ad) Bind(r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &a)
	if err != nil {
		return err
	}

	if a.Name == "" || len([]rune(a.Name)) > 200 {
		return errors.New("Name can't be binded")
	}

	if a.Description == "" || len([]rune(a.Description)) > 1000 {
		return errors.New("Description can't be binded")
	}

	if a.Photos == nil || len(a.Photos) > 3 {
		return errors.New("Photos can't be binded")
	}

	if a.Cost.IsZero() || a.Cost.IsNegative() {
		return errors.New("Cost can't be binded")
	}

	return nil
}
