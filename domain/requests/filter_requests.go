package requests

type FilterRequest struct {
	PerPage int    `query:"per_page, omitempty"`
	Page    int    `query:"page, omitempty"`
	Search  string `query:"search, omitempty"`
	Order   string `query:"order, omitempty"`
	Sort    string `query:"sort, omitempty"`
}
