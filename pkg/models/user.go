package models

type User struct {
	BaseModel
	Name           string `json:"name"`
	Email          string `gorm:"uniqueindex" json:"email"`
	HashedPassword string `json:"-"`
	IsActive       bool   `json:"is_active"`
}

type UserList struct {
	Items []User `json:"items"`
}

type UserInput struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required,min=3,max=50"`
}

type UserUpdate struct {
	Name string `form:"name" json:"name" binding:"required"`
}
