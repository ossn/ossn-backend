package ossn_backend

import (
	"context"
	"errors"
	"strconv"

	"github.com/ossn/ossn-backend/helpers"
	"github.com/ossn/ossn-backend/models"
)

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) EditUser(ctx context.Context, input models.UserInput) (*models.User, error) {
	user, err := helpers.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}
	user.Name = input.Name
	user.GithubURL = input.GithubURL
	user.PersonalURL = input.PersonalURL
	user.ReceiveNewsletter = &input.ReceiveNewsletter
	user.SortDescription = input.SortDescription
	user.Description = input.Description
	user.IsOverTheLegalLimit = input.IsOverTheLegalLimit

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

	query := tx.Unscoped().Where("user_id = ?", user.ID)
	if len(input.Clubs) > 0 {
		query = query.Where("club_id NOT IN (?)", input.Clubs)
	}
	err = query.Delete(&models.ClubUserRole{}).Error
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

func (r *mutationResolver) EditClub(ctx context.Context, clubID string, input models.ClubInput) (*models.Club, error) {
	user, err := helpers.GetUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	count := 0
	err = models.DBSession.Where("user_id = ? and club_id = ? and (role = 'admin' or role = 'club_owner')", user.ID, clubID).Table("club_user_roles").Count(&count).Error
	if err != nil || count < 1 {
		return nil, errors.New("You don't have permission to edit this club")
	}

	club := &models.Club{
		Email:           &input.Email,
		ImageURL:        input.ImageURL,
		Title:           &input.Name,
		Location:        nil,
		Description:     &input.Description,
		CodeOfConduct:   input.CodeOfConduct,
		SortDescription: input.SortDescription,
		GithubURL:       input.GithubURL,
		ClubURL:         input.ClubURL,
		BannerImageURL:  input.BannerImageURL,
	}

	id, err := strconv.Atoi(clubID)
	if err != nil {
		return nil, err
	}
	club.ID = uint(id)

	tx := models.DBSession.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if input.Location != nil {
		loc := models.Location{
			Address: input.Location.Address,
			Lat:     input.Location.Lat,
			Lng:     input.Location.Lng,
		}
		err = tx.FirstOrCreate(&loc, loc).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		club.LocationID = &loc.ID
		club.Location = &loc
	}

	err = tx.Save(&club).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	return club, nil
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

	count := 0
	err = models.DBSession.Where("user_id = ? and club_id = ?", user.ID, clubID).Table("club_user_roles").Count(&count).Error
	if err != nil {
		return false, errors.New("You don't have permission to edit this club")
	}
	if count > 0 {
		return false, errors.New("User is already a member of this club")
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
	err = models.RedisClient.Del(models.SESSION_PREFIX + session.Token).Err()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) Event(ctx context.Context, eventID *string, input *models.EventInput) (*models.Event, error) {
	panic("not implemented")
}

func (r *mutationResolver) CreateLocation(ctx context.Context, input models.LocationInput) (*models.Location, error) {
	panic("not implemented")
}
