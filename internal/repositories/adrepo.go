package repositories

import (
	"context"

	"github.com/Reywaltz/avito_advertising/internal/models"
	"github.com/google/uuid"
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
		err = rows.Scan(&tmp.ID, &tmp.Name, &tmp.Description, &tmp.Photos, &tmp.Cost)
		if err != nil {
			return nil, err
		}
		out = append(out, tmp)
	}

	return out, nil
}

const (
	CreateAdQuery = `INSERT INTO advertisement VALUES ($1, $2, $3, $4, $5) returning id`
)

func (r *AdRepo) Create(newAd models.Ad) (uuid.UUID, error) {
	row := r.DB.QueryRow(context.Background(), CreateAdQuery, &newAd.ID, &newAd.Name, &newAd.Description, &newAd.Photos, &newAd.Cost)
	var tmp string
	if err := row.Scan(&tmp); err != nil {
		return uuid.Nil, err
	}

	out, err := uuid.Parse(tmp)
	if err != nil {
		return uuid.Nil, err
	}

	return out, nil
}

const (
	GetOneQuery = `SELECT id, name, description, photos, cost FROM advertisement WHERE id = $1`
)

func (r *AdRepo) GetOne(reqUUID uuid.UUID) (models.Ad, error) {
	row := r.DB.QueryRow(context.Background(), GetOneQuery, reqUUID)
	var out models.Ad
	if err := row.Scan(&out.ID, &out.Name, &out.Description, &out.Photos, &out.Cost); err != nil {
		return models.Ad{}, err
	}

	return out, nil
}
