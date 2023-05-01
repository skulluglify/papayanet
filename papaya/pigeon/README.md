# PIGEON

### database, repository, more

- CORS
- Basic Auth


### TODOs
- Session Model
```go
package models
type SessionModel struct {
	
	// RSA, JWE requirements
	PubKey []byte `gorm:"unique" config:"jsonable:false"`
	PrivKey []byte `gorm:"unique" config:"jsonable:false"`
}
```

configuration from structure
- jsonable, can convert into json by item