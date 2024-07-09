package req_dto_user

type Login struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type ViewDashboard struct {
	Partial string `query:"partial" validate:"omitempty,oneof=1 0"`
}

type EditProfile struct {
	Username    string `form:"username" validate:"omitempty"`
	Name        string `form:"name" validate:"omitempty"`
	OldPassword string `form:"old_password" validate:"required_with=NewPassword,omitempty,min=5"`
	NewPassword string `form:"new_password"  validate:"required_with=OldPassword,omitempty,nefield=OldPassword,min=5"`
}
