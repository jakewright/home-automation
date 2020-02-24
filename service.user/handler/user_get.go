package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	userdef "github.com/jakewright/home-automation/service.user/def"
)

// HandleGetUser reads a user by ID
func HandleGetUser(r *Request, body *userdef.GetUserRequest) (*userdef.GetUserResponse, error) {
	user := &userdef.User{}
	if err := database.Find(user, body.UserId); err != nil {
		return nil, errors.WithMessage(err, "failed to find")
	}

	return &userdef.GetUserResponse{
		User: user,
	}, nil
}
