package utils

func Pagination(offset, limit, total int) (prev, next int) {
	if offset-limit < 0 {
		prev = -999
	} else {
		prev = offset - limit
	}

	if offset+limit >= total {
		next = -999
	} else {
		next = offset + limit
	}
	return prev, next
}
