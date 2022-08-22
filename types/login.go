package types

// Input login args
type Input struct {
	Phone string `json:"phone" validate:"required"`
	Code  string `json:"code" validate:"required"`
}

type PhoneNumber struct {
	Phone string `json:"phone" validate:"required"`
}
