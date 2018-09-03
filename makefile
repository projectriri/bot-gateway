PLUGIN_FILE := plugins.txt
BUILD_TAG := $(tag)
BUILD_DATE := `date +%FT%T%z`
COMMIT_SHA1 := `git rev-parse HEAD`
TAG := `git describe --tags`
PLUGIN_LIST := $(shell cat ${PLUGIN_FILE})
PLUGIN_OUTPUT_PATH := lib
all: so
	go get -d -v; \
	go build -ldflags " \
		-X github.com/projectriri/bot-gateway/main.BuildTag=${BUILD_TAG} \
		-X 'github.com/projectriri/bot-gateway/main.BuildDate=${BUILD_DATE}' \
		-X github.com/projectriri/bot-gateway/main.GitCommitSHA1=${COMMIT_SHA1} \
		-X github.com/projectriri/bot-gateway/main.GitTag=${TAG} \
		"
so:
	mkdir -p $(PLUGIN_OUTPUT_PATH)
	for plugin in $(PLUGIN_LIST); do \
		go get -d -v $$plugin; \
		cd $(shell go env GOPATH)/src/$$plugin; \
		go build -ldflags " \
			-X $$plugin.BuildTag=${BUILD_TAG} \
			-X '$$plugin.BuildDate=${BUILD_DATE}' \
			-X $$plugin.GitCommitSHA1=${COMMIT_SHA1} \
			-X $$plugin.GitTag=${TAG} \
			" --buildmode=plugin; \
		cd -; \
		mv $(shell go env GOPATH)/src/$$plugin/*.so $(PLUGIN_OUTPUT_PATH); \
	done
clean:
	rm -rf bot-gateway lib