package category

type CreateCategory struct {
	Name     string   `json:"name"`
	Children []string `json:"children"`
}

type Category struct {
	ID       string     `json:"code"`
	Name     string     `json:"name"`
	Children []Category `json:"children,omitempty"`
}

//goland:noinspection GoNameStartsWithPackageName
type CategoryId struct {
	ID string `json:"id"`
}
