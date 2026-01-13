package chirpyapi

import (
	"sync/atomic"

	"github.com/AugustoMagro/gowebserver/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret_key     string
}

func NewClient(db *database.Queries, platform string, secret_key string) ApiConfig {
	return ApiConfig{
		fileserverHits: atomic.Int32{},
		db:             db,
		platform:       platform,
		secret_key:     secret_key,
	}
}
