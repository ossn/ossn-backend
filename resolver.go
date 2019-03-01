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
		err := models.DBSession.Raw("Select * from clubs where id IN (SELECT club_id from club_user_roles where user_id = ?)", u.ID).Scan(&clubs).Error
		if err != nil {
			return usersWithRole, err
		}
		usersWithRole = append(usersWithRole, &models.UserWithRole{
			ID:                  strconv.FormatUint(uint64(u.ID), 10),
			Email:               u.Email,
			ImageURL:            u.ImageURL,
			Role:                models.TurnStringToRolename(user.Role),
			GithubURL:           u.GithubURL,
			UpdatedAt:           u.UpdatedAtToString(),
			CreatedAt:           u.CreatedAtToString(),
			Description:         u.Description,
			Clubs:               clubs,
			UserName:            u.UserName,
			Name:                u.Name,
			ReceiveNewsletter:   u.ReceiveNewsletter,
			SortDescription:     u.SortDescription,
			PersonalURL:         u.PersonalURL,
			IsOverTheLegalLimit: u.IsOverTheLegalLimit,
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
