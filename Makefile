APP_NAME=mrtasker

build:
	./scripts/build.sh ${APP_NAME}
run: build
	./tmp/bin/${APP_NAME}

clean:
	rm -r ./tmp

test:
	go test ./...

test-coverage:
	@mkdir -p ./tmp/testing
	go test ./... -coverprofile=./tmp/testing/coverage.out

dep:
	go mod download

staticcheck:
	staticcheck