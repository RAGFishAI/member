package member

import (
	"time"

	"github.com/RAGFishAI/member/migration"
	"github.com/google/uuid"
)

// MemberSetup used initialize member configruation
func MemberSetup(config Config) *Member {

	migration.AutoMigration(config.DB, config.DataBaseType)

	return &Member{
		AuthEnable:       config.AuthEnable,
		Permissions:      config.Permissions,
		PermissionEnable: config.PermissionEnable,
		Auth:             config.Auth,
		DB:               config.DB,
	}

}

// list member
func (member *Member) ListMembers(offset int, limit int, filter Filter, flag bool, TenantId string) (memb []Tblmember, totoalmember int64, err error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return []Tblmember{}, 0, AuthErr
	}

	Membermodel.Userid = member.UserId
	Membermodel.DataAccess = member.DataAccess

	memberlist, _, _ := Membermodel.MembersList(limit, offset, filter, flag, member.DB, TenantId)

	_, Total_users, _ := Membermodel.MembersList(0, 0, filter, flag, member.DB, TenantId)

	var memberlists []Tblmember

	for _, val := range memberlist {

		// var first = val.FirstName
		// var last = val.LastName
		// var firstn = strings.ToUpper(first[:1])
		// var lastn string
		// if val.LastName != "" {
		// 	lastn = strings.ToUpper(last[:1])
		// }
		// var Name = firstn + lastn
		// val.NameString = Name
		val.CreatedDate = val.CreatedOn.Format("02 Jan 2006 03:04 PM")
		if !val.ModifiedOn.IsZero() {
			val.ModifiedDate = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")
		} else {
			val.ModifiedDate = val.CreatedOn.Format("02 Jan 2006 03:04 PM")
		}

		memberlists = append(memberlists, val)

	}

	return memberlists, Total_users, nil

}

// Create Member
func (member *Member) CreateMember(Mc MemberCreationUpdation) (Tblmember, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return Tblmember{}, AuthErr
	}
	uvuid := (uuid.New()).String()

	var cmember Tblmember
	cmember.Uuid = uvuid
	cmember.ProfileImage = Mc.ProfileImage
	cmember.ProfileImagePath = Mc.ProfileImagePath
	cmember.MemberGroupId = Mc.GroupId
	cmember.FirstName = Mc.FirstName
	cmember.LastName = Mc.LastName
	cmember.Email = Mc.Email
	cmember.MobileNo = Mc.MobileNo
	cmember.IsActive = Mc.IsActive
	cmember.Username = Mc.Username
	if Mc.Password != "" {
		hash_pass := hashingPassword(Mc.Password)
		cmember.Password = hash_pass
	}
	cmember.CreatedBy = Mc.CreatedBy
	cmember.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	cmember.StorageType = Mc.StorageType
	cmember.TenantId = Mc.TenantId
	err := Membermodel.MemberCreate(&cmember, member.DB)
	if err != nil {

		return Tblmember{}, err
	}

	return cmember, nil

}

// Update Member
func (member *Member) UpdateMember(Mc MemberCreationUpdation, id int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return AuthErr
	}

	uvuid := (uuid.New()).String()

	var umember Tblmember
	umember.Uuid = uvuid
	umember.Id = id
	umember.MemberGroupId = Mc.GroupId
	umember.FirstName = Mc.FirstName
	umember.LastName = Mc.LastName
	umember.Email = Mc.Email
	umember.MobileNo = Mc.MobileNo
	umember.ProfileImage = Mc.ProfileImage
	umember.ProfileImagePath = Mc.ProfileImagePath
	umember.IsActive = Mc.IsActive
	umember.ModifiedBy = Mc.ModifiedBy
	umember.Username = Mc.Username
	password := Mc.Password
	if password != "" {
		hash_pass := hashingPassword(password)
		umember.Password = hash_pass
	}
	umember.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	umember.StorageType = Mc.StorageType
	err := Membermodel.UpdateMember(&umember, member.DB, tenantid)
	if err != nil {

		return err
	}

	return nil
}

// delete member
func (member *Member) DeleteMember(id int, modifiedBy int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return AuthErr
	}

	var dmember Tblmember

	dmember.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	dmember.DeletedBy = modifiedBy

	err := Membermodel.DeleteMember(&dmember, id, member.DB, tenantid)

	Membermodel.DeleteMemberProfile(id, modifiedBy, dmember.DeletedOn, member.DB, tenantid)

	if err != nil {

		return err
	}

	return nil
}

