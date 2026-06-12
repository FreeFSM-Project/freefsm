package templates

type User struct {
	ID    int64
	Name  string
	Email string
	Role  string
}

type LoginPageData struct {
	Error string
}

type SetupPageData struct {
	Error string
}

type DashboardData struct{}
