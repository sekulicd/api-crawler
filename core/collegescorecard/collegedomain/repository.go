package collegedomain

type CollegeRepository interface {
	GetAll() ([]School, error)
	Create(school School) error
}
