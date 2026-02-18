package member

import (
	"strings"
	"time"
)

type Tblmembergroup struct {
	Id          int `gorm:"primaryKey;auto_increment;"`
	Name        string
	Slug        string
	Description string
	IsActive    int
	IsDeleted   int
	CreatedOn   time.Time `gorm:"DEFAULT:NULL"`
	CreatedBy   int
	ModifiedOn  time.Time `gorm:"DEFAULT:NULL"`
	ModifiedBy  int       `gorm:"DEFAULT:NULL"`
	DeletedOn   time.Time
	DeletedBy   int
	DateString  string `gorm:"-"`
}

// Function ListMemberGroup pass the arguments of limit,offset and filter (eg. keywords)
// It will return the all membergroup lists
func (member *Member) ListMemberGroup(listreq MemberGroupListReq, tenantid string) (membergroup []Tblmembergroup, MemberGroupCount int64, err error) {

	AuthErr := AuthandPermission(member)

	if AuthErr != nil {

		return []Tblmembergroup{}, 0, AuthErr
	}

	Membermodel.Userid = member.UserId
	Membermodel.DataAccess = member.DataAccess

	_, membercounts, _ := Membermodel.MemberGroupList(MemberGroupListReq{Limit: 0, Offset: 0, Keyword: listreq.Keyword, ActiveGroupsOnly: listreq.ActiveGroupsOnly}, member.DB, tenantid)

	membergrouplist, _, _ := Membermodel.MemberGroupList(listreq, member.DB, tenantid)

	var membergrouplists []Tblmembergroup

	for _, val := range membergrouplist {

		if !val.ModifiedOn.IsZero() {

			val.DateString = val.ModifiedOn.Format("02 Jan 2006 03:04 PM")

		} else {
			val.DateString = val.CreatedOn.Format("02 Jan 2006 03:04 PM")

		}

		membergrouplists = append(membergrouplists, val)

	}

	return membergrouplists, membercounts, nil

}

/*Create Member Group*/
func (member *Member) CreateMemberGroup(membergrpc MemberGroupCreation, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return AuthErr
	}

	if membergrpc.Name == "" {

		return ErrorEmpty
	}

	var membergroup TblMemberGroup

	membergroup.Name = membergrpc.Name

	membergroup.Slug = strings.ToLower(membergrpc.Name)

	membergroup.Description = membergrpc.Description

	membergroup.CreatedBy = membergrpc.CreatedBy

	membergroup.IsActive = membergrpc.IsActive
	membergroup.TenantId = tenantid

	membergroup.CreatedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MemberGroupCreate(&membergroup, member.DB)

	if err != nil {

		return err
	}

	return nil
}

/*Update Member Group*/
func (member *Member) UpdateMemberGroup(membergrpc MemberGroupCreationUpdation, id int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return AuthErr
	}

	if membergrpc.Name == "" {

		return ErrorEmpty
	}

	var membergroup Tblmembergroup

	membergroup.Id = id

	membergroup.Name = membergrpc.Name

	membergroup.Slug = strings.ToLower(membergrpc.Name)

	membergroup.Description = membergrpc.Description

	membergroup.ModifiedBy = membergrpc.ModifiedBy

	membergroup.IsActive = membergrpc.IsActive

	membergroup.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MemberGroupUpdate(&membergroup, id, member.DB, tenantid)

	if err != nil {

		return err
	}

	return nil
}

// delete member group
func (member *Member) DeleteMemberGroup(id int, modifiedBy int, tenantid string) error {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return AuthErr
	}

	var Tblmembergroup Tblmembergroup
	Tblmembergroup.ModifiedBy = modifiedBy
	err := Membermodel.DeleteMemberGroup(&Tblmembergroup, id, member.DB, tenantid)
	sterr := Membermodel.RemoveMemberGroupInMember(id, []int{}, member.DB, tenantid)
	if err != nil {

		return err
	}
	if sterr != nil {

		return sterr
	}

	return nil
}

func (member *Member) GetGroupData(tenantid string) (membergroup []Tblmembergroup, err error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return []Tblmembergroup{}, AuthErr
	}

	var memgroup []Tblmembergroup

	membergrouplist, _ := Membermodel.GetGroupData(memgroup, member.DB, tenantid)

	return membergrouplist, nil

}

// Get User id for mrmber group

func (member *Member) GetGroupDataid(tenantid string, name string) (membergroup Tblmembergroup, err error) {

	var memgroup Tblmembergroup

	membergrouplist, _ := Membermodel.GetGroupDataid(memgroup, member.DB, tenantid, name)

	return membergrouplist, nil

}

// member group is_active
func (member *Member) MemberGroupIsActive(memberid int, status int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}

	var tblmembergroup Tblmembergroup

	tblmembergroup.ModifiedBy = modifiedby

	tblmembergroup.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MemberGroupIsActive(&tblmembergroup, memberid, status, member.DB, tenantid)

	if err != nil {

		return false, err

	}
	return true, nil

}

// Check Group Name is already exits or not
func (member *Member) CheckNameInMemberGroup(id int, name string, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}

	var tblmembergroup Tblmembergroup

	err := Membermodel.CheckNameInMemberGroup(&tblmembergroup, id, name, member.DB, tenantid)

	if err != nil {

		return false, err

	}

	return true, nil
}

// MULTI SELECT MEMBERGROUP DELETE FUNCTION//
func (member *Member) MultiSelectedMemberDeletegroup(Memberid []int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}
	var tblmembergroup Tblmembergroup

	tblmembergroup.DeletedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	tblmembergroup.DeletedBy = modifiedby

	tblmembergroup.IsDeleted = 1

	err := Membermodel.MultiSelectedMemberDeletegroup(&tblmembergroup, Memberid, member.DB, tenantid)
	srerr := Membermodel.RemoveMemberGroupInMember(0, Memberid, member.DB, tenantid)

	if err != nil {

		return false, err
	}

	if srerr != nil {

		return false, srerr
	}

	return true, nil

}

// multi select membergroup status
func (member *Member) MultiSelectMembersgroupStatus(memberid []int, status int, modifiedby int, tenantid string) (bool, error) {

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}

	var memberstatus TblMemberGroup

	memberstatus.ModifiedBy = modifiedby

	memberstatus.ModifiedOn, _ = time.Parse("2006-01-02 15:04:05", time.Now().UTC().Format("2006-01-02 15:04:05"))

	err := Membermodel.MultiMemberGroupIsActive(&memberstatus, memberid, status, member.DB, tenantid)

	if err != nil {

		return false, err
	}

	return true, nil

}

func (member *Member) Membergroupcheckmember(memgrpid int, memgrpids []int, tenantid string) (bool, error) {

	var memberchk TblMember

	if AuthErr := AuthandPermission(member); AuthErr != nil {

		return false, AuthErr
	}
	err := Membermodel.Checkmembergroup(&memberchk, memgrpid, memgrpids, member.DB, tenantid)

	if err != nil {
		return false, err
	}

	if memberchk.Id == 0 {

		return false, nil
	}
	return true, nil
}
