//go:generate go run scripts/gqlgen.go -v
package ossn_backend

import (
	"context"
	"strconv"
	"strings"

	"github.com/ossn/ossn-backend/helpers"

	"github.com/ossn/ossn-backend/models"
)

type Resolver struct{}

func (r *Resolver) Announcement() AnnouncementResolver {
	return &announcementResolver{r}
}
func (r *Resolver) Club() ClubResolver {
	return &clubResolver{r}
}
func (r *Resolver) Event() EventResolver {
	return &eventResolver{r}
}
func (r *Resolver) Job() JobResolver {
	return &jobResolver{r}
}
func (r *Resolver) Location() LocationResolver {
	return &locationResolver{r}
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) User() UserResolver {
	return &userResolver{r}
}

type announcementResolver struct{ *Resolver }

func (r *announcementResolver) ID(ctx context.Context, obj *models.Announcement) (string, error) {
	return obj.IDToString()
}
func (r *announcementResolver) PublishedAt(ctx context.Context, obj *models.Announcement) (*string, error) {
	return obj.PublishedAtToString()
}
func (r *announcementResolver) CreatedAt(ctx context.Context, obj *models.Announcement) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *announcementResolver) UpdatedAt(ctx context.Context, obj *models.Announcement) (string, error) {
	return obj.UpdatedAtToString(), nil
}

type clubResolver struct{ *Resolver }

func (r *clubResolver) ID(ctx context.Context, obj *models.Club) (string, error) {
	return obj.IDToString()
}
func (r *clubResolver) Name(ctx context.Context, obj *models.Club) (*string, error) {
	return obj.Title, nil
}
func (r *clubResolver) Users(ctx context.Context, obj *models.Club) ([]*models.UserWithRole, error) {

	clubUserRole := []*models.ClubUserRole{}
	usersWithRole := []*models.UserWithRole{}

	err := models.DBSession.Preload("User").Where("club_id = ?", obj.ID).Find(&clubUserRole).Error
	if err != nil {
		return usersWithRole, err
	}

	for _, user := range clubUserRole {
		u := &user.User
		// TODO: Improve this
		clubs := []*models.ClubWithRole{}
		err := models.DBSession.Raw("SELECT * FROM users where id IN (SELECT user_id from club_user_roles where club_id = ?)", u.ID).Scan(&clubs).Error
		if err != nil {
			return usersWithRole, err
		}
		usersWithRole = append(usersWithRole, &models.UserWithRole{
			ID:                strconv.FormatUint(uint64(u.ID), 10),
			Email:             u.Email,
			ImageURL:          u.ImageURL,
			Role:              models.TurnStringToRolename(user.Role),
			GithubURL:         u.GithubURL,
			UpdatedAt:         u.UpdatedAtToString(),
			CreatedAt:         u.CreatedAtToString(),
			Description:       u.Description,
			Clubs:             clubs,
			UserName:          u.UserName,
			Name:              u.Name,
			ReceiveNewsletter: u.ReceiveNewsletter,
			SortDescription:   u.SortDescription,
			PersonalURL:       u.PersonalURL,
		})
	}
	return usersWithRole, nil
}
func (r *clubResolver) CreatedAt(ctx context.Context, obj *models.Club) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *clubResolver) UpdatedAt(ctx context.Context, obj *models.Club) (string, error) {
	return obj.UpdatedAtToString(), nil
}

type eventResolver struct{ *Resolver }

func (r *eventResolver) ID(ctx context.Context, obj *models.Event) (string, error) {
	return obj.IDToString()
}
func (r *eventResolver) StartDate(ctx context.Context, obj *models.Event) (*string, error) {
	return obj.StartDateToString()
}
func (r *eventResolver) EndDate(ctx context.Context, obj *models.Event) (*string, error) {
	return obj.EndDateToString()
}
func (r *eventResolver) PublishedAt(ctx context.Context, obj *models.Event) (*string, error) {
	return obj.PublishedAtToString()
}
func (r *eventResolver) CreatedAt(ctx context.Context, obj *models.Event) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *eventResolver) UpdatedAt(ctx context.Context, obj *models.Event) (string, error) {
	return obj.UpdatedAtToString(), nil
}

type jobResolver struct{ *Resolver }

func (r *jobResolver) ID(ctx context.Context, obj *models.Job) (string, error) {
	return obj.IDToString()
}
func (r *jobResolver) PublishedAt(ctx context.Context, obj *models.Job) (*string, error) {
	return obj.PublishedAtToString()
}
func (r *jobResolver) CreatedAt(ctx context.Context, obj *models.Job) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *jobResolver) UpdatedAt(ctx context.Context, obj *models.Job) (string, error) {
	return obj.UpdatedAtToString(), nil
}

