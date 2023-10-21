package httpmod

import (
	"fmt"
	"path"
	"sync"

	"github.com/taerc/ezgo/conf"
)

type Resource struct {
	Id         string `json:"id" form:"id"`
	Path       string `json:"uri" form:"uri"`
	Type       int    `json:"type" form:"type" `
	Desc       string `json:"desc" form:"desc" `
	Attributes string `json:"attributes" form:"attributes" `
}
type ResourceType int

const (
	ResourceTypeFile ResourceType = iota + 0
	ResourceTypeAudio
	ResourceTypeVideo
	ResourceTypeLogo
	ResourceTypeApk
	ResourceTypeBinary
	ResourceTypeSqlite
	ResourceTypeLog
	ResourceTypeUnk
)

func isValidResourceType(t ResourceType) bool {

	if t >= ResourceTypeFile && t <= ResourceTypeUnk {
		return true
	}
	return false
}

var resourceMountPoint map[ResourceType]string

func init() {

	resourceMountPoint = make(map[ResourceType]string)
	resourceMountPoint[ResourceTypeFile] = "file"
	resourceMountPoint[ResourceTypeAudio] = "audio"
	resourceMountPoint[ResourceTypeVideo] = "video"
	resourceMountPoint[ResourceTypeLogo] = "logo"
	resourceMountPoint[ResourceTypeApk] = "apk"
	resourceMountPoint[ResourceTypeBinary] = "binary"
	resourceMountPoint[ResourceTypeSqlite] = "sqlite"
	resourceMountPoint[ResourceTypeLog] = "logs"
	resourceMountPoint[ResourceTypeUnk] = "default"
}

func WithComponentResource(c *conf.Configure) func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		var e error = nil
		var s string = ""
		for i := ResourceTypeFile; i <= ResourceTypeUnk; i += 1 {
			if s, e = Mkdirs(path.Join(conf.Config.ResourcePath, resourceMountPoint[i])); e != nil {
				Info(nil, M, fmt.Sprintf("init resource path [%s] failed, error is [%s]", s, e.Error()))
			}
		}
		Info(nil, M, "Init Resource Done!")
	}

}
