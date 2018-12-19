package ossn_backend

import (
	"encoding/base64"
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/ossn/ossn-backend/models"
)

func validateFirstAndLast(first, last *int) error {
	switch {
	case (first == nil && last == nil):
		return errors.New("Please provide first or last param")
	case (first != nil && last != nil):
		return errors.New("Please provide only first or only last")
	case (first != nil && *first > 100) || (last != nil && *last > 100):
		return errors.New("First and last can't exceed 100")
	default:
		return nil
	}

}

func max(first, second *int) *int {
	switch {
	case first == nil:
		return second
	case second == nil:
		return first
	case *first > *second:
		return first
	case *first < *second:
		return second
	default:
		zero := 0
		return &zero
	}
}

func parseBase64Str(str *string) (int, error) {
	decoded, err := base64.StdEncoding.DecodeString(*str)
	if err != nil {
		return 0, err
	}
	num, err := strconv.Atoi(string(decoded[:]))
	return num, err
}

func parseParams(query *gorm.DB, first, last *int, after, before *string, orederByCol string) (*gorm.DB, error) {
	l := 0

	if first != nil {
		query = query.Order("id desc, " + orederByCol + " desc")
		l = *first
	}
	if last != nil {
		query = query.Order("id asc, " + orederByCol + " asc")
		l = *last
	}
	query = query.Limit(l + 1)

	if after != nil {
		id, err := parseBase64Str(after)
		if err != nil {
			return query, errors.New("Invalid after option")
		}
		query = query.Where("id > ?", id)
	}

	if before != nil {
		id, err := parseBase64Str(before)
		if err != nil {
			return query, errors.New("Invalid before option")
		}
		query = query.Where("id < ?", id)
	}
	return query, nil
}

func getPageInfo(count *int, firstID, lastID *uint, first, last *int, length int) models.PageInfo {
	hasNext := false
	hasPrev := false

	switch {
	case first != nil:
		hasNext = length == (*first + 1)
	case last != nil:
		hasPrev = length == (*last + 1)
	}

	return models.PageInfo{
		TotalCount:      *count,
		HasPreviousPage: hasPrev,
		HasNextPage:     hasNext,
		EndCursor:       base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(*firstID), 10))),
		StartCursor:     base64.StdEncoding.EncodeToString([]byte(strconv.FormatUint(uint64(*lastID), 10))),
	}
}
