package ossn_backend

func getLimit(limit *int) int {
	l := 10
	if limit != nil {
		if *limit > 100 {
			l = 100
		}
		l = *limit
	}
	return l
}

func min(limit *int, first *int) *int {
	switch {
	case limit == nil:
		return first
	case first == nil:
		return limit
	case *limit > *first:
		return first
	case *limit < *first:
		return limit
	default:
		return first
	}
}
