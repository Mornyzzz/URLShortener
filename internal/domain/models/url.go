package models

type URL struct {
	ID    uint   `gorm:"primaryKey"`
	URL   string `gorm:"uniqueIndex:unique_url"`
	Alias string `gorm:"uniqueIndex:unique_alias"`
}
