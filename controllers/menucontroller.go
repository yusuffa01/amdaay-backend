package controllers

import (
	"amdaaybackend/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func CreateMenu(c *gin.Context) {
	fmt.Println("\n=== 🕵️ MULAI PROSES UPLOAD TRANSIT ===")

	nama := c.PostForm("nama")
	deskripsi := c.PostForm("deskripsi")
	hargaStr := c.PostForm("harga")
	linkIG := c.PostForm("link_ig")
	linkShopee := c.PostForm("link_shopee")
	harga, _ := strconv.Atoi(hargaStr)

	form, _ := c.MultipartForm()
	files := form.File["gambar"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Gagal! Foto tidak sampai."})
		return
	}

	var namaFileUtama string
	var listFoto []models.FotoProduk

	_ = godotenv.Load("C:/Users/ADMIN/Documents/amdaay.scarf/backend/.env")
	urlCloudinary := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(urlCloudinary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal konek Cloudinary."})
		return
	}

	ctx := context.Background()

	// 3. Proses setiap foto dengan Jalur Transit
	for i, fileHeader := range files {
		fmt.Printf("📦 Memproses Foto Ke-%d: %s\n", i+1, fileHeader.Filename)

		namaTransit := fmt.Sprintf("uploads/temp_transit_%d_%s", i, fileHeader.Filename)
		errSave := c.SaveUploadedFile(fileHeader, namaTransit)
		if errSave != nil {
			fmt.Println("❌ GAGAL TRANSIT:", errSave)
			continue
		}

		fmt.Println("✈️ Menerbangkan foto ke Cloudinary...")
		resp, err := cld.Upload.Upload(ctx, namaTransit, uploader.UploadParams{
			Folder: "amdaay_produk",
		})

		os.Remove(namaTransit)

		if err != nil {
			fmt.Println("❌ SISTEM ERROR:", err)
			continue
		}

		// 🚨 ALAT BARU: Tangkap "Surat Penolakan" dari Cloudinary! 🚨
		if resp.Error.Message != "" {
			fmt.Println("🚫 DITOLAK CLOUDINARY! Alasan:", resp.Error.Message)
			// Hentikan proses dan kirim alasan penolakannya ke layar web Bundo!
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary menolak: " + resp.Error.Message})
			return
		}

		fmt.Println("✅ SUKSES MENDARAT! Link Resmi:", resp.SecureURL)

		if resp.SecureURL != "" {
			if i == 0 {
				namaFileUtama = resp.SecureURL
			}
			listFoto = append(listFoto, models.FotoProduk{Path: resp.SecureURL})
		}
	}

	menu := models.Menu{
		Nama:       nama,
		Deskripsi:  deskripsi,
		Harga:      harga,
		Gambar:     namaFileUtama,
		DaftarFoto: listFoto,
		LinkIG:     linkIG,
		LinkShopee: linkShopee,
	}

	models.DB.Create(&menu)
	fmt.Println("🎉 SEMUA PROSES SELESAI SEMPURNA!")
	c.JSON(http.StatusOK, gin.H{"pesan": "Produk & Foto berhasil mendarat di awan!", "data": menu})
}

func GetMenus(c *gin.Context) {
	var menus []models.Menu
	models.DB.Preload("DaftarFoto").Find(&menus)
	c.JSON(http.StatusOK, gin.H{"data": menus})
}

func GetMenuByID(c *gin.Context) {
	var menu models.Menu
	id := c.Param("id")

	if err := models.DB.Preload("DaftarFoto").First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu

	if err := models.DB.First(&menu, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	menu.Nama = c.PostForm("nama")
	menu.Deskripsi = c.PostForm("deskripsi")
	hargaStr := c.PostForm("harga")
	menu.LinkIG = c.PostForm("link_ig")
	menu.LinkShopee = c.PostForm("link_shopee")

	if hargaStr != "" {
		h, _ := strconv.Atoi(hargaStr)
		menu.Harga = h
	}

	if err := models.DB.Save(&menu).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Update Berhasil!"})
}

func DeleteMenu(c *gin.Context) {
	var menu models.Menu
	id := c.Param("id")

	if err := models.DB.Preload("DaftarFoto").First(&menu, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	models.DB.Delete(&menu)
	c.JSON(200, gin.H{"pesan": "Produk dan seluruh galerinya berhasil dihapus dari etalase!"})
}

func ToggleMenuStatus(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu

	if err := models.DB.First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	menu.Tersedia = !menu.Tersedia
	models.DB.Save(&menu)

	status := "Tersedia"
	if !menu.Tersedia {
		status = "Habis"
	}

	c.JSON(http.StatusOK, gin.H{"pesan": "Status diubah ke: " + status, "tersedia": menu.Tersedia})
}