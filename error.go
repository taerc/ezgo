package ezgo

type EM struct {
	C int
	M string
}

var errorRegistry map[int]string

func init() {
	errorRegistry = make(map[int]string)
}

//func ErrorRegister(code int, message string) {
//	// code checking ....
//
//}

// ErrorRegister

func ErrorRegister(em ...EM) {
	for _, e := range em {
		errorRegistry[e.C] = e.M
	}
}

func GetErrorMessage(code int) string {

	if em, ok := errorRegistry[code]; ok {
		return em
	} else {
		return "unknown"
	}
}
