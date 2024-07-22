APP_MRTASKER=mrtasker
APP_IMAGE_REVERTER=image-reverter

build-mrtasker:
	./scripts/build.sh ${APP_MRTASKER}
	./scripts/build.sh ${APP_IMAGE_REVERTER}
run-mrtasker: build-mrtasker
	./tmp/bin/${APP_MRTASKER}

build-image-reverter:
	./scripts/build.sh ${APP_IMAGE_REVERTER}

run-image-reverter: build-image-reverter
	./tmp/bin/${APP_IMAGE_REVERTER}

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