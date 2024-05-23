package main

import (
	"database/sql"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stonoy/PriceTracker/internal/database"
)

//go:embed static/*
var staticFiles embed.FS

type apiConfig struct {
	DB         *database.Queries
	Jwt_Secret string
}

func (cfg *apiConfig) clientHandler() http.Handler {
	fsys := fs.FS(staticFiles)
	contentStatic, _ := fs.Sub(fsys, "static")
	return http.FileServer(http.FS(contentStatic))

}

func main() {
	// load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, loading default configaration...")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("No port provided")
	}

	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Println("No jwt_secret provided")
	}

	// Creating config struct
	apiConfigObj := &apiConfig{
		Jwt_Secret: jwt_secret,
	}

	db_conn := os.Getenv("DB_CONN")
	if db_conn == "" {
		log.Println("No database connected")
	} else {
		// Database
		db, err := sql.Open("postgres", db_conn)
		if err != nil {
			log.Fatal("error loading postgres database")
		}

		dbQueries := database.New(db)

		apiConfigObj.DB = dbQueries

		log.Println("database connected")
	}

	// run our scrapper forever in different go routine
	go apiConfigObj.ourScrapper(50, 10*time.Second)

	//Handlers
	mainRouter := chi.NewRouter()

	// Making it cors enable
	mainRouter.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	restApiRouter := chi.NewRouter()

	// check health
	restApiRouter.Get("/checkHealth", checkHealth)
	restApiRouter.Get("/checkError", checkError)

	// users
	restApiRouter.Post("/register", apiConfigObj.registerUser)
	restApiRouter.Post("/login", apiConfigObj.loginUser)

	// products
	restApiRouter.Post("/createproducts", apiConfigObj.authMiddleware(apiConfigObj.createProduct))
	restApiRouter.Get("/userproducts", apiConfigObj.authMiddleware(apiConfigObj.productByUsers))
	restApiRouter.Put("/updatepriority/{productId}", apiConfigObj.authMiddleware(apiConfigObj.updateProductPriority))
	restApiRouter.Delete("/deleteproduct/{productId}", apiConfigObj.authMiddleware(apiConfigObj.deleteProduct))

	mainRouter.Mount("/v1", restApiRouter)

	// provides client
	mainRouter.Handle("/*", apiConfigObj.clientHandler())

	myServer := &http.Server{
		Addr:    ":" + port,
		Handler: mainRouter,
	}

	log.Printf("Server is listening on port %v", port)
	log.Fatal(myServer.ListenAndServe())
}
