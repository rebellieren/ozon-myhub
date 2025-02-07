package db

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateTestUsersIfNotExist(db *gorm.DB) {
	testUsers := []User{
		{ID: uuid.New(), Nickname: "test_user_1"},
		{ID: uuid.New(), Nickname: "test_user_2"},
		{ID: uuid.New(), Nickname: "test_user_3"},
	}

	var users []User

	if err := db.Find(&users).Error; err != nil {
		log.Println("Ошибка при проверке пользователей:", err)
		return
	}

	if len(users) == 0 {
		for _, testUser := range testUsers {
			var existingUser User
			if err := db.First(&existingUser, "nickname = ?", testUser.Nickname).Error; err != nil {
				if err := db.Create(&testUser).Error; err != nil {
					fmt.Printf("Не удалось создать пользователя %s: %v\n", testUser.Nickname, err)
				} else {
					fmt.Printf("Тестовый пользователь %s успешно создан!\n", testUser.Nickname)
				}
			}
		}
	}

	if err := db.Find(&users).Error; err != nil {
		log.Println("Ошибка при получении пользователей:", err)
		return
	}

	log.Println(" ✅ Список всех пользователей:")
	for _, user := range users {
		fmt.Printf("Пользователь: %s (ID: %s)\n", user.Nickname, user.ID)
	}
}
