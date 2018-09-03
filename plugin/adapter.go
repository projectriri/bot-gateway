package plugin

type Adapter interface {
	BasePlugin
	Init(filename string)
	Start()
}
