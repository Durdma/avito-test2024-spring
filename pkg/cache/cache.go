package cache

import "avito-test2024-spring/internal/models"

type Cache interface {
	Set(banner models.Banner, tagId int, featureId int) error
	Get(tagId int, featureId int) (models.Banner, error)
}
