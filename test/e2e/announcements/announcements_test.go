package announcements

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/SatorNetwork/sator-api/lib/rbac"
	announcement_repository "github.com/SatorNetwork/sator-api/svc/announcement/repository"
	"github.com/SatorNetwork/sator-api/test/app_config"
	"github.com/SatorNetwork/sator-api/test/framework/client"
	"github.com/SatorNetwork/sator-api/test/framework/client/announcement"
	"github.com/SatorNetwork/sator-api/test/framework/client/auth"
	"github.com/SatorNetwork/sator-api/test/framework/user"
)

func TestAdminEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer app_config.RunAndWait()()

	c := client.NewClient()

	signUpRequest := auth.RandomSignUpRequest()
	user := user.NewInitializedUser(signUpRequest, t)
	user.SetRole(rbac.RoleAdmin)
	user.RefreshToken()

	title1 := "title-1"
	description1 := "description-1"
	actionUrl1 := "action-url-1"
	startsAt1 := time.Now().Unix()
	endsAt1 := time.Now().Unix()
	resp, err := c.Announcement.CreateAnnouncement(user.AccessToken(), &announcement.CreateAnnouncementRequest{
		Title:       title1,
		Description: description1,
		ActionUrl:   actionUrl1,
		StartsAt:    startsAt1,
		EndsAt:      endsAt1,
	})
	require.NoError(t, err)
	require.NotNil(t, resp.ID)

	announcement1, err := c.Announcement.GetAnnouncementByID(user.AccessToken(), &announcement.GetAnnouncementByIDRequest{
		ID: resp.ID,
	})
	require.NoError(t, err)
	require.Equal(t, resp.ID, announcement1.ID)
	require.Equal(t, title1, announcement1.Title)
	require.Equal(t, description1, announcement1.Description)
	require.Equal(t, actionUrl1, announcement1.ActionUrl)
	require.Equal(t, startsAt1, announcement1.StartsAt)
	require.Equal(t, endsAt1, announcement1.EndsAt)

	announcements, err := c.Announcement.ListAnnouncements(user.AccessToken())
	require.NoError(t, err)
	announcement1, err = findAnnouncementByID(announcements, resp.ID)
	require.NoError(t, err)
	require.Equal(t, resp.ID, announcement1.ID)
	require.Equal(t, title1, announcement1.Title)
	require.Equal(t, description1, announcement1.Description)
	require.Equal(t, actionUrl1, announcement1.ActionUrl)
	require.Equal(t, startsAt1, announcement1.StartsAt)
	require.Equal(t, endsAt1, announcement1.EndsAt)

	titleUpd := "title-upd"
	descriptionUpd := "description-upd"
	actionUrlUpd := "action-url-upd"
	startsAtUpd := time.Now().Unix()
	endsAtUpd := time.Now().Unix()
	err = c.Announcement.UpdateAnnouncement(user.AccessToken(), &announcement.UpdateAnnouncementRequest{
		ID:          resp.ID,
		Title:       titleUpd,
		Description: descriptionUpd,
		ActionUrl:   actionUrlUpd,
		StartsAt:    startsAtUpd,
		EndsAt:      endsAtUpd,
	})
	require.NoError(t, err)

	announcement1, err = c.Announcement.GetAnnouncementByID(user.AccessToken(), &announcement.GetAnnouncementByIDRequest{
		ID: resp.ID,
	})
	require.NoError(t, err)
	require.Equal(t, resp.ID, announcement1.ID)
	require.Equal(t, titleUpd, announcement1.Title)
	require.Equal(t, descriptionUpd, announcement1.Description)
	require.Equal(t, actionUrlUpd, announcement1.ActionUrl)
	require.Equal(t, startsAtUpd, announcement1.StartsAt)
	require.Equal(t, endsAtUpd, announcement1.EndsAt)

	err = c.Announcement.DeleteAnnouncement(user.AccessToken(), &announcement.DeleteAnnouncementRequest{
		ID: resp.ID,
	})
	require.NoError(t, err)
	_, err = c.Announcement.GetAnnouncementByID(user.AccessToken(), &announcement.GetAnnouncementByIDRequest{
		ID: resp.ID,
	})
	require.Error(t, err)
}

