package handler

import (
	"net/http"

	"github.com/funcTomas/hermes/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(factory service.Factory) http.Handler {
	r := chi.NewRouter()
	// 记录请求日志
	r.Use(middleware.Logger)
	// 恢复 panic
	r.Use(middleware.Recoverer)

	/*
		r.Route("/plan", func(r chi.Router) {
			planHandler := NewPlanHandler()
			r.Post("/add", planHandler.CreatePlan)
			r.Get("/detail/{id}", planHandler.UpdatePlan)
			r.Delete("/{id}", planHandler.DeletePlan)
		})
		r.Route("/monitor", func(r chi.Router) {
			monitorHandler := NewMonitorHandler()
			r.Get("/health", monitorHandler.GetHealth)
			r.Get("/metrics", monitorHandler.GetMetrics)
		})
	*/

	r.Route("/user", func(r chi.Router) {
		handler := NewUserHandler(factory)
		r.Post("/add", handler.UserAdd)
		r.Post("/entergroup", handler.UserEnterGroup)
	})

	return r
}
