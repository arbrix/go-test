# go-test
test assigment

Short coding exercise as part of interview process - please use this repo to commit all of your work. Additional question may be provided via issues to this repo. Good luck and have fun! :)

Create HTTP Rest API:

* Use echo or gin for web handler
* Implement login endpoint with JWT token, and simple middleware that checks header for 'Authorization: Bearer %jwt_token%' in each request. Otherwise return 403 and json struct with error
* Implement endpoint that will use Facebook as OAuth2 provider
    * Implement authorize and create local user feature, using data from FB. Newly created user should be issued jwt token
    * For already existing users implement login feature.
* Log each request including status code
* Implement persistence with MySQL and Gorm (https://github.com/jinzhu/gorm)
* Use Goose or other tool of choice for DB migration (https://bitbucket.org/liamstask/goose)
* Implement save endpoint for Task object
* Implement update endpoint for Task object
* Implement get endpoint for Task object
* Implement delete endpoint for Task object (just update IsDeleted field)
* Use CORS (reply with header Access-Control-Allow-Origin: *)
* Add support for OPTION HTTP method for each endpoints
* Configure daemon over simple JSON config. Specify path as process flag for daemon. Required params: ListenAddress, DatabaseUri.
* Put in comments below description of taken architectural decisions

Task:

```
type Task struct {
    Id          int64
    Title       string
    Description string
    Priority    int
    CreatedAt   *time.Time
    UpdatedAt   *time.Time
    CompletedAt *time.Time
    IsDeleted   bool
    IsCompleted bool
}
```

### Architecture

copy from github.com/dorajistyle/goyangi

Structure
```
- api           //REST API, routes, resource
- app           //Core of application
- common        //General interfaces
- config        //general system configuration
- db            //gorm init and for goose migration component
- model         //description of DB entity
- service       //business logic
- util          //helpers and tools
- main.go       //application
- README.md
```

### Install
```
go get github.com/arbrix/go-test
go get github.com/labstack/echo
go get github.com/jinzhu/gorm
go get github.com/go-sql-driver/mysql
go get bitbucket.org/liamstask/goose/cmd/goose
go get golang.org/x/crypto/bcrypt
go get github.com/dgrijalva/jwt-go
go get golang.org/x/oauth2
go get github.com/stretchr/testify/assert
goose -path="src/github.com/arbrix/go-test/db" up #table: test; user: test; password: test; should exists or change dbconf.yml and conf.json
```

### Run
```
cp src/github.com/arbrix/go-test/conf.json ./
# vim conf.json
go-test
```
