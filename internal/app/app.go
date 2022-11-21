package app

import (
	"log"
	"net/http"
	"txp/restapistarter/internal/app/module/auth"
	"txp/restapistarter/internal/app/module/content"
	"txp/restapistarter/internal/app/module/user"
	"txp/restapistarter/internal/pkg/constant"
	"txp/restapistarter/internal/pkg/middleware"
	_routerInternal "txp/restapistarter/internal/pkg/router"
	"txp/restapistarter/pkg/data/nosql/mongodb"
	"txp/restapistarter/pkg/data/sql/postgres"
	"txp/restapistarter/pkg/router"
)

// App struct
type App struct {
	DBClient      *mongodb.DBClient
	router        *router.Router
	Middlewares   []any
	AuthModule    *auth.AuthModule
	UserModule    *user.UserModule
	ContentModule *content.ContentModule
}

func (a *App) initDB() {
	postgres.InitDBClient()
	a.DBClient = mongodb.NewDBClient()
}

func (a *App) initMiddlewares() {
	authMiddleWare := middleware.NewAuthMiddleware(a.AuthModule.Service)
	a.Middlewares = append(a.Middlewares, authMiddleWare)
}

func (a *App) initModules() {
	a.UserModule = user.NewUserModule(a.DBClient.DB)
	a.AuthModule = auth.NewAuthModule(a.UserModule.Service)
	a.ContentModule = content.NewContentModule(a.DBClient.DB)
}

func (a *App) initRouters() {
	m := a.Middlewares[0].(*middleware.AuthMiddleware)
	_routerInternal.RegisterUserRoutes(a.router, constant.V1, a.UserModule, m)
	_routerInternal.RegisterContentRoutes(a.router, constant.V1, a.ContentModule)
}

// Init app
func (a *App) InitComponents() {
	a.initDB()
	a.initDB()
	a.router = router.NewRouter()
	a.initModules()
	a.initMiddlewares()
}

// Run app
func (a *App) Run() {
	err := http.ListenAndServe(
		":8080",
		a.router.Mux,
	)
	if err != nil {
		log.Fatal(err)
	}
}
