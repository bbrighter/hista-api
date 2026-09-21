package shared

import "encore.dev/types/uuid"

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name     string    `gorm:"uniqueIndex"`
	Password string
	Settings UserSettings `gorm:"embedded;embedded_prefix:settings_"`
}

type Users []User

type UserProductInstance struct {
	ID                uint
	ProductId         string
	UserId            uuid.UUID `gorm:"index:idx_user_instance,unique;type:uuid"`
	ProductInstanceId uuid.UUID `gorm:"index:idx_user_instance,unique;type:uuid"`
}

type UserAppPermission struct {
	ID                    uint
	UserProductInstanceId uint `gorm:"index:idx_user_app_instance,unique"`
	UserProductInstance   UserProductInstance
	UserId                uuid.UUID `gorm:"index:idx_user_app_instance,unique;type:uuid"`
	User                  User
	App                   string `gorm:"index:idx_user_app_instance,unique"`
	PermissionLevel       string
}

type UserProductInstanceList []UserProductInstance

type UserAppPermissionList []UserAppPermission

func (uapl UserAppPermissionList) ToMap() map[uuid.UUID][]string {
	permMap := make(map[uuid.UUID][]string)
	for _, u := range uapl {
		piid := u.UserProductInstance.ProductInstanceId
		if appIds, ok := permMap[piid]; ok {
			appIds = append(appIds, u.App)
		} else {
			permMap[piid] = []string{u.App}
		}
	}
	return permMap
}

type UserSettings struct {
	LoadingMode string `gorm:"default:spinner;not null"`
	Language    string `gorm:"default:de-DE;not null"`
}
