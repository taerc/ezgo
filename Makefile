phony:init

export GO111MODULE = on
# alpha,release,final,auto
MESSAGE?="更新版本号, 增加钉钉默认告警模块"
MAJOR?="0"
MINOR?="0"
PATCH?="3"
TAG_TYPE?="alpha"
TYPE_VERSION?="0"
DATETIME=`date +%Y%m%d%H%M`
GIT_TAG=v$(MAJOR).$(MINOR).$(PATCH)-$(TAG_TYPE).$(TYPE_VERSION)


version:
	@echo "package ezgo" > version.go
	@echo "var version=\"$(GIT_TAG)\"" >> version.go

init:
	@rm -f go.mod go.sum
	@git mod init ezgo
	@git mod download
	@git mod tidy

publish:version
#linux系统 build
	git add .
	git commit -m $(MESSAGE)
	git push
	git tag -a $(GIT_TAG) -m $(MESSAGE)
	git push origin --tags