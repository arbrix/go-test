package tests

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token string `json:"token"`
}

func createNewsUser(username, password string) *User {
	return &User{username, password}
}

func TestLogin(t *testing.T) {
	Convey("Should be able to login", t, func() {
		user := createNewsUser("jonas", "1234")
		jsondata, _ := json.Marshal(user)
		post_data := strings.NewReader(string(jsondata))
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/v1/authentications", post_data)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, _ := client.Do(req)
		So(res.StatusCode, ShouldEqual, 200)

		Convey("Should be able to parse body", func() {
			body, err := ioutil.ReadAll(res.Body)
			defer res.Body.Close()
			So(err, ShouldBeNil)
			Convey("Should be able to get json back", func() {
				responseData := new(Response)
				err := json.Unmarshal(body, responseData)
				So(err, ShouldBeNil)

				Convey("Should be able to be authorized", func() {
					token := responseData.Token
					req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/tasks", nil)
					req.Header.Set("Authorization", "Bearer "+token)
					client = &http.Client{}
					res, _ := client.Do(req)
					So(res.StatusCode, ShouldEqual, 200)
				})
			})
		})
	})
	Convey("Should not be able to login with false credentials", t, func() {
		user := createNewsUser("jnwfkjnkfneknvjwenv", "wenknfkwnfknfknkfjnwkfenw")
		jsondata, _ := json.Marshal(user)
		post_data := strings.NewReader(string(jsondata))
		req, _ := http.NewRequest("POST", "http://127.0.0.1:8080/api/v1/authentications", post_data)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		res, _ := client.Do(req)
		So(res.StatusCode, ShouldEqual, 401)
	})

	Convey("Should not be able to authorize with false credentials", t, func() {
		token := ""
		req, _ := http.NewRequest("GET", "http://127.0.0.1:8080/api/v1/tasks", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		res, _ := client.Do(req)
		So(res.StatusCode, ShouldEqual, 401)
	})
}
