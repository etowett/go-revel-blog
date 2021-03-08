package forms

import "github.com/revel/revel"

type (
	Post struct {
		UserID  int64  `json:"user_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Tag     string `json:"tag"`
	}
)

func (post *Post) Validate(v *revel.Validation) {
	ValidatePostTitle(v, post.Title)
	ValidatePostContent(v, post.Content)
	ValidatePostTag(v, post.Tag)
	return
}

func ValidatePostTitle(v *revel.Validation, title string) *revel.ValidationResult {
	result := v.Required(title).Message("Post title is required")
	if !result.Ok {
		return result
	}

	return result
}

func ValidatePostContent(v *revel.Validation, content string) *revel.ValidationResult {
	result := v.Required(content).Message("Post content is required")
	if !result.Ok {
		return result
	}

	return result
}

func ValidatePostTag(v *revel.Validation, tag string) *revel.ValidationResult {
	result := v.Required(tag).Message("Post tag is required")
	if !result.Ok {
		return result
	}

	return result
}
