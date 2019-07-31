package urlshortener

import (
	"github.com/oliverpauffley/urlshortener/database"
	"net/http"
)

// environment struct will be passed around and used for mocking
type Env struct {
	Db     database.Store
	Router *http.ServeMux
}
