phony:ezgo

export GO111MODULE = on
# alpha,release,final,auto
MAJOR?="0"
MINOR?="0"
PATCH?="1"
TAG_TYPE?="alpha"
TYPE_VERSION?="5"
MESSAGE?="update log"
DATETIME=`date +%Y%m%d%H%M`
GIT_TAG=v$(MAJOR).$(MINOR).$(PATCH)-$(TAG_TYPE).$(TYPE_VERSION)


ezgo:
	go test -v *.go
publish:
#linux系统 build
	git tag -a $(GIT_TAG) -m $(MESSAGE)
	git push origin --tags