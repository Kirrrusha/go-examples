package mapper

import (
	"gateway/internal/model"
	accountpb "gateway/pkg/account/go"
	authpb "gateway/pkg/auth/go"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbToUser(userpb *accountpb.User) model.User {
	if userpb == nil {
		return model.User{}
	}

	return model.User{
		ID:         userpb.Id,
		Login:      userpb.Login,
		Email:      userpb.Email,
		Phone:      userpb.Phone,
		FirstName:  userpb.FirstName,
		LastName:   userpb.LastName,
		MiddleName: userpb.MiddleName,
		Age:        userpb.Age,
		CreatedAt:  userpb.CreatedAt.AsTime(),
		UpdatedAt:  userpb.UpdatedAt.AsTime(),
	}
}

func UserToPb(user model.User) *accountpb.User {
	return &accountpb.User{
		Id:         user.ID,
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
		CreatedAt:  timestamppb.New(user.CreatedAt),
		UpdatedAt:  timestamppb.New(user.UpdatedAt),
	}
}

func UsersToPbs(users []model.User) []*accountpb.User {
	pbs := make([]*accountpb.User, len(users))
	for i, user := range users {
		pbs[i] = UserToPb(user)
	}

	return pbs
}

func CreateUserToUser(user model.CreateUser) model.User {
	return model.User{
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func PbToUserCreate(userpb *accountpb.CreateUser) model.CreateUser {
	if userpb == nil {
		return model.CreateUser{}
	}

	return model.CreateUser{
		Login:      userpb.Login,
		Email:      userpb.Email,
		Phone:      userpb.Phone,
		FirstName:  userpb.FirstName,
		LastName:   userpb.LastName,
		MiddleName: userpb.MiddleName,
		Age:        userpb.Age,
	}
}

func PbToUserUpdate(userpb *accountpb.UpdateUser) model.UpdateUser {
	if userpb == nil {
		return model.UpdateUser{}
	}

	return model.UpdateUser{
		Email:      userpb.Email,
		Phone:      userpb.Phone,
		FirstName:  userpb.FirstName,
		LastName:   userpb.LastName,
		MiddleName: userpb.MiddleName,
		Age:        userpb.Age,
	}
}

func UserCreateToPb(user model.User) *accountpb.CreateUser {
	return &accountpb.CreateUser{
		Login:      user.Login,
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func UserUpdateToPb(user model.UpdateUser) *accountpb.UpdateUser {
	return &accountpb.UpdateUser{
		Email:      user.Email,
		Phone:      user.Phone,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		MiddleName: user.MiddleName,
		Age:        user.Age,
	}
}

func PbToTokenPair(tokenPair *authpb.TokenPair) model.TokenPair {
	if tokenPair == nil {
		return model.TokenPair{}
	}

	return model.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}

func TokenPairToPb(tokenPair model.TokenPair) *authpb.TokenPair {
	return &authpb.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}
}
