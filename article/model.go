package article

import "time"

// ArticleModel represents the database table structure with GORM tags
type ArticleModel struct {
	ID        string    `gorm:"primaryKey;column:id"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

// TableName specifies the table name for GORM
func (ArticleModel) TableName() string {
	return "articles"
}

// ToEntity converts a database model to a domain entity
func (m ArticleModel) ToEntity() Article {
	return Article{
		ID:        m.ID,
		Title:     m.Title,
		Content:   m.Content,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// FromEntity converts a domain entity to a database model
func (m *ArticleModel) FromEntity(entity Article) {
	m.ID = entity.ID
	m.Title = entity.Title
	m.Content = entity.Content
	m.CreatedAt = entity.CreatedAt
	m.UpdatedAt = entity.UpdatedAt
}
