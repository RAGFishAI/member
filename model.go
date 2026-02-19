package member

import (
	"time"

	"gorm.io/gorm"
)

type Filter struct {
	Keyword       string
	Category      string
	Status        string
	FromDate      string
	ToDate        string
	FirstName     string
	MemberProfile bool
}

type MemberGroupListReq struct {
	Limit            int
	Offset           int
	Keyword          string
	Category         string
	Status           string
	FromDate         string
	ToDate           string
	FirstName        string
	ActiveGroupsOnly bool
}

type MemberGroupCreation struct {
	Name        string
	Description string
	CreatedBy   int
	IsActive    int
}

type MemberCreationUpdation struct {
	FirstName        string
	LastName         string
	Email            string
	MobileNo         string
	CreatedBy        int
	ModifiedBy       int
	IsActive         int
	ProfileImage     string
	ProfileImagePath string
	Username         string
	Password         string
	GroupId          int
	StorageType      string
	TenantId         string
}

type MemberGroupCreationUpdation struct {
	Name        string
	Description string
	CreatedBy   int
	ModifiedBy  int
	IsActive    int
}

type TblMemberGroup struct {
	Id          int
	Name        string
	Slug        string
	Description string
	IsActive    int
	IsDeleted   int
	CreatedOn   time.Time
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"default:null"`
	ModifiedBy  int       `gorm:"default:null"`
	DeletedOn   time.Time `gorm:"default:null"`
	DeletedBy   int       `gorm:"default:null"`
	TenantId    string
}

type TblMember struct {
	Id               int
	Uuid             string
	FirstName        string
	LastName         string
	Email            string
	MobileNo         string
	IsActive         int
	ProfileImage     string
	ProfileImagePath string
	StorageType      string
	LastLogin        int
	MemberGroupId    int
	Password         string
	Username         string
	Otp              int
	OtpExpiry        time.Time
	LoginTime        time.Time
	IsDeleted        int
	DeletedOn        time.Time `gorm:"default:null"`
	DeletedBy        int       `gorm:"default:null"`
	CreatedOn        time.Time
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"default:null"`
	ModifiedBy       int       `gorm:"default:null"`
	TenantId         string
}

type Tblmember struct {
	Id               int `gorm:"primaryKey;auto_increment;"`
	Uuid             string
	FirstName        string
	LastName         string
	Email            string
	MobileNo         string
	IsActive         int
	ProfileImage     string
	ProfileImagePath string
	StorageType      string
	LastLogin        int
	IsDeleted        int
	DeletedOn        time.Time `gorm:"DEFAULT:NULL"`
	DeletedBy        int       `gorm:"DEFAULT:NULL"`
	CreatedOn        time.Time `gorm:"DEFAULT:NULL"`
	CreatedDate      string    `gorm:"-"`
	CreatedBy        int
	ModifiedOn       time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy       int       `gorm:"DEFAULT:NULL"`
	MemberGroupId    int
	GroupName        string `gorm:"-:migration;<-:false"`
	Password         string
	DateString       string    `gorm:"-"`
	Username         string    `gorm:"DEFAULT:NULL"`
	Otp              int       `gorm:"DEFAULT:NULL"`
	OtpExpiry        time.Time `gorm:"DEFAULT:NULL"`
	ModifiedDate     string    `gorm:"-"`
	NameString       string    `gorm:"-"`
	LoginTime        time.Time `gorm:"DEFAULT:NULL"`
	Token            string    `gorm:"-"`
	Claimstatus      int       `gorm:"-"`

	TenantId string
}
type TblMemberSetting struct {
	Id                int
	AllowRegistration int
	MemberLogin       string // otp/password
	ModifiedBy        int
	ModifiedOn        time.Time
	NotificationUsers string //notification team users id
	TenantId          string
}
type MemberSettings struct {
	AllowRegistration int
	MemberLogin       string // otp/password
	ModifiedBy        int
	NotificationUsers string //notification team users id
}

// soft delete check
func IsDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("is_deleted = 0")
}

type MemberModel struct {
	Userid     int
	DataAccess int
}

var Membermodel MemberModel

