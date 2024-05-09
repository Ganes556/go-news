package req_dto_categories

type CreateCategory struct {
	Name string `form:"name" validate:"required"`
}

type UpdateCategory struct {
	ID   uint   `params:"id" validate:"required"`
	Name string `form:"name" validate:"omitempty"`
}

type ViewCategories struct {
	Partial string `query:"partial" validate:"omitempty,oneof=1 0"`
}

type DeleteNews struct {
	ID uint `params:"id" validate:"required"`
}