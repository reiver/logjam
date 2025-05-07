package users

type IUserRepository interface {
	Create(input CreateUserDTO) (string, error)
	CreateOTP(input OTPDTO) error
	CompleteSignUp(input CompleteSignUpDTO) (*CompleteSignUpResponse, error)
	SignIn(input SignInDTO) (*CompleteSignUpResponse, error)
	GetById(id string) (*UserDTO, error)
	GetByEmail(email string) (*UserDTO, error)
	UpdateProfile(input UpdateProfileDTO) error
}