type locationResolver struct{ *Resolver }

func (r *locationResolver) ID(ctx context.Context, obj *models.Location) (string, error) {
	return obj.IDToString()
}
func (r *locationResolver) CreatedAt(ctx context.Context, obj *models.Location) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *locationResolver) UpdatedAt(ctx context.Context, obj *models.Location) (string, error) {
	return obj.UpdatedAtToString(), nil
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) EditUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	user, err := helpers.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user.GithubURL = input.GithubURL
	user.PersonalURL = input.PersonalURL
	user.ReceiveNewsletter = &input.ReceiveNewsletter
	user.SortDescription = input.SortDescription
	user.Description = input.Description

	tx := models.DBSession.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Unscoped().Where("user_id = ? and club_id NOT IN (?)", user.ID, input.Clubs).Delete(&models.ClubUserRole{}).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	clubs := []*models.ClubUserRole{}
	for _, clubID := range input.Clubs {
		id, err := strconv.Atoi(clubID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		club := models.ClubUserRole{
			Role:   "member",
			UserID: user.ID,
			ClubID: uint(id),
		}

		err = tx.Preload("Club").Where("user_id = ? and club_id = ?", user.ID, id).FirstOrCreate(&club).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		club.User = *user
		clubs = append(clubs, &club)
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	user.Clubs = clubs
	return user, nil
}
func (r *mutationResolver) CreateClub(ctx context.Context, input models.ClubInput) (*models.Club, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateLocation(ctx context.Context, input models.LocationInput) (*models.Location, error) {
	panic("not implemented")
}
func (r *mutationResolver) JoinClub(ctx context.Context, clubID string) (bool, error) {
	id, err := strconv.Atoi(clubID)
	if err != nil {
		return false, err
	}
	user, err := helpers.GetUserFromContext(ctx)
	if err != nil {
		return false, err
	}
	clubUser := &models.ClubUserRole{
		UserID: user.ID,
		ClubID: uint(id),
		Role:   "member",
	}
	err = models.DBSession.Save(clubUser).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	session, err := helpers.GetSessionFromContext(ctx)
	if err != nil {
		return false, err
	}
	err = models.DBSession.Unscoped().Where("id = ?", session.ID).Delete(&models.Session{}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

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

	err = query.Find(&clubs).Error
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
	err := models.DBSession.Where("id = ?", id).First(club).Error
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

type userResolver struct{ *Resolver }

func (r *userResolver) ID(ctx context.Context, obj *models.User) (string, error) {
	return obj.IDToString()
}
func (r *userResolver) Clubs(ctx context.Context, obj *models.User) ([]*models.ClubWithRole, error) {
	clubUserRole := []models.ClubUserRole{}
	clubWithRole := []*models.ClubWithRole{}

	err := models.DBSession.Preload("Club").Preload("Club.Location").Preload("Club.Events").Where("user_id = ?", obj.ID).Find(&clubUserRole).Error
	if err != nil {
		return clubWithRole, err
	}

	for _, club := range clubUserRole {
		c := &club.Club
		users := []*models.User{}

		err := models.DBSession.Raw("SELECT * FROM users where id IN (SELECT user_id from club_user_roles where club_id = ?)", c.ID).Scan(&users).Error
		if err != nil {
			return clubWithRole, err
		}
		clubWithRole = append(clubWithRole, &models.ClubWithRole{
			ID:             strconv.FormatUint(uint64(c.ID), 10),
			Email:          c.Email,
			Location:       c.Location,
			Name:           c.Title,
			ImageURL:       c.ImageURL,
			Role:           models.TurnStringToRolename(club.Role),
			GithubURL:      c.GithubURL,
			UpdatedAt:      c.UpdatedAtToString(),
			CreatedAt:      c.CreatedAtToString(),
			Description:    c.Description,
			CodeOfConduct:  c.CodeOfConduct,
			ClubURL:        c.ClubURL,
			Events:         c.Events,
			Users:          users,
			BannerImageURL: c.BannerImageURL,
		})
	}
	return clubWithRole, nil
}
func (r *userResolver) CreatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.CreatedAtToString(), nil
}
func (r *userResolver) UpdatedAt(ctx context.Context, obj *models.User) (string, error) {
	return obj.UpdatedAtToString(), nil
}
