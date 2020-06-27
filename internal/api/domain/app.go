package domain

import "time"

// Structs to store API responses

// Login
type LoginApp struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

// Response of the successful login
type LoginResponseApp struct {
	Token      string  `json:"token"`
	Expiration string  `json:"expiration"`
	User       UserApp `json:"user`
}

// Response of the successful refresh
type RefreshResponseApp struct {
	Token      string `json:"token"`
	Expiration string `json:"expiration"`
}

type RegistrationApp struct {
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	BirthDate  string `json:"birth_date"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	IsStudent  bool   `json:"is_student"`
	SchoolYear uint   `json:"school_year"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	Department string `json:"department"`
	Location   string `json:"location"`
	School     string `json:"school"`
}

// everything but password
type UpdateProfileApp struct {
	UserName   string `json:"user_name"`
	Name       string `json:"name"`
	LastName   string `json:"last_name"`
	BirthDate  string `json:"birth_date"`
	Email      string `json:"email"`
	Gender     string `json:"gender"`
	IsStudent  bool   `json:"is_student"`
	SchoolYear uint   `json:"school_year"`
	Country    string `json:"country"`
	Province   string `json:"province"`
	Department string `json:"department"`
	Location   string `json:"location"`
	School     string `json:"school"`
}

type UserApp struct {
	UserID   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
}

type PasswordWrappedApp struct {
	Password string `json:"password"`
}

type EmailWrappedApp struct {
	Email string `json:"email"`
}

type ProblemSummaryApp struct {
	ProblemId uint `json:"problem_id"`
	Tags      []string
}

type ProblemsApp struct {
	Problems []ProblemSummaryApp `json:"problems"`
}

type ProblemStatsApp struct {
	ProblemId           uint      `json:"problem_id"`
	Attempts            uint      `json:"attempts"`
	Solved              bool      `json:"solved"`
	SolvedDuringContest bool      `json:"solved_during_contest"`
	DateSolved          time.Time `json:"date_solved"`
	Tags                []string
}
type AlbumApp struct {
	Album []ProblemStatsApp `json:"album"`
}