// Member Group List
func (membermodel MemberModel) MemberGroupList(listre MemberGroupListReq, DB *gorm.DB, tenantid string) (membergroup []Tblmembergroup, TotalMemberGroup int64, err error) {

	query := DB.Table("tbl_member_groups").Where("tenant_id=?", tenantid).Scopes(IsDeleted).Order("id desc")

	if membermodel.DataAccess == 1 {

		query = query.Where("tbl_member_groups.created_by =?", membermodel.Userid)

	}

	if listre.Keyword != "" {

		query = query.Where("LOWER(TRIM(name)) LIKE LOWER(TRIM(?))", "%"+listre.Keyword+"%")

	}

	if listre.ActiveGroupsOnly {

		query = query.Where("is_active=1")

	}

	if listre.Limit != 0 {

		query.Limit(listre.Limit).Offset(listre.Offset).Find(&membergroup)

		return membergroup, 0, err

	}

	query.Find(&membergroup).Count(&TotalMemberGroup)

	return membergroup, TotalMemberGroup, err

}

// Member Group Insert
func (membermodel MemberModel) MemberGroupCreate(membergroup *TblMemberGroup, DB *gorm.DB) error {

	if err := DB.Table("tbl_member_groups").Create(&membergroup).Error; err != nil {

		return err
	}

	return nil
}

// Member list
func (membermodel MemberModel) MembersList(limit int, offset int, filter Filter, flag bool, DB *gorm.DB, tenantid string) (member []Tblmember, Total_Member int64, err error) {

	query := DB.Debug().Table("tbl_members").Select("tbl_members.id,tbl_members.uuid,tbl_members.member_group_id,tbl_members.first_name,tbl_members.last_name,tbl_members.email,tbl_members.mobile_no,tbl_members.profile_image,tbl_members.profile_image_path,tbl_members.created_on,tbl_members.created_by,tbl_members.modified_on,tbl_members.modified_by,tbl_members.is_active,tbl_members.is_deleted,tbl_members.deleted_on,tbl_members.deleted_by,tbl_member_groups.name as group_name,tbl_members.storage_type,tbl_members.username").Joins("left join tbl_member_groups on tbl_members.member_group_id = tbl_member_groups.id").Where("tbl_members.is_deleted=? and tbl_members.tenant_id=?", 0, tenantid).Order("id desc")

	if membermodel.DataAccess == 1 {

		query = query.Where("tbl_members.created_by =?", membermodel.Userid)

	}

	if filter.Keyword != "" {

		query = query.Where("LOWER(TRIM(tbl_members.first_name)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_members.last_name)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_member_profiles.company_name)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_member_groups.name)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_members.email)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_members.mobile_no)) LIKE LOWER(TRIM(?))  OR LOWER(TRIM(tbl_member_profiles.profile_slug)) LIKE LOWER(TRIM(?)) OR LOWER(TRIM(tbl_member_profiles.company_location)) LIKE LOWER(TRIM(?)) AND tbl_members.is_deleted=0 AND tbl_member_groups.is_deleted=0", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%")

	}

	if filter.FirstName != "" {

		query = query.Debug().Where("LOWER(TRIM(tbl_members.first_name)) LIKE LOWER(TRIM(?))"+" OR LOWER(TRIM(tbl_members.last_name)) LIKE LOWER(TRIM(?))", "%"+filter.FirstName+"%", "%"+filter.FirstName+"%")

	}

	if flag {

		query.Find(&member)

		return member, 0, err

	}

	if limit != 0 && !flag {

		query.Offset(offset).Limit(limit).Order("id desc").Find(&member)

		return member, 0, err

	}
	query.Find(&member).Count(&Total_Member)

	return member, Total_Member, nil

}

// Member Insert
func (membermodel MemberModel) MemberCreate(member *Tblmember, DB *gorm.DB) error {

	if err := DB.Table("tbl_members").Create(&member).Error; err != nil {

		return err
	}

	return nil
}

// Update Member
func (membermodel MemberModel) UpdateMember(member *Tblmember, DB *gorm.DB, tenantid string) error {

	query := DB.Table("tbl_members").Where("id=? and tenant_id=?", member.Id, tenantid)

	if member.Password == "" || member.ProfileImage == "" {

		if member.Password == "" && member.ProfileImage == "" {

			query.Omit("password , profile_image , profile_image_path").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy})

		} else if member.ProfileImage == "" {

			query.Omit(" profile_image , profile_image_path").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "password": member.Password})

		} else if member.Password == "" {

			query.Omit("password").UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath})

		}
	} else {

		query.UpdateColumns(map[string]interface{}{"first_name": member.FirstName, "last_name": member.LastName, "member_group_id": member.MemberGroupId, "email": member.Email, "username": member.Username, "mobile_no": member.MobileNo, "is_active": member.IsActive, "modified_on": member.ModifiedOn, "modified_by": member.ModifiedBy, "profile_image": member.ProfileImage, "profile_image_path": member.ProfileImagePath, "password": member.Password, "storage_type": member.StorageType})
	}
	return nil
}

