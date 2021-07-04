package repositories

import (
	"context"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	QueryRow(ctx context.Context, sql string, arguments ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, arguments ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

type AdRepo struct {
	DB PgxIface
}

func NewRepo(db PgxIface) *AdRepo {
	return &AdRepo{
		DB: db,
	}
}

const (
	AdsFields = `name, description, photos, cost`
	AllFileds = `id, ` + AdsFields
	GetAds    = `SELECT ` + AllFileds + ` FROM advertisement`
)

func (r *AdRepo) GetAll() ([]models.Ad, error) {
	rows, err := r.DB.Query(context.Background(), GetAds)
	if err != nil {
		return nil, err
	}

	out := make([]models.Ad, 0)
	for rows.Next() {
		var tmp models.Ad
		rows.Scan(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Photos, &tmp.Cost)
		out = append(out, tmp)
	}

	return out, nil
}
