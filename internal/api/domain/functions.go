package domain

import "time"

func (problem *Problem) IsViewable() bool {
	now := time.Now()
	if problem.IsDraft || now.Before(problem.DateContestStart) {
		return false
	}
	return true
}

func (problem *Problem) IsCurrentProblem() bool {
	now := time.Now()
	if problem.IsDraft || now.Before(problem.DateContestStart) || problem.DateContestEnd.Before(now) {
		return false
	}
	return true
}

func (problemAdminApp *ProblemAdminApp) IsCurrentProblemAdminApp() bool {
	now := time.Now()
	if problemAdminApp.IsDraft || now.Before(problemAdminApp.ReleaseDate) || problemAdminApp.Deadline.Before(now) {
		return false
	}
	return true
}

func (problem *Problem) IsContestFinished() bool {
	now := time.Now()
	if problem.IsDraft || now.Before(problem.DateContestEnd) {
		return false
	}
	return true
}
