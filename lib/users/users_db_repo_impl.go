package users

import (
	"errors"
	"github.com/reiver/logjam/lib/marshal"
	"github.com/reiver/logjam/lib/tokens"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type userRepo struct {
}

const (
	userTbl = "userstbl"
	otpTbl  = "otps"

	nameKey  = "name"
	emailKey = "email"
	otpKey   = "otp"
)

func GetNewUsersRepository() IUserRepository {
	return &userRepo{}
}

func (u *userRepo) Create(input CreateUserDTO) (string, error) {
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return "", err
	}
	return dbsrv.Repository.Insert(userTbl, data)
}

func (u *userRepo) CreateOTP(input OTPDTO) error {
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return err
	}
	_, err = dbsrv.Repository.Insert(otpTbl, data)
	return err
}

func (u *userRepo) CompleteSignUp(input CompleteSignUpDTO) (*CompleteSignUpResponse, error) {
	rows, err := dbsrv.Repository.GetByFilter(otpTbl, map[string]any{
		emailKey: input.Email,
		otpKey:   input.Code,
	})
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return nil, errors.New("invalid email or code")
	}
	rec := rows[0]
	var otpDTO OTPDTO
	err = marshal.MapToObj(rec, &otpDTO)
	if err != nil {
		return nil, err
	}
	if len(otpDTO.Email) == 0 {
		return nil, errors.New("invalid email, try again or call support")
	}
	userId, err := dbsrv.Repository.Insert(userTbl, map[string]any{
		nameKey: input.Name,
	})
	token, err := tokens.CreateToken(userId, nil)
	response := CompleteSignUpResponse{
		UserID: userId,
		Token:  token,
	}
	return &response, nil
}

func (u *userRepo) GetMe(id string) (*UserDTO, error) {
	row, err := dbsrv.Repository.GetById(userTbl, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, errors.New("user not found")
	}
	var user UserDTO
	err = marshal.MapToObj(row, &user)
	if err != nil {
		return nil, err
	}
	if len(user.ID) == 0 {
		return nil, errors.New("invalid user obj")
	}
	return &user, nil
}

func (u *userRepo) UpdateProfile(input UpdateProfileDTO) error {
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return err
	}
	return dbsrv.Repository.Update(userTbl, input.UserID, data)
}
