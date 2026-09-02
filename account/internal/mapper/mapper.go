package mapper

import (
	"account/internal/model"
	accountpb "account/pkg/account/go"
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
		Balance:    user.Balance,
		IsDeleted:  user.IsDeleted,
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

func PbsToUsers(pbs []*accountpb.User) []model.User {
	users := make([]model.User, len(pbs))
	for i, pb := range pbs {
		users[i] = PbToUser(pb)
	}

	return users
}

func PbToUserCreate(accountpbUser *accountpb.CreateUser) model.CreateUser {
	if accountpbUser == nil {
		return model.CreateUser{}
	}

	return model.CreateUser{
		Login:      accountpbUser.Login,
		Email:      accountpbUser.Email,
		Password:   accountpbUser.Password,
		Phone:      accountpbUser.Phone,
		FirstName:  accountpbUser.FirstName,
		LastName:   accountpbUser.LastName,
		MiddleName: accountpbUser.MiddleName,
		Age:        accountpbUser.Age,
		Balance:    accountpbUser.Balance,
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
		Balance:    userpb.Balance,
	}
}
