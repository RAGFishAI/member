package postgres

import (
	"time"

	"gorm.io/gorm"
)

type TblMemberGroup struct {
	Id          int       `gorm:"primaryKey;auto_increment;type:serial"`
	Name        string    `gorm:"type:character varying"`
	Slug        string    `gorm:"type:character varying"`
	Description string    `gorm:"type:character varying"`
	IsActive    int       `gorm:"type:integer"`
	IsDeleted   int       `gorm:"type:integer"`
	CreatedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy   int       `gorm:"type:integer"`
	ModifiedOn  time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	DeletedOn   time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy   int       `gorm:"DEFAULT:NULL"`
	TenantId    int       `gorm:"type:integer;"`
}

type TblMember struct {
	Id               int       `gorm:"primaryKey;auto_increment;type:serial"`
	Uuid             string    `gorm:"type:character varying"`
	FirstName        string    `gorm:"type:character varying"`
	LastName         string    `gorm:"type:character varying"`
	Email            string    `gorm:"type:character varying"`
	MobileNo         string    `gorm:"type:character varying"`
	IsActive         int       `gorm:"type:integer"`
	ProfileImage     string    `gorm:"type:character varying"`
	ProfileImagePath string    `gorm:"type:character varying"`
	StorageType      string    `gorm:"type:character varying"`
	LastLogin        int       `gorm:"type:integer"`
	MemberGroupId    int       `gorm:"type:integer"`
	Password         string    `gorm:"type:character varying"`
	Username         string    `gorm:"DEFAULT:NULL"`
	Otp              int       `gorm:"DEFAULT:NULL"`
	OtpExpiry        time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	LoginTime        time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	IsDeleted        int       `gorm:"type:integer"`
	DeletedOn        time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	CreatedOn        time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	CreatedBy        int       `gorm:"type:integer"`
	ModifiedOn       time.Time `gorm:"type:timestamp without time zone;DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	TenantId         int       `gorm:"type:integer;"`
}

// MigrateTable creates this package related tables in your database
func MigrateTables(db *gorm.DB) {

	db.AutoMigrate(&TblMemberGroup{}, &TblMember{})

}
