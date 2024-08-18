APP_MRTASKER=mrtasker
APP_IMAGE_REVERTER=image-reverter

ifneq (,$(wildcard ./.secret))
    include .secret
    export
endif

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

dynamodb-create:
	terraform -chdir=infra/dynamodb apply -auto-approve

dynamodb-destroy:
	terraform -chdir=infra/dynamodb destroy -auto-approve

genetareMock:
	go generate -v -run="mockgen" ./...

# terraform -chdir=infra/dynamodb init
# terraform -chdir=infra/dynamodb plan
# terraform -chdir=infra/dynamodb apply
# terraform -chdir=infra/dynamodb destroy