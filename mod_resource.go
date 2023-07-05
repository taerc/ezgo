package ezgo

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	"path"
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

	resourceMountPoint = make(map[ResourceType]string, ResourceTypeUnk+1)
	resourceMountPoint[ResourceTypeFile] = "file"
	resourceMountPoint[ResourceTypeAudio] = "audio"
	resourceMountPoint[ResourceTypeVideo] = "video"
	resourceMountPoint[ResourceTypeLogo] = "logo"
	resourceMountPoint[ResourceTypeApk] = "apk"
	resourceMountPoint[ResourceTypeBinary] = "binary"
	resourceMountPoint[ResourceTypeSqlite] = "sqlite"
	resourceMountPoint[ResourceTypeUnk] = "unknown"

	var e error = nil
	var s string = ""
	for i := ResourceTypeFile; i <= ResourceTypeUnk; i += 1 {
		if s, e = Mkdirs(path.Join(conf.Config.ResourcePath, resourceMountPoint[i])); e != nil {
			Info(nil, M, fmt.Sprintf("init resource path [%s] failed, error is [%s]", s, e.Error()))
		}
	}
}

func getResourceTypePath(t ResourceType) (string, error) {
	dt := GetLocalDate()
	return Mkdirs(path.Join(conf.Config.ResourcePath, resourceMountPoint[t], dt))
}

func (r *Resource) Proc(ctx *gin.Context) ([]Resource, error) {
	//
	if e := ctx.BindJSON(r); e != nil {
		return nil, e
	}
	if !isValidResourceType(ResourceType(r.Type)) {
		return nil, errors.New("invalid resource type")
	}

	form, e := ctx.MultipartForm()
	if e != nil {
		return nil, e
	}

	files := form.File["file"]
	results := make([]Resource, 0)
	for _, f := range files {
		if d, e := getResourceTypePath(ResourceType(r.Type)); e == nil {
			if e = ctx.SaveUploadedFile(f, path.Join(d, f.Filename)); e != nil {
				return nil, e
			} else {
				results = append(results, Resource{
					Type:       r.Type,
					Path:       path.Join(d, f.Filename),
					Desc:       r.Desc,
					Attributes: r.Attributes,
					Id:         r.Id,
				})
			}

		} else {
			return nil, e
		}

		return results, nil
	}

	return results, nil
}
