package models

type Menu struct {
    ID         uint         `gorm:"primaryKey"`
    Nama       string       `gorm:"type:varchar(100);not null"`
    Deskripsi  string       `gorm:"type:text"`
    Harga      int          `gorm:"not null"`
    Gambar     string       `gorm:"type:varchar(255)"` 
    DaftarFoto []FotoProduk `gorm:"foreignKey:MenuID;constraint:OnDelete:CASCADE;" json:"daftar_foto"`
    
    Tersedia   bool         `gorm:"default:true"`
    LinkIG     string       `gorm:"type:varchar(255)"`
    LinkShopee string       `gorm:"type:varchar(255)"`
}

type FotoProduk struct {
    ID     uint   `gorm:"primaryKey" json:"id"`
    MenuID uint   `json:"menu_id"`
    Path   string `json:"path"`
}