package enum

type Permission string

const (
	Admin   Permission = "admin"
	Student Permission = "student"
)

func (p Permission) String() string {
	return string(p)
}
