package forms

import "github.com/revel/revel"

type (
	Category struct {
		UserID      int64  `json:"user_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func (cat *Category) Validate(v *revel.Validation) {
	ValidateCategoryName(v, cat.Name)
	ValidateCategoryDescription(v, cat.Description)
	return
}

func ValidateCategoryName(v *revel.Validation, name string) *revel.ValidationResult {
	result := v.Required(name).Message("Category name is required")
	if !result.Ok {
		return result
	}

	return result
}

func ValidateCategoryDescription(v *revel.Validation, description string) *revel.ValidationResult {
	result := v.Required(description).Message("Category description is required")
	if !result.Ok {
		return result
	}

	return result
}
