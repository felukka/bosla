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
	Created_At  string    `json:"created_at"`
	Updated_At  string    `json:"updated_at,omitempty"`
	Service     []Service `json:"services"`
}

func ToProject(p database.Project) Project {
	return Project{
		Name:        p.Name,
		Status:      p.Status,
		Link:        p.Link.String,
		Description: p.Description.String,
		Created_At:  p.CreatedAt.Format(time.RFC1123),
		Updated_At:  p.UpdatedAt.Format(time.RFC1123),
	}
}

func ToProjects(ps []database.Project) []Project {
	result := make([]Project, 0, len(ps))
	for _, p := range ps {
		result = append(result, ToProject(p))
	}

	return result
}
