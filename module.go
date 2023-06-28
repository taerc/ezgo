package ezgo

import "sync"

// ComponentLoad @description 应用系统依赖的第三方组件，如数据库等
type Component func(wg *sync.WaitGroup)

// ModuleRegister @description 应用系统的业务模块
type Module func(wg *sync.WaitGroup)

func LoadComponent(coms ...Component) {

	var wg sync.WaitGroup

	for _, c := range coms {
		wg.Add(1)
		go c(&wg)
	}
	wg.Wait()

}

func LoadModule(modules ...Module) {

	var wg sync.WaitGroup

	for _, m := range modules {
		wg.Add(1)
		go m(&wg)
	}
	wg.Wait()

}
