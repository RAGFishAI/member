package member

import (
	"fmt"
	"log"
	"testing"

	"github.com/RAGFishAI/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var SecretKey = "Secret123"

// Db connection
func DBSetup() (*gorm.DB, error) {

	dbConfig := map[string]string{
		"username": "postgres",
		"password": "postgres",
		"host":     "localhost",
		"port":     "5432",
		"dbname":   "nov_14",
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "user=" + dbConfig["username"] + " password=" + dbConfig["password"] +
			" dbname=" + dbConfig["dbname"] + " host=" + dbConfig["host"] +
			" port=" + dbConfig["port"] + " sslmode=disable TimeZone=Asia/Kolkata",
	}), &gorm.Config{})

	if err != nil {

		log.Fatal("Failed to connect to database:", err)

	}
	if err != nil {

		return nil, err

	}

	return db, nil
}

// test listmembergroup function
func TestListMemberGroup(t *testing.T) {

	db, _ := DBSetup()

	config := auth.Config{
		UserId: 1,
		// ExpiryTime: 2,
		ExpiryFlg: false,
		SecretKey: "Secret123",
		DB:        db,
		RoleId:    1,
	}

	Auth := auth.AuthSetup(config)

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Members Group", auth.CRUD, "1")

	member := MemberSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		membergroup, count, err := member.ListMemberGroup(MemberGroupListReq{Limit: 10, Offset: 0}, "1")

		if err != nil {

			panic(err)
		}

		fmt.Println(membergroup, count)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test createmembergroup function
func TestCreateMemberGroup(t *testing.T) {

	db, _ := DBSetup()

	config := auth.Config{
		UserId: 1,
		// ExpiryTime: 2,
		ExpiryFlg: false,
		SecretKey: "Secret123",
		DB:        db,
		RoleId:    1,
	}

	Auth := auth.AuthSetup(config)

	token, _ := Auth.CreateToken()

	Auth.VerifyToken(token, SecretKey)

	permisison, _ := Auth.IsGranted("Members Group", auth.CRUD, "1")

	member := MemberSetup(Config{
		DB:               db,
		AuthEnable:       true,
		PermissionEnable: true,
		Auth:             Auth,
	})
	if permisison {

		err := member.CreateMemberGroup(MemberGroupCreation{"sports", "indian team", 1, 1}, "1")

		if err != nil {

			panic(err)
		}

		fmt.Println(err)
	} else {

		log.Println("permissions enabled not initialised")

	}

}

// test createmembergroup function
func TestUpdateMemberGroup(t *testing.T) {

	db, _ := DBSetup()

	member := MemberSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := member.UpdateMemberGroup(MemberGroupCreationUpdation{Name: "Default", Description: "default group2", ModifiedBy: 1, IsActive: 1}, 2, "1")

	if err != nil {

		panic(err)
	}

	fmt.Println(err)

}

// test createmembergroup function
func TestDeleteMemberGroup(t *testing.T) {

	db, _ := DBSetup()

	member := MemberSetup(Config{
		DB:               db,
		AuthEnable:       false,
		PermissionEnable: false,
	})

	err := member.DeleteMemberGroup(2, 1, "1")

	if err != nil {

		panic(err)
	}

	fmt.Println(err)

}
