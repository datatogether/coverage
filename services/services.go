package services

import (
	"github.com/archivers-space/coverage/services/archivers2"
	"github.com/archivers-space/coverage/services/archivers_space"
	"github.com/archivers-space/coverage/services/eot"
	"github.com/archivers-space/coverage/services/ia"
	"github.com/archivers-space/coverage/services/svalbard"
	"github.com/archivers-space/coverage/tree"
)

// Services can describe themselves via an Info Method
type Service interface {
	Info() map[string]interface{}
}

// A Coverage Service is any service that can also provide coverage information
type CoverageService interface {
	Service
	AddUrls(tree *tree.Node) error
	AddCoverage(tree *tree.Node)
}

var Services = []CoverageService{
	archivers2.Service,
	archivers_space.Service,
	eot.Service,
	ia.Service,
	svalbard.Service,
}

// Service is anything that can provide information about archving
// status for a given url
// type Service struct {
//  Id          string
//  Name        string
//  Description string
//  HomeUrl     string
// }

// func (s *Service) Read() error {
//  for _, ser := range services {
//    if s.Id == ser.Id {
//      *s = *ser
//      return nil
//    }
//  }
//  return ErrNotFound
// }
