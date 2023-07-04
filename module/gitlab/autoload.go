package gitlab

import (
	"github.com/taerc/ezgo"
	"sync"
)

var M string = "GITLAB"

func WithModuleGitLab() func(wg *sync.WaitGroup) {
	return func(wg *sync.WaitGroup) {
		wg.Done()
		route := ezgo.Group("/gitlab/hook/event/")
		ezgo.SetPostProc(route, "/push", &PushEventPayload{})
		ezgo.SetPostProc(route, "/push_tag", &TagEventsLoad{})
		ezgo.Info(nil, M, "Load finished!")
	}
}
