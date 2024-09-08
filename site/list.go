package site

import "context"

type ListResponse struct {
	// Sites is the list of monitored sites.
	Sites []*Site `json:"sites"`
}

// List lists hte monitor websites.
//
//encore:api public method=GET path=/site
func (s *Service) List(ctx context.Context) (*ListResponse, error) {
	var sites []*Site
	err := s.db.WithContext(ctx).Find(&sites).Error
	if err != nil {
		return nil, err
	}
	return &ListResponse{Sites: sites}, nil
}
