package pg

import "github.com/jackc/pgtype"

func ToJsonb[T any](payload *T) (*pgtype.JSONB, error) {
	if payload == nil {
		return getEmptyJSON()
	}
	jsonb := pgtype.JSONB{}
	err := jsonb.Set(payload)
	if err != nil {
		return nil, err
	}
	return &jsonb, nil
}

func FromJsonb[T any](j *pgtype.JSONB) (*T, error) {
	if j == nil {
		return nil, nil
	}
	var v T
	err := j.AssignTo(&v)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func getEmptyJSON() (*pgtype.JSONB, error) {
	emptyJson := pgtype.JSONB{}
	err := emptyJson.Scan("{}")
	if err != nil {
		return nil, err
	}
	return &emptyJson, nil
}
