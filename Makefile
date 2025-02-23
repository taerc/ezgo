phony:publish

export GO111MODULE = on
# alpha,release,final,auto
MAJOR?="0"
MINOR?="0"
PATCH?="5"
TAG_TYPE?="alpha"
TYPE_VERSION?="0"
DATETIME=`date +%Y%m%d%H%M`
GIT_TAG=v$(MAJOR).$(MINOR).$(PATCH)-$(TAG_TYPE).$(TYPE_VERSION)
MESSAGE?="优化Application代码"
BUILD?=build


update:
	@git pull
version:Makefile
	@echo "package ezgo" > version.go
	@echo "var version=\"$(GIT_TAG)\"" >> version.go

init:
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

## sqlmonitor
sqlmonitor:cmd/sqlmonitor/main.go
	@go build -o $(BUILD)/sqlmonitor cmd/sqlmonitor/main.go

## httpmod
httpmod:cmd/simphttp/main.go
	@go build -o $(BUILD)/httpmod cmd/simphttp/main.go

## noteme
noteme:cmd/noteme/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD}/noteme cmd/noteme/main.go
## gitlab
gitlabnote:cmd/gitlab/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD}/gitlabnote cmd/gitlab/main.go

## token
token:cmd/token/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD}/token cmd/token/main.go

expthreads:cmd/example/main.go | update
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD}/thread cmd/example/main.go

zentao:cmd/zentao/main.go
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -extldflags '-static'" -o ${BUILD}/zentao cmd/zentao/main.go
$(BUILD): 
	@mkdir -p $(BUILD)