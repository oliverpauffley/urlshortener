package hashing

import (
	"github.com/speps/go-hashids"
	"log"
)

// using the package https://hashids.org/go to turn integer ids to short alphanumeric hashes
func NewHashId(id int) (hashId string) {
	hashData := hashids.NewData()
	hashData.Salt = "this is the salt"
	hashData.MinLength = 5
	hash, err := hashids.NewWithData(hashData)
	if err != nil {
		log.Fatalf("Error setting hash id, %v", err)
	}
	hashId, err = hash.Encode([]int{id})
	if err != nil {
		log.Fatalf("error creating hashId from int Id, %v", err)
	}
	return
}
