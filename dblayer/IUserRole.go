package dblayer

type IUserRole struct {
	userRoleID     int32
	userRoleUserID int32
	userRoleRoleID int32
}

func NewUserRole(userRoleID int32, userRoleUserID int32, userRoleRoleID int32) *IUserRole {
	return &IUserRole{
		userRoleID:     userRoleID,
		userRoleUserID: userRoleUserID,
		userRoleRoleID: userRoleRoleID,
	}
}
