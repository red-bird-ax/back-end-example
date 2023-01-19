module github.com/red-bird-ax/poster/auth

go 1.18

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-chi/chi v1.5.4
	github.com/red-bird-ax/poster/utils v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible // indirect
	github.com/go-redis/redis/v9 v9.0.0-rc.1 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/google/uuid v1.3.0 // indirect
)

replace github.com/red-bird-ax/poster/utils => ../utils
