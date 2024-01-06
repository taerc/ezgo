phony:publish

export GO111MODULE = on
# alpha,release,final,auto
MAJOR?="0"
MINOR?="0"
PATCH?="4"
TAG_TYPE?="alpha"
TYPE_VERSION?="24"
DATETIME=`date +%Y%m%d%H%M`
GIT_TAG=v$(MAJOR).$(MINOR).$(PATCH)-$(TAG_TYPE).$(TYPE_VERSION)
MESSAGE?="优化Application代码"
BUILD?=build


version:Makefile
	@echo "package ezgo" > version.go
	@echo "var version=\"$(GIT_TAG)\"" >> version.go

init:
	@rm -f go.mod go.sum
	@go mod init github.com/taerc/ezgo
	@go generate ./db/ent
	@go  mod download
	@go mod tidy



publish:version
#linux系统 build
	git add .
	git commit -m $(MESSAGE)
	git push
	git tag -a $(GIT_TAG) -m $(MESSAGE)
	git push origin --tags

## testing
server:cmd/chatserver/main.go | $(BUILD)
	@go build -o $(BUILD)/server cmd/chatserver/main.go

client:cmd/tcpclient/main.go | $(BUILD)
	@go build -o $(BUILD)/client cmd/tcpclient/main.go

test:
	@/usr/local/go/bin/go test -all -timeout 1h -run ^TestGenerateLicence$ github.com/taerc/ezgo/licence github.com/taerc/ezgo -count=1 -v

## columns
columns:cmd/columns/main.go
	@go build -o $(BUILD)/columns cmd/columns/main.go

$(BUILD): 
	@mkdir -p $(BUILD)