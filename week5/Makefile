.PHONY: mock
mock:
	@mockgen -source=./webook/internal/service/user.go -package=svcmocks -destination=./webook/internal/service/mocks/user.mock.go
	@mockgen -source=./webook/internal/service/code.go -package=svcmocks -destination=./webook/internal/service/mocks/code.mock.go
	@mockgen -source=./webook/internal/service/sms/types.go -package=smsmocks -destination=./webook/internal/service/sms/mocks/sms.mock.go
	@mockgen -source=./webook/internal/service/sms/circuit_breaker/types.go -package=smsmocks -destination=./webook/internal/service/sms/mocks/circuit_breaker.mock.go
	@mockgen -source=./webook/internal/service/passwd_validate.go -package=svcmocks -destination=./webook/internal/service/mocks/passwd_validate.mock.go
	@mockgen -source=./webook/internal/repository/code.go -package=repomocks -destination=./webook/internal/repository/mocks/code.mock.go
	@mockgen -source=./webook/internal/repository/user.go -package=repomocks -destination=./webook/internal/repository/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/msg.go -package=repomocks -destination=./webook/internal/repository/mocks/msg.mock.go
	@mockgen -source=./webook/internal/repository/dao/user.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/dao/msg.go -package=daomocks -destination=./webook/internal/repository/dao/mocks/msg.mock.go
	@mockgen -source=./webook/internal/repository/cache/user.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/user.mock.go
	@mockgen -source=./webook/internal/repository/cache/code/types.go -package=cachemocks -destination=./webook/internal/repository/cache/mocks/code.mock.go
	@mockgen -package=redismocks -destination=./webook/internal/repository/cache/redismocks/cmd.mock.go github.com/redis/go-redis/v9 Cmdable
	@go mod tidy