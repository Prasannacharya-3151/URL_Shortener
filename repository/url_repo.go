package repository

import (
	"url-shorter/config"
	"url-shorter/models"
)

func CreateURL(code, original string) (models.URL, error) {
	var url models.URL

	query := `
	INSERT INTO urls (code, original)
	VALUES ($1, $2)
	RETURNING id, code, original, created_at
	`

	row := config.DB.QueryRow(query, code, original)
	err := row.Scan(&url.ID, &url.Code, &url.Original, &url.CreatedAt)
	return url, err
}

func GetURLByCode(code string) (models.URL, error) {
	var url models.URL

	query := `SELECT id, code , original, created_at FROM urls WHERE code = $1`

	row := config.DB.QueryRow(query, code)
	err := row.Scan(&url.ID, &url.Code, &url.Original, &url.CreatedAt)
	return url, err
}