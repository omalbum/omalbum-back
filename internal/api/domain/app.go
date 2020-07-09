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
	UserID           uint      `json:"user_id"`
	UserName         string    `json:"user_name"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	LastName         string    `json:"last_name"`
	BirthDate        time.Time `json:"birth_date"`
	Gender           string    `json:"gender"`
	IsStudent        bool      `json:"is_student"`
	SchoolYear       uint      `json:"school_year"`
	Country          string    `json:"country"`
	Province         string    `json:"province"`
	Department       string    `json:"department"`
	Location         string    `json:"location"`
	School           string    `json:"school"`
	RegistrationDate time.Time `json:"registration_date"`
}

type PasswordWrappedApp struct {
	Password string `json:"password"`
}

type EmailWrappedApp struct {
	Email string `json:"email"`
}

type ProblemSummaryApp struct {
	ProblemId uint   `json:"problem_id"`
	Series    string `json:"series"`
	Tags      []string
}

type ProblemsApp struct {
	Problems []ProblemSummaryApp `json:"problems"`
}

type ProblemStatsApp struct {
	ProblemId uint `json:"problem_id"`

	Attempts            uint      `json:"attempts"`
	Solved              bool      `json:"solved"`
	SolvedDuringContest bool      `json:"solved_during_contest"`
	DateSolved          time.Time `json:"date_solved"`
	Series              string    `json:"series"`
	Tags                []string  `json:"tags"`
}
type AlbumApp struct {
	Album []ProblemStatsApp `json:"album"`
}

type NextProblemsApp struct {
	NextProblems []ProblemNextApp `json:"next_problems"`
}

type CurrentProblemsApp struct {
	CurrentProblems []ProblemApp `json:"current_problems"`
}

type AllProblemsApp struct {
	Problems []ProblemApp `json:"all_problems"`
}

type ProblemAdminApp struct {
	ProblemId        uint      `json:"problem_id"`
	Statement        string    `json:"statement"`
	Answer           int       `json:"answer"`
	OmaforosPostId   uint      `json:"omaforos_post_id"`
	Annotations      string    `json:"annotations"`
	Hint             string    `json:"hint"`
	Series           string    `json:"series"`
	Tags             []string  `json:"tags"`
	OfficialSolution string    `json:"official_solution"`
	ReleaseDate      time.Time `json:"release_date"`
	Deadline         time.Time `json:"deadline"`
	IsDraft          bool      `json:"is_draft"`
}

type ProblemApp struct {
	ProblemId      uint      `json:"problem_id"`
	Statement      string    `json:"statement"`
	OmaforosPostId uint      `json:"omaforos_post_id"`
	Series         string    `json:"series"`
	Tags           []string  `json:"tags"`
	ReleaseDate    time.Time `json:"release_date"`
	Deadline       time.Time `json:"deadline"`
}

type ProblemNextApp struct {
	ProblemId   uint      `json:"problem_id"`
	ReleaseDate time.Time `json:"release_date"`
	Deadline    time.Time `json:"deadline"`
	Series      string    `json:"series"`
}

type ProblemAttemptApp struct {
	ProblemId uint `json:"problem_id"`
	Answer    int  `json:"answer"`
}
type AttemptResultApp struct {
	AttemptId uint      `json:"attempt_id"`
	Result    string    `json:"result"`
	Deadline  time.Time `json:"deadline"`
}
