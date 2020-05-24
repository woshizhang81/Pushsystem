package module

/*
	代表module 的基类
*/
type Module interface {
	Start()
	Stop()
}


/*
	代表module 的基类
*/
type Manager interface {
	ModuleNotice()
}
