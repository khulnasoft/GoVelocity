package services

import (
	"errors"

	"numtostr/gotodo/app/dal"
	"numtostr/gotodo/app/types"
	"numtostr/gotodo/utils"
	"numtostr/gotodo/utils/jwt"
	"numtostr/gotodo/utils/password"

	"go.khulnasoft.com/velocity"
	"gorm.io/gorm"
)

// Login service logs in a user
func Login(ctx *velocity.Ctx) error {
	b := new(types.LoginDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}

	u := &types.UserResponse{}

	err := dal.FindUserByEmail(u, b.Email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return velocity.NewError(velocity.StatusUnauthorized, "Invalid email or password")
	}

	if err := password.Verify(u.Password, b.Password); err != nil {
		return velocity.NewError(velocity.StatusUnauthorized, "Invalid email or password")
	}

	t := jwt.Generate(&jwt.TokenPayload{
		ID: u.ID,
	})

	return ctx.JSON(&types.AuthResponse{
		User: u,
		Auth: &types.AccessResponse{
			Token: t,
		},
	})
}

// Signup service creates a user
func Signup(ctx *velocity.Ctx) error {
	b := new(types.SignupDTO)

	if err := utils.ParseBodyAndValidate(ctx, b); err != nil {
		return err
	}

	err := dal.FindUserByEmail(&struct{ ID string }{}, b.Email).Error

	// If email already exists, return
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return velocity.NewError(velocity.StatusConflict, "Email already exists")
	}

	user := &dal.User{
		Name:     b.Name,
		Password: password.Generate(b.Password),
		Email:    b.Email,
	}

	// Create a user, if error return
	if err := dal.CreateUser(user); err.Error != nil {
		return velocity.NewError(velocity.StatusConflict, err.Error.Error())
	}

	// generate access token
	t := jwt.Generate(&jwt.TokenPayload{
		ID: user.ID,
	})

	return ctx.JSON(&types.AuthResponse{
		User: &types.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Auth: &types.AccessResponse{
			Token: t,
		},
	})
}
