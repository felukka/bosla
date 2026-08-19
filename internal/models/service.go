package models

import (
	"time"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type Service struct {
	Name        string `json:"name"`
	Status      string `json:"status,omitempty"`
	RepoURL     string `json:"repo_url,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

func ToService(s database.Service) Service {
	return Service{
		Name:        s.Name,
		Status:      s.Status,
		RepoURL:     s.RepoUrl,
		Description: s.Description,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   s.UpdatedAt.Format(time.RFC3339),
	}
}

func ToServices(services []database.Service) []Service {
	results := make([]Service, 0, len(services))
	for _, service := range services {
		results = append(results, ToService(service))
	}

	return results
}
