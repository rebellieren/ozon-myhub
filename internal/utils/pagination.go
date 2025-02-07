package utils

func Paginate[T any](items []T, limit, offset int32) ([]T, bool) {

	if offset >= int32(len(items)) {
		return []T{}, false
	}

	end := offset + limit
	if end > int32(len(items)) {
		end = int32(len(items))
	}

	paginatedItems := items[offset:end]

	hasNextPage := end < int32(len(items))

	return paginatedItems, hasNextPage
}
