package models

type ProductPopularity struct {
	ProductID int    `json:"product_id"`
	Name      string `json:"name"`
	SoldCount int    `json:"sold_count"`
}

type SalesSummary struct {
	TotalRevenue      int                 `json:"total_revenue"`
	TotalTransaction  int                 `json:"total_transaction"`
	ProductPopularity []ProductPopularity `json:"product_popularity"`
}
