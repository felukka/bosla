package models

import (
	"time"

	"github.com/mahmoudk1000/bosla/internal/database"
)

type Project struct {
	Name        string    `json:"name"`
	Status      string    `json:"status,omitempty"`
	Link        string    `json:"link,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at,omitempty"`
	Services    []Service `json:"services,omitempty"`
}

func ToProject(p database.Project) Project {
	return Project{
		Name:        p.Name,
		Status:      p.Status,
		Link:        p.Link,
		Description: p.Description,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}
}

func ToProjects(ps []database.Project) []Project {
	result := make([]Project, 0, len(ps))
	for _, p := range ps {
		result = append(result, ToProject(p))
	}

	return result
}
