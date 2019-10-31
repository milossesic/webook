SERVICE_NAME=notificationsloadtest
BUILD_PATH=/go/src/git.corp.adobe.com/dc/notifications_load_test
BUILDER_TAG?=$(or $(sha),$(SERVICE_NAME)-builder)
IMAGE_TAG=$(SERVICE_NAME)-img
CODECOVERAGE_TAG?=$(or $(sha),$(SERVICE_NAME)-cc)
COVERAGE_OUT_FILE=coverage.out

default: ci

# login to artifactory using artifactory credentials
login:
		@docker login -u $(ARTIFACTORY_USER) -p $(ARTIFACTORY_API_TOKEN) docker-asr-release.dr.corp.adobe.com

pre-deploy-build:
		@echo "nothing is defined in pre-deploy-build step"

post-deploy-build:
		@echo "nothing is defined in post-deploy-build step"

# runs code coverage
code-coverage:
		@docker build -t $(CODECOVERAGE_TAG) -f Dockerfile.codecoverage .
		@docker run --rm -v `pwd`:$(BUILD_PATH):z -e ARTIFACTORY_USER -e ARTIFACTORY_API_TOKEN $(CODECOVERAGE_TAG) \
		sh -c 'echo "Per pkg coverage:"; go test -cover -coverprofile $(COVERAGE_OUT_FILE) $$(go list ./...) | grep -v "no test files"; echo "Per file coverage"; go tool cover -func=$(COVERAGE_OUT_FILE)'

# runs unit test
unit-test:
		@docker run --rm -v `pwd`:$(BUILD_PATH):z -e ARTIFACTORY_USER -e ARTIFACTORY_API_TOKEN $(BUILDER_TAG) \
		sh -c "ginkgo -r -trace -failFast -v --randomizeAllSpecs --randomizeSuites -p"

# runs unit test and code coverage
test: unit-test code-coverage

# builds the buildtime and runtime image
build: login
		@docker build -t $(BUILDER_TAG) -f Dockerfile.build .
		@docker run --rm -v `pwd`:$(BUILD_PATH):z -e ARTIFACTORY_USER -e ARTIFACTORY_API_TOKEN $(BUILDER_TAG)
		@docker build -t $(IMAGE_TAG) .

# The image tag for ci will be different with BYOJ, see https://jira.corp.adobe.com/browse/EON-4685
ci: IMAGE_TAG := $(if $(sha),$(IMAGE_TAG)-ci-$(sha),$(IMAGE_TAG))
ci: build test
