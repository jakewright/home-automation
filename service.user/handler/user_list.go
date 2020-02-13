package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/errors"
	userdef "github.com/jakewright/home-automation/service.user/def"
)

// HandleListUsers lists all users
func HandleListUsers(_ *userdef.ListUsersRequest) (*userdef.ListUsersResponse, error) {
	var users []userdef.User
	if err := database.Find(&users); err != nil {
		return nil, errors.WithMessage(err, "failed to find")
	}

	return &userdef.ListUsersResponse{
		Users: users,
	}, nil
}
