package util

const (
	D = "dealer"
	A = "admin"
	C = "customer"
)

func IsValidRole(role string) bool {
	switch role {
	case D, A, C:
		return true
	default:
		return false
	}
}
