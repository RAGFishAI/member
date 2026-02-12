package mysql

import (
	"time"

	"gorm.io/gorm"
)

type TblMemberGroup struct {
	Id          int       `gorm:"primaryKey;auto_increment"`
	Name        string    `gorm:"type:varchar(255)"`
	Slug        string    `gorm:"type:varchar(255)"`
	Description string    `gorm:"type:varchar(255)"`
	IsActive    int       `gorm:"type:int"`
	IsDeleted   int       `gorm:"type:int"`
	CreatedOn   time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	CreatedBy   int       `gorm:"type:int"`
	ModifiedOn  time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL;type:int"`
	DeletedBy   int       `gorm:"type:int"`
	DeletedOn   time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	TenantId    int       `gorm:"type:int;"`
}

type TblMember struct {
	Id               int       `gorm:"primaryKey;auto_increment"`
	Uuid             string    `gorm:"type:varchar(255)"`
	FirstName        string    `gorm:"type:varchar(255)"`
	LastName         string    `gorm:"type:varchar(255)"`
	Email            string    `gorm:"type:varchar(255)"`
	MobileNo         string    `gorm:"type:varchar(255)"`
	IsActive         int       `gorm:"type:int"`
	ProfileImage     string    `gorm:"type:varchar(255)"`
	ProfileImagePath string    `gorm:"type:varchar(255)"`
	StorageType      string    `gorm:"type:varchar(255)"`
	LastLogin        int       `gorm:"type:int"`
	MemberGroupId    int       `gorm:"type:int"`
	Password         string    `gorm:"type:varchar(255)"`
	Username         string    `gorm:"type:varchar(255)"`
	Otp              int       `gorm:"DEFAULT:NULL;type:int"`
	OtpExpiry        time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	LoginTime        time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	IsDeleted        int       `gorm:"type:int"`
	DeletedOn        time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL;type:int"`
	CreatedOn        time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	CreatedBy        int       `gorm:"type:int"`
	ModifiedOn       time.Time `gorm:"type:datetime;DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL;type:int"`
	TenantId         int       `gorm:"type:int;"`
}

// MigrateTable creates this package related tables in your database
func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(&TblMemberGroup{}, &TblMember{})

}
