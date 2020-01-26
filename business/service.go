package business

// Business contains information about one business
type Business struct {
	ID   int64
	Name string
}

// Service ...
type Service struct{}

// NewBusinessService ...
func NewBusinessService() Service {
	return Service{}
}

// GetByID gets the business by id
func (s *Service) GetByID(id string) (*Business, error) {

	return nil, nil
}

// CreateBusinessInput ...
type CreateBusinessInput struct {
	Name string
}

// CreateBusiness gets the business by id
func (s *Service) CreateBusiness(input CreateBusinessInput) (*Business, error) {
	return nil, nil
}
