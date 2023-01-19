module github.com/red-bird-ax/poster/users

go 1.18

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/go-chi/chi v1.5.4
	github.com/go-ozzo/ozzo-dbx v1.5.0
	github.com/go-ozzo/ozzo-validation v3.6.0+incompatible
	github.com/google/uuid v1.3.0
	github.com/lib/pq v1.10.7
	github.com/red-bird-ax/poster/utils v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.1.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
)

replace github.com/red-bird-ax/poster/utils => ../utils
