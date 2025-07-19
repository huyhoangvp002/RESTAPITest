package util

const (
	D = "buyer"
	A = "admin"
	C = "seller"
)

func IsValidRole(role string) bool {
	switch role {
	case D, A, C:
		return true
	default:
		return false
	}
}
