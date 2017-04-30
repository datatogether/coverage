package main

var services = []*Service{
	iaService,
	eotService,
	archiversService,
	archivers2Service,
}

// Service is anything that can provide information about archving
// status for a given url
type Service struct {
	Id          string
	Name        string
	Description string
	HomeUrl     string
}

func (s *Service) Read() error {
	for _, ser := range services {
		if s.Id == ser.Id {
			*s = *ser
			return nil
		}
	}
	return ErrNotFound
}