func TestUserEndpoints(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer app_config.RunAndWait()()

	c := client.NewClient()
	err := c.DB.AnnouncementsDB().Repository().CleanUpReadAnnouncements(context.Background())
	require.NoError(t, err)
	err = c.DB.AnnouncementsDB().Repository().CleanUpAnnouncements(context.Background())
	require.NoError(t, err)

	user1 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)
	user1.SetRole(rbac.RoleAdmin)
	user1.RefreshToken()
	userID, err := c.DB.AuthDB().GetUserIDByEmail(context.Background(), user1.Email())
	require.NoError(t, err)
	user2 := user.NewInitializedUser(auth.RandomSignUpRequest(), t)

	var announcementID1 string
	{
		resp, err := c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-1",
			Description: "description-1",
			ActionUrl:   "action-url-1",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)

		isRead, err := c.DB.AnnouncementsDB().Repository().IsRead(context.Background(), announcement_repository.IsReadParams{
			AnnouncementID: uuid.MustParse(resp.ID),
			UserID:         userID,
		})
		require.NoError(t, err)
		require.False(t, isRead)
		isNotRead, err := c.DB.AnnouncementsDB().Repository().IsNotRead(context.Background(), announcement_repository.IsNotReadParams{
			AnnouncementID: uuid.MustParse(resp.ID),
			UserID:         userID,
		})
		require.NoError(t, err)
		require.True(t, isNotRead)

		err = c.Announcement.MarkAsRead(user1.AccessToken(), &announcement.MarkAsReadRequest{
			AnnouncementID: resp.ID,
		})
		require.NoError(t, err)

		isRead, err = c.DB.AnnouncementsDB().Repository().IsRead(context.Background(), announcement_repository.IsReadParams{
			AnnouncementID: uuid.MustParse(resp.ID),
			UserID:         userID,
		})
		require.NoError(t, err)
		require.True(t, isRead)
		isNotRead, err = c.DB.AnnouncementsDB().Repository().IsNotRead(context.Background(), announcement_repository.IsNotReadParams{
			AnnouncementID: uuid.MustParse(resp.ID),
			UserID:         userID,
		})
		require.NoError(t, err)
		require.False(t, isNotRead)

		announcementID1 = resp.ID
	}

	var announcementID2 string
	{
		resp, err := c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-2",
			Description: "description-2",
			ActionUrl:   "action-url-2",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)

		announcements, err := c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Equal(t, 1, len(announcements))

		err = c.Announcement.MarkAsRead(user1.AccessToken(), &announcement.MarkAsReadRequest{
			AnnouncementID: resp.ID,
		})
		require.NoError(t, err)

		announcements, err = c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Equal(t, 0, len(announcements))

		announcementID2 = resp.ID
	}

	{
		resp3, err := c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-3",
			Description: "description-3",
			ActionUrl:   "action-url-3",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)
		resp4, err := c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-4",
			Description: "description-4",
			ActionUrl:   "action-url-4",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)
		announcements, err := c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 2)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 4)

		markAsRead(t, c, user2, announcementID1)
		markAsRead(t, c, user2, announcementID2)
		markAsRead(t, c, user2, resp3.ID)
		markAsRead(t, c, user2, resp4.ID)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 2)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)

		markAsRead(t, c, user1, resp3.ID)
		markAsRead(t, c, user1, resp4.ID)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)
	}

	{
		_, err := c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-5",
			Description: "description-5",
			ActionUrl:   "action-url-5",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)
		_, err = c.Announcement.CreateAnnouncement(user1.AccessToken(), &announcement.CreateAnnouncementRequest{
			Title:       "title-6",
			Description: "description-6",
			ActionUrl:   "action-url-6",
			StartsAt:    time.Now().Unix(),
			EndsAt:      time.Now().Unix(),
		})
		require.NoError(t, err)
		announcements, err := c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 2)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 2)

		markAllAsRead(t, c, user1)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 2)

		markAllAsRead(t, c, user2)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user1.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)
		announcements, err = c.Announcement.ListUnreadAnnouncements(user2.AccessToken())
		require.NoError(t, err)
		require.Len(t, announcements, 0)
	}
}

func findAnnouncementByID(announcements []*announcement.Announcement, announcementID string) (*announcement.Announcement, error) {
	for _, a := range announcements {
		if a.ID == announcementID {
			return a, nil
		}
	}

	return nil, errors.Errorf("can't find announcement by ID")
}

func markAsRead(t *testing.T, c *client.Client, user *user.User, announcementID1 string) {
	err := c.Announcement.MarkAsRead(user.AccessToken(), &announcement.MarkAsReadRequest{
		AnnouncementID: announcementID1,
	})
	require.NoError(t, err)
}

func markAllAsRead(t *testing.T, c *client.Client, user *user.User) {
	err := c.Announcement.MarkAllAsRead(user.AccessToken())
	require.NoError(t, err)
}
