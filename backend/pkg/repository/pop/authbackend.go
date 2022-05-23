package pop

import (
	"errors"

	"github.com/apex/log"
	"github.com/gobuffalo/pop/v6"

	"github.com/spezifisch/rueder3/backend/pkg/authbackend/controller"
	"github.com/spezifisch/rueder3/backend/pkg/repository/pop/models"
)

// AuthBackendPopRepository internal state
type AuthBackendPopRepository struct {
	pop *pop.Connection
}

// NewAuthBackendPopRepository returns a AuthBackendPopRepository that wraps a pop DB
func NewAuthBackendPopRepository(db string) *AuthBackendPopRepository {
	// connect using pop for ORM stuff
	popTx, err := pop.Connect(db)
	if err != nil {
		log.WithError(err).WithField("db", db).Error("couldn't connect with pop")
		return nil
	}

	return &AuthBackendPopRepository{
		pop: popTx,
	}
}

func (r *AuthBackendPopRepository) getUser(authOrigin, authSubject string) (user models.User, exists bool, err error) {
	users := []models.User{}
	err = r.pop.Select("id").Where("auth_origin = ?", authOrigin).Where("auth_subject = ?", authSubject).Limit(1).All(&users)
	if err != nil {
		return
	}
	if len(users) != 1 {
		return
	}
	user = users[0]
	exists = true
	return
}

// GetOrCreateUser returns user details or creates the user if it doesn't exist
func (r *AuthBackendPopRepository) GetOrCreateUser(authOrigin, authSubject string) (ret controller.User, err error) {
	user, exists, err := r.getUser(authOrigin, authSubject)
	if err != nil {
		return
	}
	if !exists {
		log.WithField("origin", authOrigin).WithField("sub", authSubject).Info("creating new user")
		newUser := models.User{
			AuthOrigin:  authOrigin,
			AuthSubject: authSubject,
		}
		_, err = r.pop.ValidateAndCreate(&newUser)
		if err != nil {
			return
		}

		user, exists, err = r.getUser(authOrigin, authSubject)
		if err != nil {
			return
		}
		if !exists {
			err = errors.New("couldn't select user after creating it")
			return
		}
	}

	ret = controller.User{
		ID:          user.ID,
		AuthOrigin:  authOrigin,
		AuthSubject: authSubject,
	}
	return
}
