package additions

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Query struct {
	Offset  int
	Cost    string
	Created string
}

func (q *Query) Bind(r *http.Request) error {
	query := r.URL.Query()

	offset, err := handleOffset(query)
	if err != nil {
		return err
	}
	q.Offset = offset

	cost, err := handleCost(query)
	if err != nil {
		return err
	}
	q.Cost = cost

	created, err := handleCreate(query)
	if err != nil {
		return err
	}
	q.Created = created

	return nil
}

func handleOffset(query url.Values) (int, error) {
	offset := query.Get("offset")
	if offset == "" {
		return 0, nil
	}

	value, err := strconv.Atoi(offset)
	if err != nil {
		return 0, err
	}
	if value < 0 {
		return 0, errors.New("Offset can't be negative")
	}

	return value, nil
}

func handleCost(query url.Values) (string, error) {
	cost := query.Get("cost")
	if cost == "" {
		return "", nil
	}

	if cost != "asc" && cost != "desc" {
		return "", errors.New("Wrong value of cost")
	}

	return cost, nil
}

func handleCreate(query url.Values) (string, error) {
	cost := query.Get("created")
	if cost == "" {
		return "", nil
	}

	if cost != "asc" && cost != "desc" {
		return "", errors.New("Wrong value of created")
	}

	return cost, nil
}
