package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/teliaz/goapi/app/handlers"
	"github.com/teliaz/goapi/app/middlewares"
	"github.com/teliaz/goapi/app/models"
	"github.com/teliaz/goapi/config"
)

// App Structure
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize initializes the app with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
		config.DB.Username,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
		config.DB.Charset)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	if err != nil {
		log.Fatal("Could not connect database")
	}

	a.DB = models.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
}

func (a *App) setRouters() {

	// Health Check
	a.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(handlers.Ping)).Methods("GET")

	// Login Route
	a.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(handlers.Login)).Methods("POST")

	//Users routes
	a.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(handlers.CreateUser)).Methods("POST")
	a.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(handlers.GetUsers)).Methods("GET")
	a.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(handlers.GetUser)).Methods("GET")
	a.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(handlers.UpdateUser))).Methods("PUT")
	a.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(handlers.DeleteUser)).Methods("DELETE")

	/*
		// Assets Routes
		a.Get("/assets", a.handleRequest(handlers.GetAssets))
		a.Get("/assets/{id:[0-9]+}", a.handleRequest(handlers.GetAsset))
		a.Put("/assets/{id:[0-9]+}/isFavorite/{isFavorite}", a.handleRequest(handlers.UpdateAssetIsFavorite))
		a.Put("/assets/{id:[0-9]+}/title/{title}", a.handleRequest(handlers.UpdateAssetTitle))
		a.Delete("/assets/{id:[0-9]+}", a.handleRequest(handlers.DeleteAsset))
	*/
}

// Run the app on Mux router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}

// RequestHandlerFunction HandlerRequest extension
type RequestHandlerFunction func(db *gorm.DB, w http.ResponseWriter, r *http.Request)

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// Post wraps the router for POST method
func (a *App) Post(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("POST")
}

// Put wraps the router for PUT method
func (a *App) Put(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("PUT")
}

// Delete wraps the router for DELETE method
func (a *App) Delete(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("DELETE")
}

func (a *App) handleRequest(handler RequestHandlerFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(a.DB, w, r)
	}
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&models.Asset{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Data{})
	// db.Model(&Chart{}).AddForeignKey("AssetId", "Assets(id)", "CASCADE", "CASCADE")
	return db
}
