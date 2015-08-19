package config

const (
	PrivKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	PubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	SecretKey   = "7Ml|j1(8apfeywbo|=<9'P8#7X<Na;&BJqmTqOXqphT4?$g@V!e2zuiH%M3#8zu3"
)
