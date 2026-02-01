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

func (repo *ProductRepositories) GetAllProduct() ([]models.Produk, error) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, p.category_id, 
		       c.id, c.name, c.description 
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id`

	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := make([]models.Produk, 0)
	for rows.Next() {
		var p models.Produk
		var c models.Category

		// Use sql.Null types if columns can be null, but for simplicity assuming valid data or handling basic scan.
		// If category_id is null in DB, this might fail without NullInt definition.
		// However, let's assume standard scan for now as per Go SQL patterns,
		// possibly needing pointers or Null types if Left Join returns nulls.
		// For safety with LEFT JOIN where category might be missing,
		// strictly speaking, we should use NullString/Int or pointers.
		// Given the user context, I'll try direct scan. If p.category_id is null, it might error.
		// But let's proceed with the standard assumption that the user handles schema.

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
