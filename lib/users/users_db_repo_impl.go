package users

import (
	"errors"
	cErrors "github.com/reiver/logjam/lib/errors"
	"github.com/reiver/logjam/lib/marshal"
	"github.com/reiver/logjam/lib/tokens"
	dbsrv "github.com/reiver/logjam/srv/db"
	"net/http"
)

type userRepo struct {
}

const (
	userTbl = "userstbl"
	otpTbl  = "otps"

	nameKey     = "name"
	emailKey    = "email"
	otpKey      = "otp"
	usernameKey = "username"
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

func (u *userRepo) SignIn(input SignInDTO) (*CompleteSignUpResponse, error) {
	user, err := u.GetByEmail(input.Email)
	if err != nil || user == nil {
		return nil, cErrors.NewErrorWithMsg(http.StatusUnauthorized, "couldnt find user with this email")
	}
	rows, err := dbsrv.Repository.GetByFilter(otpTbl, map[string]any{
		emailKey: input.Email,
		otpKey:   input.Code,
	})
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return nil, cErrors.NewErrorWithMsg(http.StatusUnauthorized, "invalid email or code")
	}
	rec := rows[0]
	var otpDTO OTPDTO
	err = marshal.MapToObj(rec, &otpDTO)
	if err != nil {
		return nil, err
	}
	if len(otpDTO.Email) == 0 {
		return nil, cErrors.NewErrorWithMsg(http.StatusUnauthorized, "invalid email, try again or call support")
	}
	token, err := tokens.CreateToken(user.ID, nil)
	response := CompleteSignUpResponse{
		UserID: user.ID,
		Token:  token,
	}
	return &response, nil
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
		return nil, cErrors.NewErrorWithMsg(http.StatusUnauthorized, "invalid email or code")
	}
	rec := rows[0]
	var otpDTO OTPDTO
	err = marshal.MapToObj(rec, &otpDTO)
	if err != nil {
		return nil, err
	}
	if len(otpDTO.Email) == 0 {
		return nil, cErrors.NewErrorWithMsg(http.StatusUnauthorized, "invalid email, try again or call support")
	}
	userId, err := u.Create(CreateUserDTO{
		Email:    input.Email,
		Name:     input.Name,
		UserName: "",
		Bio:      "",
	})
	token, err := tokens.CreateToken(userId, nil)
	response := CompleteSignUpResponse{
		UserID: userId,
		Token:  token,
	}
	return &response, nil
}

func (u *userRepo) GetById(id string) (*UserDTO, error) {
	row, err := dbsrv.Repository.GetById(userTbl, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, cErrors.NewErrorWithMsg(http.StatusNotFound, "user not found")
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

func (u *userRepo) GetByEmail(email string) (*UserDTO, error) {
	rows, err := dbsrv.Repository.GetByFilter(userTbl, map[string]any{
		"email": email,
	})
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return nil, nil
	}
	var user UserDTO
	err = marshal.MapToObj(rows[0], &user)
	if err != nil {
		return nil, err
	}
	if len(user.ID) == 0 {
		return nil, errors.New("invalid user obj")
	}
	return &user, nil
}

func (u *userRepo) UpdateProfile(input UpdateProfileDTO) error {
	if input.Username != nil {
		existingUserNameRows, err := dbsrv.Repository.GetByFilter(userTbl, map[string]any{
			usernameKey: *input.Username,
		})
		if err != nil {
			return err
		}
		if existingUserNameRows != nil && len(existingUserNameRows) > 0 {
			return cErrors.NewErrorWithMsg(http.StatusConflict, "username is taken")
		}
	}
	data, err := marshal.ObjToMap(input)
	if err != nil {
		return err
	}
	return dbsrv.Repository.Update(userTbl, input.UserID, data)
}
