package gvapi

type Ticket struct {
}

type AdminUser struct {
	Id          int    `json:"-" db:"admin_id"`
	Email       string `json:"email" db:"admin_email" binding:"required"`
	Password    string `json:"password" db:"admin_password" binding:"required"`
	Priveledges string `db:"privileges_level"`
}

type Client struct {
}
