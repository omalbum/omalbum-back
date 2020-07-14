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
		&domain.UserProblemAttempt{},
		&domain.ProblemTag{},
		&domain.Problem{},
		&domain.User{},
	)
}

func CreateTables(db *db.Database) {
	// Creates all the tables
	db.DB.CreateTable(
		&domain.User{},
		&domain.Problem{},
		&domain.UserProblemAttempt{},
		&domain.ProblemTag{},
		&domain.UserAction{},
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
		BirthDate:        time.Now(),
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

func RefreshViews(db *db.Database) {
	db.DB.Exec("DROP VIEW IF EXISTS expanded_user_problem_attempts")
	createExpandedUserProblemAttempts := `
	CREATE VIEW expanded_user_problem_attempts AS
		SELECT
			problems.answer as answer,
			problems.date_contest_end as date_contest_end,
			problems.date_contest_start as date_contest_start,
			user_problem_attempts.id as attempt_id,
			user_problem_attempts.user_id as user_id,
			user_problem_attempts.date as attempt_date,
			user_problem_attempts.problem_id as problem_id,
			user_problem_attempts.user_answer as user_answer,
			problems.answer = user_problem_attempts.user_answer as is_correct,
			(user_problem_attempts.date <  problems.date_contest_end) and (user_problem_attempts.date >  problems.date_contest_start)   as during_contest
		FROM
			problems
		INNER JOIN
			user_problem_attempts
		ON
			problems.id = user_problem_attempts.problem_id `
	db.DB.Exec(createExpandedUserProblemAttempts)
}
