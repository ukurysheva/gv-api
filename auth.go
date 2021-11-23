package gvapi

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
