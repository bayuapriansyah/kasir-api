package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetSalesSummary() (*models.SalesSummary, error) {
	var summary models.SalesSummary

	//total revenue
	err := repo.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE DATE(created_at) = CURRENT_DATE").Scan(&summary.TotalRevenue)
	if err != nil {
		return nil, err
	}

	//total transaction
	err = repo.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE DATE(created_at) = CURRENT_DATE").Scan(&summary.TotalTransaksi)
	if err != nil {
		return nil, err
	}

	// produk laris
	err = repo.db.QueryRow(`
	SELECT p.name, COALESCE(SUM(td.quantity), 0) as sold_count
	FROM transaction_details td
	JOIN products p ON p.id = td.product_id
	JOIN transactions t ON t.id = td.transaction_id
	WHERE DATE(t.created_at) = CURRENT_DATE
	GROUP BY p.id, p.name
	ORDER BY sold_count DESC
	LIMIT 1
`).Scan(
		&summary.ProdukTerlaris.Nama,
		&summary.ProdukTerlaris.QtyTerjual,
	)

	if err != nil {
		return nil, err
	}
	return &summary, nil
}
