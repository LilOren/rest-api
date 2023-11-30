package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lil-oren/rest/internal/dto"
	"github.com/lil-oren/rest/internal/model"
)

type (
	CategoryRepository interface {
		FindByLevel(ctx context.Context, level int) ([]model.Category, error)
		FindByParentCategoryID(ctx context.Context, parentCategoryID int64) ([]model.Category, error)
		FindByPairCategoryID(ctx context.Context) ([]dto.HomePageCategoryPair, error)
	}
	categoryRepository struct {
		db *sqlx.DB
	}
)

// FindByPairCategoryID implements CategoryRepository.
func (r *categoryRepository) FindByPairCategoryID(ctx context.Context) ([]dto.HomePageCategoryPair, error) {
	qs := `
	SELECT
		DISTINCT ON 
		(c.name) 
		c.id AS first_level_id,
		c.name as first_level_name,
		c2.id AS second_level_id,
		c.image_url as image_url
	FROM
		categories c
	JOIN
		categories c2 
		ON
		c2.parent_category = c.id
	WHERE
			c.level = 1
	ORDER BY
		c.name
	`

	cpair := make([]dto.HomePageCategoryPair, 0)
	if err := r.db.SelectContext(ctx, &cpair, qs); err != nil {
		return nil, err
	}

	return cpair, nil
}

// FindByParentCategoryID implements CategoryRepository.
func (r *categoryRepository) FindByParentCategoryID(ctx context.Context, parentCategoryID int64) ([]model.Category, error) {
	res := make([]model.Category, 0)

	qs := `
	SELECT
		c.*
	FROM categories c
	WHERE
		parent_category = $1
	`

	if err := r.db.SelectContext(ctx, &res, qs, parentCategoryID); err != nil {
		return nil, err
	}

	return res, nil
}

// FindByLevel implements CategoryRepository.
func (r *categoryRepository) FindByLevel(ctx context.Context, level int) ([]model.Category, error) {
	res := make([]model.Category, 0)

	qs := `
	SELECT
		c.*
	FROM categories c
	WHERE
		c.level = $1
	`

	if err := r.db.SelectContext(ctx, &res, qs, level); err != nil {
		return nil, err
	}

	return res, nil
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}