// Get member data
func (member *Member) GetMemberDetails(id int, tenantid string) (members Tblmember, err error) {

	var memberdata Tblmember

	err = Membermodel.MemberDetails(&memberdata, id, member.DB, tenantid)
	if err != nil {

		return Tblmember{}, err
	}

	return memberdata, nil

}

// member is_active
func (member *Member) MemberStatus(memberid int, status int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return false, AuthErr
	}

	var memberstatus TblMember
	memberstatus.ModifiedBy = modifiedby
	memberstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MemberStatus(memberstatus, memberid, status, member.DB, tenantid)
	if err != nil {
		return false, err
	}

	return true, nil

}

// multiselecte member delete
func (member *Member) MultiSelectedMemberDelete(Memberid []int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return false, AuthErr
	}

	var members TblMember
	members.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	members.DeletedBy = modifiedby
	members.IsDeleted = 1

	err := Membermodel.MultiSelectedMemberDelete(&members, Memberid, member.DB, tenantid)
	if err != nil {

		return false, err
	}

	return true, nil
}

// multiselecte member status change
func (member *Member) MultiSelectMembersStatus(memberid []int, status int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}

	var memberStatus TblMember
	memberStatus.ModifiedBy = modifiedby
	memberStatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MultiMemberIsActive(&memberStatus, memberid, status, member.DB, tenantid)
	if err != nil {

		return false, err
	}

	return true, nil

}

func (member *Member) GetMemberAndProfileData(memberId int, emailid string, profileId int, profileSlug string, tenantid string) (Tblmember, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return Tblmember{}, AuthErr
	}

	profile, err := Membermodel.GetMemberProfile(memberId, emailid, profileId, profileSlug, member.DB, tenantid)
	if err != nil {
		return Tblmember{}, err
	}

	return profile, nil
}

// Active MemberList Function//
func (member *Member) ActiveMemberList(limit int, tenantid string) (memberdata []Tblmember, err error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return []Tblmember{}, AuthErr
	}

	var members []Tblmember
	activememlist, err := Membermodel.ActiveMemberList(members, limit, member.DB, tenantid)

	var memberlist []Tblmember
	for _, val := range activememlist {
		val.DateString = val.LoginTime.Format("02 Jan 2006 03:04 PM")
		memberlist = append(memberlist, val)
	}

	if err != nil {
		return []Tblmember{}, err
	}

	return memberlist, nil

}

// Member flexible update functionality
func (member *Member) MemberFlexibleUpdate(memberData map[string]interface{}, memberId, modifiedBy int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return AuthErr
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	memberData["modified_on"] = currentTime
	memberData["modified_by"] = modifiedBy
	err := Membermodel.FlexibleMemberUpdate(memberData, memberId, member.DB, tenantid)
	if err != nil {
		return err
	}

	return nil

}

// Memeber profile flexible update
func (member *Member) MemberProfileFlexibleUpdate(memberProfileData map[string]interface{}, memberId, modifiedBy int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return AuthErr
	}

	currentTime, _ := time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	memberProfileData["modified_on"] = currentTime
	memberProfileData["modified_by"] = modifiedBy
	err := Membermodel.FlexibleMemberProfileUpdate(memberProfileData, memberId, member.DB, tenantid)
	if err != nil {
		return err
	}

	return nil
}

// Member password update functionality
func (member *Member) MemberPasswordUpdate(newPassword, confirmPassword, oldPassword string, memberId, modifiedBy int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {
		return AuthErr
	}

	var memberData TblMember
	if err := Membermodel.GetMemberDetailsByMemberId(&memberData, memberId, member.DB, tenantid); err != nil {
		return err
	}

	// if err := bcrypt.CompareHashAndPassword([]byte(memberData.Password), []byte(oldPassword)); err != nil {
	// 	return err
	// }

	if newPassword != confirmPassword {
		return ErrorPassMissMatch
	}

	hash_pass := hashingPassword(confirmPassword)
	// if err := bcrypt.CompareHashAndPassword([]byte(hash_pass), []byte(oldPassword)); err != nil {
	// 	return err
	// }

	memberData.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))
	memberData.ModifiedBy = modifiedBy
	memberData.Password = hash_pass
	err := Membermodel.MemberPasswordUpdate(memberData, memberId, member.DB, tenantid)
	if err != nil {
		return err
	}

	return nil

}
