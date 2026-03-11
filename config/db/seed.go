package db

import (
	"libraryOnline/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {

	var count int64
	db.Model(&models.User{}).Count(&count)
	if count > 0 {
		log.Println("Seed already applied, skipping...")
		return
	}

	pass, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	users := []models.User{
		{Name: "George", LastName: "Albarado",
			Email: "george.test@expple.com", Password: string(pass), Role: "PROFESOR"},
		{Name: "Andres Felipe", LastName: "Manzur",
			Email: "andres.test@expple.com", Password: string(pass), Role: "ESTUDIANTE"},
		{Name: "Jose", LastName: "Banquez",
			Email: "admin.test@expple.com", Password: string(pass), Role: "ADMIN"},
	}

	for _, user := range users {
		db.FirstOrCreate(&user, models.User{Email: user.Email})
	}

	editorials := []models.Editorial{
		{Name: "Editorial A"},
		{Name: "Editorial B"},
		{Name: "Edtorial C"},
		{Name: "Editorial D"},
	}
	for _, edit := range editorials {
		db.FirstOrCreate(&edit, models.Editorial{Name: edit.Name})
	}

	authors := []models.Author{
		{Name: "Jose", LastName: "Gaitan"},
		{Name: "Andres", LastName: "Barrios"},
		{Name: "Gabriel", LastName: "Marquez"},
		{Name: "Enrique", LastName: "Valdes"},
	}

	for i, author := range authors {
		db.FirstOrCreate(&authors[i], models.Author{Name: author.Name, LastName: author.LastName})
	}
	books := []models.Book{
		{Title: "Dios de la guerra", AvailableQuantity: 3, TotalQuantity: 5, EditorialId: 1, Authors: []models.Author{authors[0], authors[3]}},
		{Title: "Amanecer Dorado", AvailableQuantity: 2, TotalQuantity: 3, EditorialId: 3, Authors: []models.Author{authors[1], authors[2]}},
		{Title: "Los super heores", AvailableQuantity: 4, TotalQuantity: 5, EditorialId: 2, Authors: []models.Author{authors[0]}},
	}

	for _, book := range books {
		var existingBook = models.Book{}
		if err := db.Where("title = ?", book.Title).
			First(&existingBook).Error; err != nil {
			db.Create(&book)
		}
	}

	log.Println("Database seeded successfully!")
}

func SeedLoans(db *gorm.DB) {

	var count int64
	db.Model(&models.Loand{}).Count(&count)
	if count > 0 {
		log.Println("SeedLoans already applied, skipping...")
		return
	}

	var book1, book2 models.Book
	var user models.User

	db.First(&book1, 1) // Libro: Dios de la guerra
	db.First(&book2, 2) // Libro: Amanecer Dorado
	db.First(&user, 1)  // Usuario de prueba

	// fechas base
	now := time.Now()
	future := now.AddDate(0, 0, 7) // +7 días para devolver

	borrowed14 := now.AddDate(0, 0, -14)
	borrowed19 := now.AddDate(0, 0, -19)
	returned := now.AddDate(0, 0, -10)

	loans := []models.Loand{
		// Caso 1: Una reserva pendiente (Aparece en pestaña "Gestión de Reservas")
		{
			UserId:             user.ID,
			BookId:             book1.ID,
			Status:             "RESERVED",
			Quantity:           1,
			ReservationDate:    now,
			ExpectedReturnDate: &future,
		},
		// Caso 2: Un préstamo ya entregado (Aparece en pestaña "Control de Préstamos")
		{
			UserId:             user.ID,
			BookId:             book2.ID,
			Status:             "ACTIVE",
			Quantity:           1,
			ReservationDate:    now.AddDate(0, 0, -2), // Reservado hace 2 días
			BorrowedDate:       &now,                  // Entregado hoy
			ExpectedReturnDate: &future,
		},
		// Caso 3: Un préstamo vencido (Para probar las alertas en rojo del Admin)
		{
			UserId:          user.ID,
			BookId:          book1.ID,
			Status:          "ACTIVE",
			Quantity:        1,
			ReservationDate: now.AddDate(0, 0, -15),
			BorrowedDate:    &borrowed14,
			ExpectedReturnDate: func() *time.Time {
				t := now.AddDate(0, 0, -7)
				return &t
			}(), // Venció hace una semana
		},
		// Caso 4: Otra reserva de un libro diferente
		{
			UserId:             user.ID,
			BookId:             book2.ID,
			Status:             "RESERVED",
			Quantity:           1,
			ReservationDate:    now,
			ExpectedReturnDate: &future,
		},
		// Caso 5: Un préstamo ya devuelto (Historial)
		{
			UserId:           user.ID,
			BookId:           book1.ID,
			Status:           "RETURNED",
			Quantity:         1,
			ReservationDate:  now.AddDate(0, 0, -20),
			BorrowedDate:     &borrowed19,
			ActualReturnDate: &returned,
		},
	}

	for _, l := range loans {
		db.Create(&l)
	}
}
