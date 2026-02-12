package member

// Check Email is already exits or not
func (member *Member) CheckEmailInMember(id int, email string, tenantid string) (bool, error) {

	var cmember TblMember

	err := Membermodel.CheckEmailInMember(&cmember, email, id, member.DB, tenantid)

	if err != nil {

		return false, err
	}

	return true, nil

}

// function to Check Number is already exits or not
func (member *Member) CheckNumberInMember(id int, number string, tenantid string) (bool, error) {

	var cmember TblMember

	err := Membermodel.CheckNumberInMember(&cmember, number, id, member.DB, tenantid)

	if err != nil {

		return false, err

	}

	return true, nil
}

// Check Name is already exits or not
func (member *Member) CheckNameInMember(id int, name string, tenantid string) (bool, error) {

	cmember, err := Membermodel.CheckNameInMember(id, name, member.DB, tenantid)

	if err != nil {
		return false, err
	}
	if cmember.Id == 0 {

		return false, err
	}

	return true, nil
}
