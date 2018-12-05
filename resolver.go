//go:generate go run scripts/gqlgen.go -v
package ossn_backend

import (
	"context"
	"strconv"

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
			ID:                strconv.Itoa(u.ID),
			Email:             u.Email,
			ImageURL:          u.ImageURL,
			Role:              &models.Role{Name: models.TurnStringToRolename(user.Role)},
			GithubURL:         u.GithubURL,
			UpdatedAt:         u.UpdatedAtToString(),
			CreatedAt:         u.CreatedAtToString(),
			Description:       u.Description,
			Clubs:             clubs,
			UserName:          u.UserName,
			FirstName:         u.FirstName,
			LastName:          u.LastName,
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

func (r *mutationResolver) CreateUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateClub(ctx context.Context, input models.ClubInput) (*models.Club, error) {
	panic("not implemented")
}
func (r *mutationResolver) CreateLocation(ctx context.Context, input *models.LocationInput) (*models.Location, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	err := models.DBSession.Where("id = ?", id).First(user).Error
	return user, err
}
func (r *queryResolver) Users(ctx context.Context, first *int, before *string, after *string, limit *int) (*models.Users, error) {
	query := models.DBSession
	query, err := parseParams(query, min(first, limit), after, before)
	if err != nil {
		return nil, err
	}
	users := []models.User{}
	count := 0
	err = query.Find(&users).Limit(-1).Offset(-1).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if len(users) < 1 {
		return &models.Users{Users: users, PageInfo: models.PageInfo{
			TotalCount:      count,
			HasPreviousPage: count > 0,
			HasNextPage:     false,
			StartCursor:     "",
			EndCursor:       "",
		}}, err
	}
	firstID := &users[len(users)-1].ID
	lastID := &users[0].ID
	return &models.Users{
		Users:    users,
		PageInfo: getPageInfo(&count, firstID, lastID, min(limit, first), 0),
	}, err
}

func (r *queryResolver) Clubs(ctx context.Context, first *int, userID *string, ids []*string, before *string, after *string, limit *int) (*models.Clubs, error) {
	query := models.DBSession
	query, err := parseParams(query, min(first, limit), after, before)
	if err != nil {
		return nil, err
	}
	clubs := []models.Club{}
	i := []string{}
	for _, id := range ids {
		if id != nil {
			i = append(i, *id)
		}
	}
	if len(i) > 0 {
		query = query.Where("id in (?)", i)
	}
	if userID != nil {
		query = query.Where("id in (SELECT club_id from club_user_roles where user_id = ?)", userID)
	}
	err = query.Find(&clubs).Error
	if err != nil {
		return nil, err
	}

	count := 0
	err = models.DBSession.Find(&[]models.Club{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if len(clubs) < 1 {
		return &models.Clubs{
			Clubs: clubs,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	}

	firstID := &clubs[len(clubs)-1].ID
	lastID := &clubs[0].ID
	return &models.Clubs{
		Clubs:    clubs,
		PageInfo: getPageInfo(&count, firstID, lastID, min(limit, first), len(i)),
	}, err
}

func (r *queryResolver) Club(ctx context.Context, id string) (*models.Club, error) {
	club := &models.Club{}
	err := models.DBSession.Where("id =?", id).First(club).Error
	return club, err
}

func (r *queryResolver) Events(ctx context.Context, first *int, clubId *string, before *string, after *string, limit *int) (*models.Events, error) {
	query := models.DBSession.Order("id desc, published_at")
	query, err := parseParams(query, min(first, limit), after, before)
	if err != nil {
		return nil, err
	}

	events := []models.Event{}
	if clubId != nil {
		query = query.Where("club_id =", clubId)
	}
	err = query.Find(&events).Error
	if err != nil {
		return nil, err
	}

	count := 0
	err = models.DBSession.Find(&[]models.Event{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if len(events) < 1 {
		return &models.Events{
			Events: events,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	}

	firstID := &events[len(events)-1].ID
	lastID := &events[0].ID
	return &models.Events{
		Events:   events,
		PageInfo: getPageInfo(&count, firstID, lastID, min(limit, first), 0),
	}, err

}

func (r *queryResolver) Event(ctx context.Context, id string) (*models.Event, error) {
	event := &models.Event{}
	err := models.DBSession.Where("id =?", id).First(event).Error
	return event, err
}

func (r *queryResolver) Jobs(ctx context.Context, first *int, before *string, after *string, limit *int) (*models.Jobs, error) {
	query := models.DBSession.Order("id desc, published_at")
	query, err := parseParams(query, min(first, limit), after, before)
	if err != nil {
		return nil, err
	}

	jobs := []models.Job{}
	err = query.Find(&jobs).Error
	if err != nil {
		return nil, err
	}

	count := 0
	err = models.DBSession.Find(&[]models.Job{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if len(jobs) < 1 {
		return &models.Jobs{
			Jobs: jobs,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	}

	firstID := &jobs[len(jobs)-1].ID
	lastID := &jobs[0].ID
	return &models.Jobs{
		Jobs:     jobs,
		PageInfo: getPageInfo(&count, firstID, lastID, min(limit, first), 0),
	}, err
}

func (r *queryResolver) Announcements(ctx context.Context, first *int, before *string, after *string, limit *int) (*models.Announcements, error) {
	query := models.DBSession.Order("id desc, published_at")
	query, err := parseParams(query, min(first, limit), after, before)
	if err != nil {
		return nil, err
	}
	announcements := []models.Announcement{}
	err = query.Find(&announcements).Error
	if err != nil {
		return nil, err
	}

	count := 0
	err = models.DBSession.Find(&[]models.Announcement{}).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if len(announcements) < 1 {
		return &models.Announcements{
			Announcements: announcements,
			PageInfo: models.PageInfo{
				TotalCount:      count,
				HasPreviousPage: count > 0,
				HasNextPage:     false,
				StartCursor:     "",
				EndCursor:       "",
			}}, err
	}

	firstID := &announcements[len(announcements)-1].ID
	lastID := &announcements[0].ID
	return &models.Announcements{
		Announcements: announcements,
		PageInfo:      getPageInfo(&count, firstID, lastID, min(limit, first), 0),
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
			ID:            strconv.Itoa(c.ID),
			Email:         c.Email,
			Location:      c.Location,
			Name:          c.Title,
			ImageURL:      c.ImageURL,
			Role:          &models.Role{Name: models.TurnStringToRolename(club.Role)},
			GithubURL:     c.GithubURL,
			UpdatedAt:     c.UpdatedAtToString(),
			CreatedAt:     c.CreatedAtToString(),
			Description:   c.Description,
			CodeOfConduct: c.CodeOfConduct,
			ClubURL:       c.ClubURL,
			Events:        c.Events,
			Users:         users,
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
