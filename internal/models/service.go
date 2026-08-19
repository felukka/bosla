package models

import (
	"github.com/mahmoudk1000/bosla/internal/database"
)

type Service struct {
	Name        string `json:"name"`
	Status      string `json:"status,omitempty"`
	Repo_Url    string `json:"repo,omitempty"`
	Description string `json:"description,omitempty"`
	Created_At  string `json:"created_at"`
	Updated_At  string `json:"updated_at,omitempty"`
}

func ToService(s database.Service) Service {
	return Service{
		Name:        s.Name,
		Repo_Url:    s.RepoUrl.String,
		Description: s.Description.String,
		Created_At:  s.CreatedAt.Format("2006-01-02T15:04:05 -07:00:00"),
	}
}

func ToServices(apps []database.Service) []Service {
	results := make([]Service, 0, len(apps))
	for _, a := range apps {
		results = append(results, ToService(a))
	}

	return results
}
