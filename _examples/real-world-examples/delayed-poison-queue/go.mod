module delayed-poison-queue

go 1.23.0

require (
	github.com/ThreeDotsLabs/watermill v1.3.7
	github.com/ThreeDotsLabs/watermill-sql/v3 v3.0.4-0.20240906122508-e0de57ad1d8e
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.2
)

require (
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/lithammer/shortuuid/v3 v3.0.7 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sony/gobreaker v1.0.0 // indirect
)

replace github.com/ThreeDotsLabs/watermill => ../../../
replace github.com/ThreeDotsLabs/watermill-sql/v3 => ../../../../watermill-sql
