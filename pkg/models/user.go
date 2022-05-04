package models

type User struct {
	BaseModel
	Name           string `json:"name"`
	Email          string `gorm:"uniqueindex" json:"email"`
	HashedPassword string `json:"-"`
	IsActive       bool   `json:"is_active"`
}

func HashPassword(password string) (string, error) {
	return "", nil
}

func CheckPassword(password, hash string) bool {
	return true
}
