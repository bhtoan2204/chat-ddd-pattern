package elasticsearch

import "io"

func readBody(body io.Reader) string {
	if body == nil {
		return ""
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return ""
	}
	return string(data)
}
