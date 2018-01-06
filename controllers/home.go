package controllers

type HomeController struct {
	BaseController
}

func (c *HomeController) Prepare() {
	c.BaseController.Prepare()

	c.Layout = "jeasyui/layout.html"
	c.Data["StaticPath"] = "/static/jeasyui"
}

func (c *HomeController) Finish() {
	c.Controller.Finish()
}

func (c *HomeController) Get() {
	c.RedirectWithURL("/admincp/login")
}

func (c *HomeController) Index() {
	c.TplName = "jeasyui/index.html"
}
