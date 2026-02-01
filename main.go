package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"

	// "strconv"
	"strings"

	"github.com/spf13/viper"
)

// type Kategori struct {
// 	ID           int    `json:"id"`
// 	NamaKategori string `json:"namaKategori"`
// 	Deskripsi    string `json:"deskripsi"`
// }

// var produk = []models.Produk{
// 	{ID: 1, Name: "Indomie Goreng", Price: 3500, Stock: 10},
// 	{ID: 2, Name: "Vit 1000ml", Price: 3000, Stock: 20},
// 	{ID: 3, Name: "Silver Queen", Price: 15000, Stock: 15},
// }

// var kategori = []Kategori{
// 	{ID: 1, NamaKategori: "Makanan Ringan", Deskripsi: "Maknan yang cocok untuk cemilan"},
// 	{ID: 2, NamaKategori: "Minuman", Deskripsi: "Minuman segar untuk menghilangkan dahaga"},
// }

//categories

// func getKategoryByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	for _, k := range kategori {
// 		if k.ID == id {
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(k)
// 			return
// 		}
// 	}
// 	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
// }

// func updateKategori(w http.ResponseWriter, r *http.Request) {
// 	//get id
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid Kategori ID", http.StatusBadRequest)
// 		return
// 	}
// 	var updateKategori Kategori
// 	err = json.NewDecoder(r.Body).Decode(&updateKategori)
// 	if err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	for i, k := range kategori {
// 		if k.ID == id {
// 			updateKategori.ID = id
// 			kategori[i] = updateKategori
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(updateKategori)
// 			return
// 		}
// 	}

// 	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
// }

// func deleteKategori(w http.ResponseWriter, r *http.Request) {
// 	// get id
// 	idStr := strings.TrimPrefix(r.URL.Path, "/api/kategori/")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	for i, k := range kategori {
// 		if k.ID == id {
// 			kategori = append(kategori[:i], kategori[i+1:]...)
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "Kategori berhasil dihapus",
// 			})
// 			return
// 		}
// 	}
// 	http.Error(w, "Kategori tidak ditemukan", http.StatusNotFound)
// }

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	//setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepositories(db)
	categoryRepo := repositories.NewCategoryRepository(db)

	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	http.HandleFunc("/api/produk", productHandler.HandleProduk)
	http.HandleFunc("/api/produk/", productHandler.HandleProdukByID)

	http.HandleFunc("/api/category", categoryHandler.HandleCategory)
	http.HandleFunc("/api/category/", categoryHandler.HandleCategoryByID)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running di port : " + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)

	if err != nil {
		fmt.Println("Gagal runing server")
	}
}
