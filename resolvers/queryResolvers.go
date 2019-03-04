package resolvers

import (
	"context"
	"strings"

	"github.com/ossn/ossn-backend/helpers"
	"github.com/ossn/ossn-backend/models"
)

type queryResolver struct{ *Resolver }

func (q *queryResolver) Session(ctx context.Context) (*models.User, error) {
	return helpers.GetUserFromContext(ctx)
}

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := models.DBSession.Where("id = ?", id).First(user).Error
	return user, err
}

func (r *queryResolver) Users(ctx context.Context, first, last *int, before, after, search *string) (*models.Users, error) {
	err := validateFirstAndLast(first, last)
	if err != nil {
		return nil, err
	}
	query := models.DBSession
	count := 0
	if search != nil {
		str := "%" + strings.ToLower(*search) + "%"
		query = query.Where("name ILIKE ? OR user_name ILIKE ?", str, str)
	}
	err = query.Find(&[]models.User{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	query, err = parseParams(query, first, last, after, before, "name")
	if err != nil {
		return nil, err
	}
	users := []models.User{}
	err = query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	length := len(users)
	switch {
	case length < 1:
		return &models.Users{Users: users, PageInfo: models.PageInfo{
			TotalCount:      count,
			HasPreviousPage: count > 0,
			HasNextPage:     false,
			StartCursor:     "",
			EndCursor:       "",
		}}, err
	case length > 1 && length >= (*max(first, last)+1):
		users = users[:length-1]
	}
	firstID := &users[len(users)-1].ID
	lastID := &users[0].ID
	return &models.Users{
		Users:    users,
		PageInfo: getPageInfo(&count, firstID, lastID, first, last, length),
	}, err
}

func (r *queryResolver) Clubs(ctx context.Context, first, last *int, userID *string, ids []*string, before, after, search *string) (*models.Clubs, error) {
	err := validateFirstAndLast(first, last)
	if err != nil {
		return nil, err
	}
	count := 0
	safeIDS := []string{}
	for _, id := range ids {
		if id != nil {
			safeIDS = append(safeIDS, *id)
		}
	}

	query := models.DBSession
	if search != nil {
		query = query.Where("title ILIKE ?", "%"+strings.ToLower(*search)+"%")
	}
	if len(safeIDS) > 0 {
		query = query.Where("id in (?)", safeIDS)
	}
	if userID != nil {
		query = query.Where("id in (SELECT club_id from club_user_roles where user_id = ?)", userID)
	}
	err = query.Find(&[]models.Club{}).Count(&count).Error
	if err != nil {
		return nil, err
	}

	query, err = parseParams(query, first, last, after, before, "title")
	if err != nil {
		return nil, err
	}
	clubs := []models.Club{}

	err = query.Preload("Location").Preload("Events").Preload("Events.Location").Find(&clubs).Error
	if err != nil {
		return nil, err
	}

	length := len(clubs)
	switch {
	case length < 1:
		return &models.Clubs{
			Clubs: clubs,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	case length >= (*max(first, last) + 1):
		clubs = clubs[:length-1]
	}

	firstID := &clubs[len(clubs)-1].ID
	lastID := &clubs[0].ID
	return &models.Clubs{
		Clubs:    clubs,
		PageInfo: getPageInfo(&count, firstID, lastID, first, last, length),
	}, err
}

func (r *queryResolver) Club(ctx context.Context, id string) (*models.Club, error) {
	club := &models.Club{}
	err := models.DBSession.Where("id = ?", id).Preload("Location").Preload("Events").First(club).Error
	return club, err
}

func (r *queryResolver) Events(ctx context.Context, first, last *int, clubId, before, after *string) (*models.Events, error) {
	err := validateFirstAndLast(first, last)
	if err != nil {
		return nil, err
	}
	query := models.DBSession
	count := 0
	if clubId != nil {
		query = query.Where("club_id = ?", clubId)
	}
	err = models.DBSession.Find(&[]models.Event{}).Count(&count).Error
	if err != nil {
		return nil, err
	}

	query, err = parseParams(query, first, last, after, before, "published_at")
	if err != nil {
		return nil, err
	}

	events := []models.Event{}
	err = query.Find(&events).Error
	if err != nil {
		return nil, err
	}

	length := len(events)
	switch {
	case length < 1:
		return &models.Events{
			Events: events,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	case length > 1 && length >= (*max(first, last)+1):
		events = events[:length-1]
	}

	firstID := &events[len(events)-1].ID
	lastID := &events[0].ID
	return &models.Events{
		Events:   events,
		PageInfo: getPageInfo(&count, firstID, lastID, first, last, length),
	}, err

}

func (r *queryResolver) Event(ctx context.Context, id string) (*models.Event, error) {
	event := &models.Event{}
	err := models.DBSession.Where("id =?", id).First(event).Error
	return event, err
}

func (r *queryResolver) Jobs(ctx context.Context, first, last *int, before, after *string) (*models.Jobs, error) {
	err := validateFirstAndLast(first, last)
	if err != nil {
		return nil, err
	}
	query := models.DBSession
	count := 0
	err = models.DBSession.Find(&[]models.Job{}).Count(&count).Error
	if err != nil {
		return nil, err
	}

	query, err = parseParams(query, first, last, after, before, "published_at")
	if err != nil {
		return nil, err
	}

	jobs := []models.Job{}
	err = query.Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	length := len(jobs)
	switch {
	case length < 1:
		return &models.Jobs{
			Jobs: jobs,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	case length > 1 && length >= (*max(first, last)+1):
		jobs = jobs[:length-1]
	}
	firstID := &jobs[len(jobs)-1].ID
	lastID := &jobs[0].ID
	return &models.Jobs{
		Jobs:     jobs,
		PageInfo: getPageInfo(&count, firstID, lastID, first, last, length),
	}, err
}

func (r *queryResolver) Announcements(ctx context.Context, first, last *int, before *string, after *string) (*models.Announcements, error) {
	err := validateFirstAndLast(first, last)
	if err != nil {
		return nil, err
	}
	query := models.DBSession
	count := 0
	err = models.DBSession.Find(&[]models.Announcement{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	query, err = parseParams(query, first, last, after, before, "published_at")
	if err != nil {
		return nil, err
	}
	announcements := []models.Announcement{}
	err = query.Find(&announcements).Error
	if err != nil {
		return nil, err
	}

	length := len(announcements)
	switch {
	case length < 1:
		return &models.Announcements{
			Announcements: announcements,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	case length > 1 && length >= (*max(first, last)+1):
		announcements = announcements[:length-1]
	}

	firstID := &announcements[len(announcements)-1].ID
	lastID := &announcements[0].ID
	return &models.Announcements{
		Announcements: announcements,
		PageInfo:      getPageInfo(&count, firstID, lastID, first, last, length),
	}, err
}
