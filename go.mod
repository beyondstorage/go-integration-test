module github.com/aos-dev/go-integration-test/v3

go 1.15

require (
	github.com/aos-dev/go-storage/v3 v3.4.2
	github.com/google/uuid v1.2.0
	github.com/smartystreets/goconvey v1.6.4
)

replace github.com/aos-dev/go-storage/v3 => ../go-storage
