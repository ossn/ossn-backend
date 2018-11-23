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
