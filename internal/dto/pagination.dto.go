package dto

type PaginationParams struct {
	Page   int
	Size   int
	Search string
	Status string
	SortBy string
	Order  string
}
