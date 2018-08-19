package bingo

// 控制器接口
type ControllerInterface interface {
	Index(con *Context)
	Store(con *Context)
	Update(con *Context)
	Create(con *Context)
	Destroy(con *Context)
}

// 控制器父类
type Controller struct {
}

func (c *Controller) Index(con *Context) {

}

func (c *Controller) Create(con *Context) {

}

func (c *Controller) Store(con *Context) {

}

func (c *Controller) Update(con *Context) {

}

func (c *Controller) Destroy(con *Context) {

}
