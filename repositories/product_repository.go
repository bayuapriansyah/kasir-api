package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepositories struct {
	db *sql.DB
}

func NewProductRepositories(db *sql.DB) *ProductRepositories {
	return &ProductRepositories{db: db}
}

func (repo *ProductRepositories) GetAllProduct(nameFilter string) ([]models.Produk, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, 
		       c.id, c.name, c.description 
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id`

	args := []interface{}{}

	if nameFilter != "" {
		query += " WHERE p.name ILIKE $1"
		args = append(args, "%"+nameFilter+"%")
	}
	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]models.Produk, 0)
	for rows.Next() {
		var p models.Produk
		var c models.Category

		var categoryID sql.NullInt64
		var categoryName sql.NullString
		var categoryDesc sql.NullString
		var categoryID_Join sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID,
			&categoryID_Join, &categoryName, &categoryDesc,
		)
		if err != nil {
			return nil, err
		}

		if categoryID.Valid {
			p.CategoryID = int(categoryID.Int64)
		}

		if categoryID_Join.Valid {
			c.ID = int(categoryID_Join.Int64)
			c.Name = categoryName.String
			c.Description = categoryDesc.String
			p.Category = &c
		}

		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepositories) CreateProduct(product *models.Produk) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

func (repo *ProductRepositories) GetProdukByID(id int) (*models.Produk, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id,
		       c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1`

	var p models.Produk
	var c models.Category

	var categoryID sql.NullInt64
	var categoryName sql.NullString
	var categoryDesc sql.NullString
	var categoryID_Join sql.NullInt64

	err := repo.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID,
		&categoryID_Join, &categoryName, &categoryDesc,
	)
	if err != nil {
		return nil, err
	}

	if categoryID.Valid {
		p.CategoryID = int(categoryID.Int64)
	}

	if categoryID_Join.Valid {
		c.ID = int(categoryID_Join.Int64)
		c.Name = categoryName.String
		c.Description = categoryDesc.String
		p.Category = &c
	}

	return &p, nil
}

func (repo *ProductRepositories) UpdateProduk(product *models.Produk) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepositories) DeleteProduk(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return nil
}
