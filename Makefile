PROJECT="gophkeeper"

default:
	echo ${PROJECT}

test:
	go test -v -count=1 ./...

.PHONY: default test cover
# покрытие тестами (исключены моки и protobuf сгененрированные файлы)
cover:
	go test -count=1 -coverprofile coverage.temp.out -coverpkg=./... ./...
	cat coverage.temp.out | grep -v "mock" | grep -v "pb.go" > coverage.out
	go tool cover -func coverage.out | grep total | awk '{print $3}'

# html представление покрытия (исключены моки и protobuf сгененрированные файлы)
cover-html:
	go test -count=1 -coverprofile coverage.temp.out -coverpkg=./... ./...
	cat coverage.temp.out | grep -v "mock" | grep -v "pb.go" > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

# покрытие с моками и protobuf файлами
cover-all:
	go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./...
	go tool cover -func coverage.out | grep total | awk '{print $3}'

.PHONY: protogen
protogen:
	protoc --go_out=. --go_opt=paths=source_relative   --go-grpc_out=. --go-grpc_opt=paths=source_relative   internal/proto/keeper.proto


.PHONY: migration_add
migration_add:
	# командная строка: make migration_add name=<migration_file_name>
	# бинарник migrate взят тут https://github.com/golang-migrate/migrate/releases
	cmd/server/migrations/migrate create -ext sql -dir cmd/server/migrations -seq $(name)

gen_keeper:
	mockgen -source=internal/server/domain/user/user.go -destination=internal/server/mocks/mock_user.go -package="mocks"
	mockgen -source=internal/server/domain/entity_code/entity_code.go -destination=internal/server/mocks/mock_entity_code.go  -package="mocks"
	mockgen -source=internal/server/domain/field/field.go -destination=internal/server/mocks/mock_field.go  -package="mocks"
	mockgen -source=internal/server/domain/entity/entity.go -destination=internal/server/mocks/mock_entity.go  -package="mocks"
