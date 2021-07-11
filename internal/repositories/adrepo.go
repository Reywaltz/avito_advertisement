package repositories

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Reywaltz/avito_advertising/cmd/advert-api/additions"
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
	DefaultLimit           = `10`
	AdsFieldsWithMainPhoto = `name, photos[1], cost, created`
)

var GetAds = `SELECT ` + AdsFieldsWithMainPhoto + ` FROM advertisement limit ` + DefaultLimit + ` offset $1`

func (r *AdRepo) GetAll(queries additions.Query) ([]models.AdMainPhoto, error) {
	limit, _ := strconv.Atoi(DefaultLimit)

	getquery := BuildSQLQuery(queries)

	rows, err := r.DB.Query(context.Background(), getquery, queries.Offset*limit)
	if err != nil {
		return nil, err
	}

	out := make([]models.AdMainPhoto, 0)
	for rows.Next() {
		var tmp models.AdMainPhoto
		err := rows.Scan(&tmp.Name, &tmp.Photo, &tmp.Cost, &tmp.Created)
		if err != nil {
			return nil, err
		}
		out = append(out, tmp)
	}

	return out, nil
}

const (
	CreateAdQuery = `INSERT INTO advertisement VALUES ($1, $2, $3, $4, $5, $6) returning id`
)

func (r *AdRepo) Create(newAd models.Ad) (uuid.UUID, error) {
	row := r.DB.QueryRow(context.Background(), CreateAdQuery,
		&newAd.ID, &newAd.Name, &newAd.Description,
		&newAd.Photos, &newAd.Cost, &newAd.Created)
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
	GetOneQuery = `SELECT id, name, description, photos, 
	cost FROM advertisement WHERE id = $1`
)

func (r *AdRepo) GetOne(reqUUID uuid.UUID) (models.Ad, error) {
	row := r.DB.QueryRow(context.Background(), GetOneQuery, reqUUID)
	var out models.Ad
	if err := row.Scan(&out.ID, &out.Name, &out.Description, &out.Photos,
		&out.Cost, out.Created); err != nil {
		return models.Ad{}, err
	}

	return out, nil
}

func BuildSQLQuery(queries additions.Query) string {
	if queries.Cost != "" && queries.Created != "" {
		GetAds = fmt.Sprintf(`SELECT * from (SELECT name, photos[1], cost, created FROM advertisement order by cost %s) as a order by a.created %s limit `+DefaultLimit+` offset $1`, queries.Cost, queries.Created)

		return GetAds
	}

	if queries.Cost != "" {
		GetAds = fmt.Sprintf(`SELECT name, photos[1], cost, created FROM advertisement order by cost %s limit `+DefaultLimit+` offset $1`, queries.Cost)

		return GetAds
	}
	if queries.Created != "" {
		GetAds = fmt.Sprintf(`SELECT name, photos[1], cost, created FROM advertisement order by created %s limit `+DefaultLimit+` offset $1`, queries.Cost)

		return GetAds
	}

	return GetAds
}
