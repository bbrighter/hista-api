package entity

import "encore.dev/types/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"uniqueIndex"`
	Password string
}

type Users []User

type UserProductInstance struct {
	ID                uint
	ProductId         string
	UserId            uuid.UUID `gorm:"index:idx_user_instance,unique"`
	ProductInstanceId uuid.UUID `gorm:"index:idx_user_instance,unique"`
}

type UserAppPermission struct {
	ID                    uint
	UserProductInstanceId uint `gorm:"index:idx_user_app_instance,unique"`
	UserProductInstance   UserProductInstance
	UserId                uuid.UUID `gorm:"index:idx_user_app_instance,unique"`
	User                  User
	App                   string `gorm:"index:idx_user_app_instance,unique"`
	PermissionLevel       string
}

func (u User) ToResponse() UserResponse {
	return UserResponse{ID: u.ID, Name: u.Name}
}

func (us Users) ToResponse() UserListResponse {
	var users = []UserResponse{}
	for _, u := range us {
		users = append(users, u.ToResponse())
	}
	return UserListResponse{Users: users}
}
