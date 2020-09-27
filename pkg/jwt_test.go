package pkg

import (
	"colossus/models"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

// ModelSuite is
type ModelSuite struct {
	*suite.Model
}

// var jwtSecret = "testunnotechjwttoken"
var email, group = "test.unnotech.com", "testgroup"
var provider, providerID = "google", "999999"

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.New("app:jwt:fixtures", "../fixtures/jwt"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}

func generateToken(jwtSecret []byte) string {
	claims := Claims{
		Email:    email,
		Provider: provider,
		Group:    group,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Id:        providerID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}

func (ms *ModelSuite) Test_CreateJWTToken() {
	token, err := CreateJWTToken(email, provider, providerID, group)
	ms.NoError(err)
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return jwtSecret, nil
	})
	claim, ok := tokenClaims.Claims.(*Claims)
	ms.Equal(ok, true)
	ms.Equal(claim.Email, email)
	ms.Equal(claim.Group, group)
	ms.Equal(claim.Provider, provider)

	// Test with wrong jwtSecret
	_, err = jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return "wrongtoken", nil
	})
	ms.Error(err)
}

func (ms *ModelSuite) Test_CheckJWTToken() {
	token := generateToken(jwtSecret)
	// Lack of JWT at begin
	_, err := CheckJWTToken(token)
	ms.EqualError(err, "tokenstring should contains 'JWT'")

	// Good token
	tokenString := fmt.Sprintf("JWT %s", token)
	_, err = CheckJWTToken(tokenString)
	ms.NoError(err)

	// Wrong JWT token
	wrongSecret := []byte{'w', 'r', 'o', 'n', 'g'}
	wrongTokenString := fmt.Sprintf("JWT %s", generateToken(wrongSecret))
	_, err = CheckJWTToken(wrongTokenString)
	ms.EqualError(err, "signature validation failed")
}

func (ms *ModelSuite) Test_CheckPermission() {
	ms.LoadFixture("groups data")
	g := models.Group{}
	err := ms.DB.Where("name = ?", "defaultgroup").First(&g)
	if err != nil {
		return
	}
	g.GetPermission = []string{
		"/v1/get/", "/v1/get2/", "/v1/get3/",
	}
	g.UpdatePermission = []string{
		"/v1/update/", "/v1/update2/",
	}
	ms.DB.Update(&g)

	// Check for group permission
	getCase1 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/get/",
		"method":  "GET",
	}
	ms.NoError(CheckPermission(ms.DB, getCase1))

	getCase2 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/get2/",
		"method":  "GET",
	}
	ms.NoError(CheckPermission(ms.DB, getCase2))

	getCase3 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/getnull/",
		"method":  "GET",
	}
	ms.Error(CheckPermission(ms.DB, getCase3))

	getCase4 := map[string]interface{}{
		"group":   "defaultgroupnull",
		"urlPath": "/v1/get/",
		"method":  "GET",
	}
	ms.Error(CheckPermission(ms.DB, getCase4))

	updateCase1 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/update/",
		"method":  "PUT",
	}
	ms.NoError(CheckPermission(ms.DB, updateCase1))

	updateCase2 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/update/",
		"method":  "PATCH",
	}
	ms.NoError(CheckPermission(ms.DB, updateCase2))

	updateCase3 := map[string]interface{}{
		"group":   "defaultgroup",
		"urlPath": "/v1/update/",
		"method":  "POST",
	}
	ms.Error(CheckPermission(ms.DB, updateCase3))

	updateCase4 := map[string]interface{}{
		"group":   "defaultgroupnull",
		"urlPath": "/v1/update/",
		"method":  "PUT",
	}
	ms.Error(CheckPermission(ms.DB, updateCase4))

	// Check for OPA permission
	opaCase1 := map[string]interface{}{
		"group": "admin",
	}
	ms.NoError(CheckPermission(ms.DB, opaCase1))

	opaCase2 := map[string]interface{}{
		"group":  "defaultgroup",
		"path":   strings.Split("testjwt", "/"),
		"method": "GET",
	}
	ms.NoError(CheckPermission(ms.DB, opaCase2))

	opaCase3 := map[string]interface{}{
		"group":  "defaultgroup",
		"path":   strings.Split("testjwt/1", "/"),
		"method": "GET",
	}
	ms.NoError(CheckPermission(ms.DB, opaCase3))

	opaCase4 := map[string]interface{}{
		"group":  "defaultgroup",
		"path":   strings.Split("testjwt", "/"),
		"method": "POST",
	}
	ms.Error(CheckPermission(ms.DB, opaCase4))

	opaCase5 := map[string]interface{}{
		"group":  "testgroup",
		"path":   strings.Split("testjwt", "/"),
		"method": "GET",
	}
	ms.Error(CheckPermission(ms.DB, opaCase5))
}
