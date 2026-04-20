package entity

func normalizePair(a, b string) (string, string) {
	if a < b {
		return a, b
	}
	return b, a
}
