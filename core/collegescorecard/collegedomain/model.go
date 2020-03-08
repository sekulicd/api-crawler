package collegedomain

import "github.com/jinzhu/gorm"

type ApiResponse struct {
	Metadata Metadata `json:"metadata"`
	Results  []School `json:"results"`
}

type Metadata struct {
	Total      int `json:"total"`
	PageNumber int `json:"page"`
	PageSize   int `json:"per_page"`
}

type School struct {
	gorm.Model
	SchoolId int    `json:"id"`
	Name     string `json:"school.name"`
	City     string `json:"school.city"`
}
