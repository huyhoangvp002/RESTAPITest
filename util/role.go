package util

const (
	U = "user"
	A = "admin"
	C = "customer"
)

func isValidRole(role string) bool {
	switch role {
	case U, A, C:
		return true
	default:
		return false
	}
}
