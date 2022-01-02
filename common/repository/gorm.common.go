package repository

import (
	"fmt"

	"github.com/Benyam-S/aseri/common"
	"github.com/Benyam-S/aseri/entity"
	"github.com/Benyam-S/aseri/tools"
	"github.com/jinzhu/gorm"
)

// CommonRepository is a type that defines a repository for common use
type CommonRepository struct {
	conn *gorm.DB
}

// NewCommonRepository is a function that returns a new common repository type
func NewCommonRepository(connection *gorm.DB) common.ICommonRepository {
	return &CommonRepository{conn: connection}
}

// IsUnique is a methods that checks if a given column value is unique in a certain table
func (repo *CommonRepository) IsUnique(columnName string, columnValue interface{}, tableName string) bool {
	return tools.IsUnique(columnName, columnValue, tableName, repo.conn)
}

// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
//    ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======

// CreateJobAttribute is a method that adds a new job attribute to the database
func (repo *CommonRepository) CreateJobAttribute(newAttribute *entity.JobAttribute, tableName string) error {

	var prefix string

	switch tableName {
	case "job_types":
		prefix = "TYPE"
	case "job_sectors":
		prefix = "SECTOR"
	case "education_levels":
		prefix = "LEVEL"
	}

	totalNumOfMembers := tools.CountMembers(tableName, repo.conn)
	newAttribute.ID = fmt.Sprintf(prefix+"-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)

	for !tools.IsUnique("id", newAttribute.ID, tableName, repo.conn) {
		totalNumOfMembers++
		newAttribute.ID = fmt.Sprintf(prefix+"-%s%d", tools.RandomStringGN(7), totalNumOfMembers+1)
	}

	err := repo.conn.Table(tableName).Create(newAttribute).Error
	if err != nil {
		return err
	}
	return nil
}

// FindJobAttribute is a method that finds a certain job attribute from the database using an identifier and table name.
// In FindJobAttribute() id and name are used as an key
func (repo *CommonRepository) FindJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {
	attribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).
		Where("id = ? || name = ?", identifier, identifier).
		First(attribute).Error

	if err != nil {
		return nil, err
	}
	return attribute, nil
}

// AllJobAttributes is a method that returns all the job attributes of a single job attribute table in the database
func (repo *CommonRepository) AllJobAttributes(tableName string) []*entity.JobAttribute {
	var attributes []*entity.JobAttribute
	err := repo.conn.Table(tableName).Find(&attributes).Error

	if err != nil {
		return []*entity.JobAttribute{}
	}
	return attributes
}

// UpdateJobAttribute is a method that updates a certain job attribute value in the database
func (repo *CommonRepository) UpdateJobAttribute(attribute *entity.JobAttribute, tableName string) error {

	prevAttribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).Where("id = ?", attribute.ID).First(prevAttribute).Error

	if err != nil {
		return err
	}

	err = repo.conn.Table(tableName).Save(attribute).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteJobAttribute is a method that deletes a certain job attribute from the database using an identifier.
// In DeleteJobAttribute() id and name are used as an key
func (repo *CommonRepository) DeleteJobAttribute(identifier, tableName string) (*entity.JobAttribute, error) {
	attribute := new(entity.JobAttribute)
	err := repo.conn.Table(tableName).Where("id = ? || name = ?", identifier, identifier).First(attribute).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Table(tableName).Delete(attribute)
	return attribute, nil
}

// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
//    ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======

// CreateReferral is a method that adds a new referral to the database
func (repo *CommonRepository) CreateReferral(newReferral *entity.Referral) error {

	newReferral.Code = tools.RandomStringGN(10)
	for !tools.IsUnique("code", newReferral.Code, "referrals", repo.conn) {
		newReferral.Code = tools.RandomStringGN(10)
	}

	err := repo.conn.Create(newReferral).Error
	if err != nil {
		return err
	}
	return nil
}

// FindRefferal is a method that finds a certain referral from the database using an identifier
// In FindReferral() only code is used as a key
func (repo *CommonRepository) FindReferral(code string) (*entity.Referral, error) {
	referral := new(entity.Referral)
	err := repo.conn.Model(referral).Where("code = ?", code).First(referral).Error

	if err != nil {
		return nil, err
	}
	return referral, nil
}

// AllReferrals is a method that returns all the referrals in the database
func (repo *CommonRepository) AllReferrals() []*entity.Referral {
	var referrals []*entity.Referral
	err := repo.conn.Model(&entity.Referral{}).Find(&referrals).Error

	if err != nil {
		return []*entity.Referral{}
	}
	return referrals
}

// UpdateReferral is a method that updates a certain referral value in the database
// UpdateReferral() uses code as a key
func (repo *CommonRepository) UpdateReferral(referral *entity.Referral) error {

	prevReferral := new(entity.Referral)
	err := repo.conn.Model(prevReferral).Where("code = ?", referral.Code).First(prevReferral).Error

	if err != nil {
		return err
	}

	/* --------------------------- can change layer if needed --------------------------- */

	referral.CreatedAt = prevReferral.CreatedAt

	/* -------------------------------------- end --------------------------------------- */

	err = repo.conn.Save(referral).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteReferral is a method that deletes a certain referral from the database using an identifier.
// In DeleteReferral() only code is used as a key
func (repo *CommonRepository) DeleteReferral(code string) (*entity.Referral, error) {
	referral := new(entity.Referral)
	err := repo.conn.Model(referral).Where("code = ?", code).First(referral).Error

	if err != nil {
		return nil, err
	}

	repo.conn.Delete(referral)
	return referral, nil
}

// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
//    ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======
// ======	======    ======   ======   ======   ======   ======   ======   ======   ======   ======   ======

// CreateReferralVisit is a method that adds a new referral visit to the database
func (repo *CommonRepository) CreateReferralVisit(newReferralVist *entity.ReferralVisit) error {

	err := repo.conn.Create(newReferralVist).Error
	if err != nil {
		return err
	}
	return nil
}

// FindReferralVisits is a method that finds referral visits from the database using an identifier
// In FindReferralVisits() only code is used as a key
func (repo *CommonRepository) FindReferralVisits(code string) []*entity.ReferralVisit {
	var referralVisits []*entity.ReferralVisit
	err := repo.conn.Model(&entity.ReferralVisit{}).Where("code = ?", code).Find(&referralVisits).Error

	if err != nil {
		return []*entity.ReferralVisit{}
	}
	return referralVisits
}

// AllReferralVisits is a method that returns all the referrals visits in the database
func (repo *CommonRepository) AllReferralVisits() []*entity.ReferralVisit {
	var referralVisits []*entity.ReferralVisit
	err := repo.conn.Model(&entity.ReferralVisit{}).Find(&referralVisits).Error

	if err != nil {
		return []*entity.ReferralVisit{}
	}
	return referralVisits
}

// CountReferralVists is a method that counts referral visits from the database using for certain code
func (repo *CommonRepository) CountReferralVists(code string) int64 {
	var count int64
	repo.conn.Model(&entity.ReferralVisit{}).Where("code = ?", code).Count(&count)
	return count
}

// CountAllReferralVists is a method that counts all referral visits from the database
func (repo *CommonRepository) CountAllReferralVists() int64 {
	var count int64
	repo.conn.Model(&entity.ReferralVisit{}).Count(&count)
	return count
}
