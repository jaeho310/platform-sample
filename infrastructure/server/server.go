package server

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"html/template"
	"io"
	"platform-sample/controller/api"
	"platform-sample/controller/web"
	_ "platform-sample/docs"
	"platform-sample/repository"
	"platform-sample/service"
)

type Server struct {
	MainDb *gorm.DB
}

type TemplateRenderer struct {
	templates *template.Template
}

func (server Server) Init() {
	e := echo.New()

	// web controller setting
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("view/templates/*.html")),
	}
	e.Renderer = renderer
	e.Static("/static", "view/static")
	web.WebController{}.Init(e)

	// api controller setting
	userController := server.InjectUserController()
	userController.Init(e.Group("/api/users"))

	// swagger setting
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":8395"))
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func (server Server) contextDB(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

func (server Server) InjectDb() *gorm.DB {
	return server.MainDb
}

func (server Server) InjectUserRepository() *repository.UserRepositoryImpl {
	return repository.UserRepositoryImpl{}.NewUserRepositoryImpl(server.InjectDb())
}

func (server Server) InjectUserService() *service.UserServiceImpl {
	return service.UserServiceImpl{}.NewUserServiceImpl(server.InjectUserRepository())
}

func (server Server) InjectUserController() *api.UserController {
	return api.UserController{}.NewUserController(server.InjectUserService())
}
