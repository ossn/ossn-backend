package ossn_backend

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/ossn/ossn-backend/models"
)

func min(limit *int, first *int) *int {
	switch {
	case limit == nil:
		return first
	case first == nil:
		return limit
	case *limit > *first:
		return first
	case *limit < *first:
		return limit
	default:
		return first
	}
}

func parseParams(query *gorm.DB, limit *int, after, before *string) (*gorm.DB, error) {
	l := 10
	if limit != nil {
		if *limit > 50 {
			l = 50
		}
		l = *limit
	}
	query = query.Limit(l)
	if after != nil {
		decoded, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return query, errors.New("Invalid after option")
		}
		id, err := strconv.Atoi(string(decoded[:]))
		if err != nil {
			return query, errors.New("Invalid after option")
		}
		query = query.Where("id < ?", id)
	}
	if before != nil {
		decoded, err := base64.StdEncoding.DecodeString(*before)
		if err != nil {
			return query, errors.New("Invalid before option")
		}
		id, err := strconv.Atoi(string(decoded[:]))
		if err != nil {
			return query, errors.New("Invalid before option")
		}
		query = query.Where("id > ?", id)
	}
	return query, nil
}

func getPageInfo(count *int, firstID, lastID *uint, first *int, length int) models.PageInfo {
	return models.PageInfo{
		TotalCount:      *count,
		HasPreviousPage: length < 1 && *count > 0 && int(*lastID) < *count,
		HasNextPage:     length < 1 && *count > 0 && int(*firstID) > 1,
		EndCursor:       base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(*firstID), 10))),
		StartCursor:     base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(*lastID), 10))),
	}
}
