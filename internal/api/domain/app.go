package domain

// Structs to store API responses

// Login
type LoginApp struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Response of the successful login
type LoginResponseApp struct {
	Token      string  `json:"token"`
	Expiration string  `json:"expiration"`
	User       UserApp `json:"user`
}

// Response of the successful refresh
type RefreshResponseApp struct {
	Token      string `json:"token"`
	Expiration string `json:"expiration"`
}

type RegistrationApp struct {
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Cellphone string `json:"cellphone"`
	Email     string `json:"email"`
}

type UserApp struct {
	UserID    uint   `json:"user_id"`
	UserName  string `json:"user_name"`
	Name      string `json:"name"`
	LastName  string `json:"last_name"`
	Cellphone string `json:"cellphone"`
	Email     string `json:"email"`
}

type PasswordWrappedApp struct {
	Password string `json:"password"`
}

type EmailWrappedApp struct {
	Email string `json:"email"`
}
