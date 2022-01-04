package hook

var (
	DefaultHooks = NewHooks("default")
)

func Exist(name string) bool {
	return DefaultHooks.Exist(name)
}

func Add(name string, handler Handler) error {
	return DefaultHooks.Add(name, handler)
}

func Del(name string) {
	DefaultHooks.Del(name)
}

func Call(name string, obj interface{}) (interface{}, error) {
	return DefaultHooks.Call(name, obj)
}
