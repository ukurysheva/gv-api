package gvapi

import "errors"

type Ticket struct {
}

type AuthAdminUser struct {
	Id          int    `json:"-" db:"admin_id"`
	Email       string `json:"email" db:"admin_email" binding:"required"`
	Password    string `json:"password" db:"admin_password" binding:"required"`
	Priveledges string `db:"privileges_level"`
}

type AuthUser struct {
	Id       int    `json:"-" db:"user_id"`
	Email    string `json:"email" db:"user_email" binding:"required"`
	Password string `json:"password" db:"user_password" binding:"required"`
}

type UpdateUserInput struct {
	Email      *string `json:"userEmail" db:"user_email"`
	Password   *string `json:"userPassword" db:"user_password"`
	FirstName  *string `json:"userFirstName" db:"user_first_name"`
	LastName   *string `json:"userLastName" db:"user_last_name"`
	MiddleName *string `json:"userMiddleName" db:"user_middle_name"`
	PhoneNum   *string `json:"userPhoneNum" db:"user_phone_number"`
	BirthDate  *string `json:"birthDate" db:"birth_date"`
	CountryId  *int    `json:"userCountryId" db:"user_country_id"`
}

func (i UpdateUserInput) Validate() error {
	if i.Email == nil && i.Password == nil && i.FirstName == nil && i.LastName == nil && i.MiddleName == nil && i.PhoneNum == nil && i.BirthDate == nil && i.CountryId == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
