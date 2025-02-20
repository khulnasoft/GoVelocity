package main

import (
	"catalog/api"
	"catalog/config"
	"catalog/repository"
	"catalog/service"
	"go.khulnasoft.com/velocity"
	"go.khulnasoft.com/velocity/middleware/logger"
	"go.khulnasoft.com/velocity/middleware/recover"
)

func main() {
	conf, _ := config.NewConfig("./config/config.yaml")
	repo, _ := repository.NewMongoRepository(conf.Database.URL, conf.Database.DB, conf.Database.Timeout)
	service := service.NewProductService(repo)
	handler := api.NewHandler(service)

	r := velocity.New()
	r.Use(recover.New())
	r.Use(logger.New(logger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	r.Get("/products/{code}", handler.Get)
	r.Post("/products", handler.Post)
	r.Delete("/products/{code}", handler.Delete)
	r.Get("/products", handler.GetAll)
	r.Put("/products", handler.Put)
	r.Listen(":8080")
}