// Delete Member
func (membermodel MemberModel) DeleteMember(member *Tblmember, id int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Where("id=? and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

// Check Email is already exists
func (membermodel MemberModel) CheckEmailInMember(member *TblMember, email string, userid int, DB *gorm.DB, tenantid string) error {

	if userid == 0 {
		if err := DB.Table("tbl_members").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and is_deleted=0 and tenant_id=?", email, tenantid).First(&member).Error; err != nil {

			return err
		}
	} else {
		if err := DB.Table("tbl_members").Where("LOWER(TRIM(email))=LOWER(TRIM(?)) and id not in (?) and is_deleted = 0 and tenant_id=?", email, userid, tenantid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

func (membermodel MemberModel) CheckNumberInMember(member *TblMember, number string, userid int, DB *gorm.DB, tenantid string) error {

	if userid == 0 {

		if err := DB.Debug().Table("tbl_members").Where("mobile_no = ? and tenant_id=? and is_deleted = 0", number, tenantid).First(&member).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Debug().Table("tbl_members").Where("mobile_no = ? and id not in (?) and tenant_id=? and is_deleted=0", number, userid, tenantid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

// Name already exists
func (membermodel MemberModel) CheckNameInMember(userid int, name string, DB *gorm.DB, tenantid string) (member Tblmember, err error) {

	if userid == 0 {

		if err := DB.Table("tbl_members").Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and tenant_id=? and is_deleted=0", name, tenantid).First(&member).Error; err != nil {

			return Tblmember{}, err
		}
	} else {

		if err := DB.Table("tbl_members").Where("LOWER(TRIM(username))=LOWER(TRIM(?)) and id not in (?) and tenant_id=?   and is_deleted=0", name, userid, tenantid).First(&member).Error; err != nil {

			return Tblmember{}, err
		}
	}

	return member, nil
}

// Member Group Update
func (membermodel MemberModel) MemberGroupUpdate(membergroup *Tblmembergroup, id int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_groups").Where("id=? and tenant_id=? ", id, tenantid).Updates(TblMemberGroup{Name: membergroup.Name, Slug: membergroup.Slug, Description: membergroup.Description, Id: membergroup.Id, ModifiedOn: membergroup.ModifiedOn, ModifiedBy: membergroup.ModifiedBy, IsActive: membergroup.IsActive}).Error; err != nil {

		return err
	}

	return nil
}

// Member Group Delete
func (membermodel MemberModel) DeleteMemberGroup(membergroup *Tblmembergroup, id int, DB *gorm.DB, tenantid string) error {

	if err := DB.Debug().Table("tbl_member_groups").Where("id=? and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "modified_by": membergroup.ModifiedBy}).Error; err != nil {

		return err

	}
	return nil
}

// get member group
func (membermodel MemberModel) GetGroupData(membergroup []Tblmembergroup, DB *gorm.DB, tenantid string) (membergrouplists []Tblmembergroup, err error) {

	var membergrouplist []Tblmembergroup

	if err := DB.Table("tbl_member_groups").Where("is_deleted = 0 and is_active = 1 and tenant_id=?", tenantid).Order("name").Find(&membergrouplist).Error; err != nil {

		return []Tblmembergroup{}, err

	}

	return membergrouplist, nil

}

// get member group
func (membermodel MemberModel) GetGroupDataid(membergroup Tblmembergroup, DB *gorm.DB, tenantid string, name string) (membergrouplists Tblmembergroup, err error) {

	var membergrouplist Tblmembergroup

	if err := DB.Table("tbl_member_groups").Where("is_deleted = 0 and is_active = 1 and tenant_id=? and name=?", tenantid, name).First(&membergrouplist).Error; err != nil {

		return Tblmembergroup{}, err

	}

	return membergrouplist, nil

}

// get member details
func (membermodel MemberModel) GetMemberDetailsByMemberId(MemberDetails *TblMember, memberId int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Where("is_deleted=0 and id = ? and tenant_id=?", memberId, tenantid).First(&MemberDetails).Error; err != nil {

		return err
	}

	return nil
}

// Get Member Details
func (membermodel MemberModel) MemberDetails(member *Tblmember, memberid int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Select("tbl_members.*,tbl_member_groups.name as group_name").Joins("left join tbl_member_groups on tbl_member_groups.id = tbl_members.member_group_id").Where("tbl_members.id=? and tbl_members.tenant_id=? and tbl_members.is_deleted=0", memberid, tenantid).First(&member).Error; err != nil {
		return err

	}

	return nil
}

// Member  IsActive Function
func (membermodel MemberModel) MemberStatus(memberstatus TblMember, memberid int, status int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Where("id=? and tenant_id=?", memberid, tenantid).UpdateColumns(map[string]interface{}{"is_active": status, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// MultiSelectedMemberDelete
func (membermodel MemberModel) MultiSelectedMemberDelete(member *TblMember, id []int, DB *gorm.DB, tenantid string) error {

	return DB.Transaction(func(tx *gorm.DB) error {

		if err := DB.Model(&member).Where("id in (?) and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

			return err

		}

		if err := DB.Table("tbl_member_profiles").Where("id in (?) and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": 1, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

			return err

		}

		return nil

	})
}

func (membermodel MemberModel) MultiMemberIsActive(memberstatus *TblMember, memberid []int, status int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Where("id in (?) and tenant_id=?", memberid, tenantid).UpdateColumns(map[string]interface{}{"is_active": status, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// Member la IsActive Function
func (membermodel MemberModel) MemberGroupIsActive(memberstatus *Tblmembergroup, memberid int, status int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_groups").Where("id=? and tenant_id=?", memberid, tenantid).UpdateColumns(map[string]interface{}{"is_active": status, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

// Group Name already exists
func (membermodel MemberModel) CheckNameInMemberGroup(member *Tblmembergroup, userid int, name string, DB *gorm.DB, tenantid string) error {

	if userid == 0 {

		if err := DB.Table("tbl_member_groups").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and tenant_id=? and is_deleted=0", name, tenantid).First(&member).Error; err != nil {

			return err
		}
	} else {

		if err := DB.Table("tbl_member_groups").Where("LOWER(TRIM(name))=LOWER(TRIM(?)) and id not in (?) and tenant_id=? and is_deleted=0", name, userid, tenantid).First(&member).Error; err != nil {

			return err
		}
	}

	return nil
}

// selected member group delete
func (membermodel MemberModel) MultiSelectedMemberDeletegroup(member *Tblmembergroup, id []int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_groups").Where("id in (?) and tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"is_deleted": member.IsDeleted, "deleted_on": member.DeletedOn, "deleted_by": member.DeletedBy}).Error; err != nil {

		return err

	}
	return nil
}

// selected member group status change
func (membermodel MemberModel) MultiMemberGroupIsActive(memberstatus *TblMemberGroup, memberid []int, status int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_groups").Where("id in (?) and tenant_id=?", memberid, tenantid).UpdateColumns(map[string]interface{}{"is_active": status, "modified_by": memberstatus.ModifiedBy, "modified_on": memberstatus.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (membermodel MemberModel) GetMemberProfile(memberId int, emailid string, profileId int, profileSlug string, DB *gorm.DB, tenantid string) (tblmember Tblmember, err error) {

	query := DB.Debug().Table("tbl_members").Preload("TblMemberProfile")

	if memberId != 0 {

		query = query.Where("is_deleted = 0 and id = ? and tenant_id=?", memberId, tenantid)

	} else if emailid != "" {

		query = query.Where("is_deleted = 0 and email = ? and tenant_id=?", emailid, tenantid)

	} else if profileSlug != "" {

		query = query.Where("tenant_id=? and is_deleted = 0 and id = (select member_id from tbl_member_profiles where is_deleted = 0 and profile_slug=?) ", tenantid, profileSlug)

	} else if profileId != 0 {

		query = query.Where("tenant_id=? and is_deleted = 0 and id = (select member_id from tbl_member_profiles where is_deleted = 0 and id=?)", tenantid, profileId)

	}

	query.First(&tblmember)

	if err := query.Error; err != nil {
		return tblmember, err
	}
	return tblmember, nil
}
func (membermodel MemberModel) AllMemberCount(DB *gorm.DB, tenantid string) (count int64, err error) {

	if err := DB.Table("tbl_members").Where("is_deleted = 0 and tenant_id=?", tenantid).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}

func (membermodel MemberModel) NewmemberCount(DB *gorm.DB, tenantid string) (count int64, err error) {

	if err := DB.Table("tbl_members").Where("tenant_id=? and is_deleted = 0 AND created_on >=?", tenantid, time.Now().AddDate(0, 0, -10)).Count(&count).Error; err != nil {

		return 0, err
	}

	return count, nil

}
func (membermodel MemberModel) ActiveMemberList(member []Tblmember, limit int, DB *gorm.DB, tenantid string) (members []Tblmember, err error) {

	if err := DB.Table("tbl_members").Where("tenant_id=? and is_deleted=0 and last_login=1 AND login_time >=?", tenantid, time.Now().UTC().Add(-8*time.Hour).Format("2006-01-02 15:04:05")).Find(&members).Limit(limit).Error; err != nil {

		return []Tblmember{}, err

	}

	return members, nil
}

func (membermodel MemberModel) FlexibleMemberUpdate(memberData map[string]interface{}, memberid int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Debug().Where("is_deleted = 0 and id = ? and tenant_id=?", memberid, tenantid).UpdateColumns(memberData).Error; err != nil {

		return err
	}

	return nil
}

func (membermodel MemberModel) FlexibleMemberProfileUpdate(memberProfileData map[string]interface{}, memberid int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_profiles").Where("is_deleted = 0 and member_id = ? and tenant_id=?", memberid, tenantid).UpdateColumns(memberProfileData).Error; err != nil {

		return err
	}

	return nil
}

func (membermodel MemberModel) MemberPasswordUpdate(memberData TblMember, memberId int, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_members").Where("is_deleted = 0 and id = ? and tenant_id=?", memberId, tenantid).UpdateColumns(map[string]interface{}{"password": memberData.Password, "modified_by": memberData.ModifiedBy, "modified_on": memberData.ModifiedOn}).Error; err != nil {

		return err
	}

	return nil
}

func (membermodel MemberModel) GetMemberSettings(DB *gorm.DB, tenantid string) (membersetting TblMemberSetting, err error) {

	if err := DB.Table("tbl_member_settings").Where(" tenant_id=?", tenantid).First(&membersetting).Error; err != nil {

		return TblMemberSetting{}, err
	}

	return membersetting, nil
}

func (membermodel MemberModel) UpdateMemberSetting(membersetting map[string]interface{}, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_settings").Where("id=1 and tenant_id=?").Updates(membersetting).Error; err != nil {

		return err
	}

	return nil
}

func (membermodel MemberModel) DeleteMemberProfile(memberid int, deletedby int, deletedOn time.Time, DB *gorm.DB, tenantid string) error {

	if err := DB.Table("tbl_member_profiles").Where("member_id=? and tenant_id=?", memberid, tenantid).UpdateColumns(map[string]interface{}{
		"is_deleted": 1, "deleted_by": deletedby, "deleted_on": deletedOn}).Error; err != nil {

		return err
	}

	return nil
}

// Remove member group in member
func (membermodel MemberModel) RemoveMemberGroupInMember(id int, ids []int, DB *gorm.DB, tenantid string) error {
	if id != 0 {
		if err := DB.Debug().Table("tbl_members").Where("member_group_id=? tenant_id=?", id, tenantid).UpdateColumns(map[string]interface{}{"member_group_id": 1}).Error; err != nil {

			return err

		}
	} else {
		if err := DB.Table("tbl_members").Where("member_group_id in (?) and tenant_id=?", ids, tenantid).UpdateColumns(map[string]interface{}{"member_group_id": 1}).Error; err != nil {

			return err

		}
	}
	return nil

}

func (membermodel MemberModel) Checkmembergroup(member *TblMember, id int, ids []int, DB *gorm.DB, tenantid string) error {

	query := DB.Table("tbl_members")
	if id != 0 {
		query = query.Where("member_group_id=? and tenant_id=? and is_deleted = 0", id, tenantid)
	} else {
		query = query.Where("member_group_id in (?) and tenant_id=? and is_deleted = 0", ids, tenantid)
	}
	if err := query.First(&member).Error; err != nil {
		return err
	}
	return nil
}
