package model

type AddressResponse struct {
	ID         int64  `json:"id"`
	Street     string `json:"street"`
	City       string `json:"city"`
	Province   string `json:"province"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	CreatedAt  int64  `json:"created_at"`
	UpdatedAt  int64  `json:"updated_at"`
}

type AddressResponseList []AddressResponse

type ListAddressRequest struct {
	UserID    int64 `json:"-" validate:"required"`
	ContactID int64 `json:"-" validate:"required"`
}

type CreateAddressRequest struct {
	UserID     int64  `json:"-"           validate:"required"`
	ContactID  int64  `json:"-"           validate:"required"`
	Street     string `json:"street"      validate:"max=255"`
	City       string `json:"city"        validate:"max=255"`
	Province   string `json:"province"    validate:"max=255"`
	PostalCode string `json:"postal_code" validate:"max=10"`
	Country    string `json:"country"     validate:"max=100"`
}

type UpdateAddressRequest struct {
	UserID     int64  `json:"-"           validate:"required"`
	ContactID  int64  `json:"-"           validate:"required"`
	ID         int64  `json:"-"           validate:"required"`
	Street     string `json:"street"      validate:"max=255"`
	City       string `json:"city"        validate:"max=255"`
	Province   string `json:"province"    validate:"max=255"`
	PostalCode string `json:"postal_code" validate:"max=10"`
	Country    string `json:"country"     validate:"max=100"`
}

type GetAddressRequest struct {
	UserID    int64 `json:"-" validate:"required"`
	ContactID int64 `json:"-" validate:"required"`
	ID        int64 `json:"-" validate:"required"`
}

type DeleteAddressRequest struct {
	UserID    int64 `json:"-" validate:"required"`
	ContactID int64 `json:"-" validate:"required"`
	ID        int64 `json:"-" validate:"required"`
}
