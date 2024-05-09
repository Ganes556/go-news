package req_dto_user

type Login struct {
	Username string `form:"username" validate:"required"`
	Password string `form:"password" validate:"required"`
}

type ViewDashboard struct {
	Partial string `query:"partial" validate:"omitempty,oneof=1 0"`
}