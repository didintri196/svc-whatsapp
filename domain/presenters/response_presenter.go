package presenters

type ResponsePresenter struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type MetaResponsePresenter struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
}

const (
	//default limit for pagination
	defaultLimit = 10

	//max limit for pagination
	maxLimit = 50

	//default order by
	defaultOrderBy = "created_at"

	//default sort
	defaultSort = "asc"

	//default last page for pagination
	defaultLastPage = 0
)

func SetPaginationParameter(page, limit int, order, sort string) (int, int, int, string, string) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > maxLimit {
		limit = defaultLimit
	}
	if order == "" {
		order = defaultOrderBy
	}
	if sort == "" {
		sort = defaultSort
	}
	offset := (page - 1) * limit

	return offset, limit, page, order, sort
}

func SetPaginationResponse(page, limit, total int) (res MetaResponsePresenter) {
	var lastPage int

	if total > 0 {
		lastPage = total / limit

		if total%limit != 0 {
			lastPage = lastPage + 1
		}
	} else {
		lastPage = defaultLastPage
	}

	res = MetaResponsePresenter{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		LastPage:    lastPage,
	}

	return res
}
