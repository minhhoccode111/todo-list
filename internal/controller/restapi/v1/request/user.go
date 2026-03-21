package request

type Register struct {
	Name     string `json:"name"     validate:"required,max=255,name"`
	Email    string `json:"email"    validate:"required,max=255,email"`
	Password string `json:"password" validate:"required,max=255,password"` //nolint:gosec // intended
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"` //nolint:gosec // intended
}
