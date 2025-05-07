package users

type IUserRepository interface {
	Create(input CreateUserDTO) (string, error)
	CreateOTP(input OTPDTO) error
	CompleteSignUp(input CompleteSignUpDTO) (*CompleteSignUpResponse, error)
	GetMe(id string) (*UserDTO, error)
	UpdateProfile(input UpdateProfileDTO) error
}
