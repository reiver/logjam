package users

import (
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type userRepo struct {
}

const (
	usertbl = "userstbl"
)

func GetNewUsersRepository() IUserRepository {
	return &userRepo{}
}

func (u *userRepo) Create(input CreateUserDTO) (string, error) {
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return "", err
	}
	return dbsrv.Repository.Insert(usertbl, data)
}

func (u *userRepo) CreateOTP(input CreateOTPDTO) error {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) CompleteSignUp(input CompleteSignUpDTO) (CompleteSignUpResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) GetMe(id string) (UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userRepo) UpdateProfile(input UpdateProfileDTO) error {
	//TODO implement me
	panic("implement me")
}
