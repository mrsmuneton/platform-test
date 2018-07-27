package user

import (
	"strconv"
	"testing"
	"time"
)

func getUserStub() User {
	return User{CurrentPassword: "P0werpuff", Email: "cinnamon@nice.com", Name: "Ray May"}
}

func TestBryptEncryptDecrypt(t *testing.T) {
	password := "thisisMyPASS123!"
	hash, err := BcryptEncrypt(password)
	if err != nil {
		t.Error("Password encryption resulted in unexpected error")
	}
	decrypted := BcryptMatchPassword(hash, password)
	if decrypted == false {
		t.Error("Password decryption error")
	}
}

func TestCreateUserReturnsSuccess(t *testing.T) {
	var userStub = getUserStub()
	ts := time.Now()
	userStub.Email = userStub.Email + ts.String()
	var error = CreateUser(userStub)
	if error != false {
		t.Error("CreateUser returned unexpected error	")
	}
}

func TestCreateUserReturnsFailure(t *testing.T) {
	var userStub = getUserStub()
	userStub.Email = ""
	var error = CreateUser(userStub)
	if error == false {
		t.Error("CreateUser should return error	")
	}
}

func TestGetUserRecordReturnsSuccess(t *testing.T) {
	var userStub = getUserStub()
	var error = CreateUser(userStub)
	if error != false {
		t.Error("CreateUser returned unexpected error	")
	}
}

func TestUpdateUserReturnsSuccess(t *testing.T) {
	var userStub = getUserStub()
	userRecord, _ := GetUserRecordByEmail(userStub.Email)
	id := strconv.Itoa(userRecord.Id)
	_, err1 := UpdateUserFields(id, userStub)
	if err1 != false {
		t.Error("Update user returned unexpected error	")
	}
}

func TestDeleteUserReturnsSuccess(t *testing.T) {
	var userStub = getUserStub()
	var userRecord, err = GetUserRecordByEmail(userStub.Email)
	if err != false {
		t.Error("Delete user unexpectedly failed to fetch user record")
	}
	id := strconv.Itoa(userRecord.Id)
	var err1 = DeleteUser(id)
	if err1 != false {
		t.Error("Delete user returned unexpected error")
	}
}

func TestValidateMinimumFieldsFailReturnsError(t *testing.T) {
	var userStub = getUserStub()
	userStub.Name = ""
	_, e := ValidateUserMinimumFields(userStub)
	if e.Code != "Please provide all fields, including: Name" {
		t.Error("Name is a required field error not thrown")
	}
}

func TestValidateMinimumFieldsPresentReturnsUser(t *testing.T) {
	var userStub = getUserStub()
	u, e := ValidateUserMinimumFields(userStub)
	if e.Code != "" {
		t.Error("Validating minimum fields failed with an error")
	}
	if userStub != u {
		t.Error("Unexpected User Mutation")
	}
}

func TestValidateEmailIsUnique(t *testing.T) {
	email := "uniqueemail1@test.com"

	unique := ValidateEmailIsUnique(email)
	if unique != true {
		t.Error("ValidateEmailIsUnique encountered unexpected result")
	}
}
