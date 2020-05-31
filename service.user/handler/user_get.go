package handler

import (
	"github.com/jakewright/home-automation/libraries/go/database"
	"github.com/jakewright/home-automation/libraries/go/oops"
	userdef "github.com/jakewright/home-automation/service.user/def"
)

// GetUser reads a user by ID
func (h *Handler) GetUser(r *Request, body *userdef.GetUserRequest) (*userdef.GetUserResponse, error) {
	user := &userdef.User{}
	if err := database.Find(user, body.UserId); err != nil {
		return nil, oops.WithMessage(err, "failed to find")
	}

	return &userdef.GetUserResponse{
		User: user,
	}, nil
}
