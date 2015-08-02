# go-test
test assigment

Test project TaskAPI

Create HTTP Rest API: 

1.	Use “echo” or “gin” frameworks for web handler.
2.	Implement simple middleware for contains header 'Authorization: Bearer testkey123' in each request. Otherwise return 403 and json struct with error
3.	Log each request include status code
4.	Implement persistence with MySQL and Gorm (https://github.com/jinzhu/gorm)
5.	Use Goose for DB migration (https://bitbucket.org/liamstask/goose)
6.	Implement save endpoint for Task object
7.	Implement update endpoint for Task object
8.	Implement get endpoint for Task object
9.	Implement delete endpoint for Task object (just update IsDeleted field)
10.	Use CORS (reply with header Access-Control-Allow-Origin: *)
11.	Add support for OPTION HTTP method for each endpoints
12.	Configure daemon over simple JSON config. Specify path as process flag for daemon. Required params: ListenAddress, DatabaseUri.

Task: 

```javascript
type Task struct {
    Id          int64
    Title       string
    Description string
    Priority    int
    CreatedAt   *time.Time 
    UpdatedAt   *time.Time 
    CompletedAt bool
    IsDeleted   bool
    IsCompleted bool
}
```
