package core

type StatusCode int

const (
	Ok       StatusCode = 200
	NotFound StatusCode = 404
)

func (s StatusCode) GetReason() string {
	switch s {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 301:
		return "Moved Permanently"
	case 308:
		return "Permanent Redirect"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "Bad Request"
	}
}
