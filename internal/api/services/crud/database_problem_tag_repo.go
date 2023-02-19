package crud

import (
	"github.com/omalbum/omalbum-back/internal/api/db"
	"github.com/omalbum/omalbum-back/internal/api/domain"
	"github.com/omalbum/omalbum-back/internal/api/messages"
)

type databaseProblemTagRepo struct {
	database *db.Database
}

func (d databaseProblemTagRepo) GetAllTags() []domain.ProblemTag {
	var tags []domain.ProblemTag
	d.database.DB.Find(&tags)
	return tags
}

func (d databaseProblemTagRepo) Create(problemTag *domain.ProblemTag) error {
	return d.database.DB.Create(problemTag).Error

}

func (d databaseProblemTagRepo) CreateByProblemIdAndTags(problemId uint, tags []string) error {
	for _, tag := range tags {
		var problemTag domain.ProblemTag
		problemTag.ProblemId = problemId
		problemTag.Tag = tag
		err := d.Create(&problemTag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d databaseProblemTagRepo) DeleteAllTagsByProblemId(problemId uint) error {
	if problemId == 0 { // esto es clave para no borrar la tabla entera accidentalmente
		return messages.New("problem_id_must_be_nonzero", "problem id must be nonzero")
	}
	return d.database.DB.Where(&domain.ProblemTag{ProblemId: problemId}).Delete(&domain.ProblemTag{ProblemId: problemId}).Error
}

func (d databaseProblemTagRepo) GetByProblemId(problemId uint) []domain.ProblemTag {
	if problemId == 0 {
		return nil //needed to avoid an empty condition in Where
	}
	var tags []domain.ProblemTag
	d.database.DB.Where("problem_id = ?", problemId).Find(&tags)
	return tags
}

func NewDatabaseProblemTagRepo(database *db.Database) domain.ProblemTagRepo {
	return &databaseProblemTagRepo{
		database: database,
	}
}
