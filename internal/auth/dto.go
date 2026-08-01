package auth

type CheckEmployeeStatus string

const (
	CheckEmployeeStatusNotFound CheckEmployeeStatus = "NOT_FOUND"
	CheckEmployeeStatusRegister CheckEmployeeStatus = "REGISTER"
	CheckEmployeeStatusLogin    CheckEmployeeStatus = "LOGIN"
)

type CheckEmployeeRequest struct {
	EmployeeNo string `json:"employee_no" binding:"required"`
}

type CheckEmployeeResponse struct {
	Status CheckEmployeeStatus `json:"status"`
}

type RegisterRequest struct {
	EmployeeNo      string `json:"employee_no" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type LoginRequest struct {
	EmployeeNo string `json:"employee_no" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`

	ExpiresIn int64 `json:"expires_in"`
}

type LoginResponse = TokenResponse

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
