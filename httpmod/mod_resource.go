package httpmod

import (
	"errors"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/taerc/ezgo/conf"
	ezgo "github.com/taerc/ezgo/pkg"
)

func getResourceTypePath(t ResourceType) (string, error) {
	dt := ezgo.GetLocalDate()
	return ezgo.Mkdirs(path.Join(conf.Config.ResourcePath, resourceMountPoint[t], dt))
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
