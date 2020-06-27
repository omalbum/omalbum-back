package crud

import (
	"github.com/miguelsotocarlos/teleoma/internal/api/db"
	"github.com/miguelsotocarlos/teleoma/internal/api/domain"
	"github.com/miguelsotocarlos/teleoma/internal/api/utils/crypto"
	"time"
)

func DropTables(db *db.Database) {
	// Creates all the tables
	db.DB.DropTableIfExists( // because of the foreing keys we must drop tables in this particular order
		&domain.UserAction{},
		&domain.User{},
	)
}

func CreateTables(db *db.Database) {
	// Creates all the tables
	db.DB.CreateTable(
		&domain.UserAction{},
		&domain.User{},
	)
}

func CreateForeignKeys(db *db.Database) {

	db.DB.Model(&domain.Problem{}).AddForeignKey("poser_id", "users(id)", "CASCADE", "CASCADE")

	db.DB.Model(&domain.UserProblemAttempt{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	db.DB.Model(&domain.UserProblemAttempt{}).AddForeignKey("problem_id", "problems(id)", "CASCADE", "CASCADE")

	db.DB.Model(&domain.UserAction{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

	db.DB.Model(&domain.ProblemTag{}).AddForeignKey("problem_id", "problems(id)", "CASCADE", "CASCADE")
}

func CreateAdminUser(db *db.Database) {
	//TODO
	userRepo := NewDatabaseUserRepo(db)

	userAdmin := domain.User{
		UserName:         "admin",
		HashedPassword:   crypto.HashAndSalt("admin123456"),
		Email:            "admin@gmail.com",
		RegistrationDate: time.Now(),
		LastActiveDate:   time.Now(),
		IsAdmin:          true,
	}
	_ = userRepo.Create(&userAdmin)
}

func CreateAnonymousUser(db *db.Database) {

	anonymousPass, _ := crypto.GenerateRandomString(20)

	anonymousUser := domain.User{
		UserName:         "anonymous",
		HashedPassword:   crypto.HashAndSalt(anonymousPass), // to prevent someone from logging in :)
		RegistrationDate: time.Now(),
		LastActiveDate:   time.Now(),
		Name:             "",
		LastName:         "",
		Email:            "anonymous@teleoma.com.ar",
		BirthDate:        time.Now(),
	}
	db.DB.Create(&anonymousUser)

	if anonymousUser.ID != 1 {
		panic("Error, the anonymous user should have ID=1!")
	}
}

func CreateSampleData(database *db.Database) {
	// data de prueba para testear

}
