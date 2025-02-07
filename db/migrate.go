package db

import (
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&Post{}, &Comment{}, &User{}, &Reply{})
	if err != nil {
		log.Fatal("❌ Ошибка миграции:", err)
	}
	log.Println("✅ Таблицы успешно созданы в БД!")
}
