package main

import (
	// "github.com/eugene/iizi_errand/pkg/models"
	"github.com/gofiber/fiber/v2"
)

func (r *Repository) Routes(app *fiber.App) {
	api := app.Group("api")
	api.Post("/", r.CreateUser)
	api.Post("/user/login", r.LoginHandler)
	// api.Use(models.JWTMiddleware())
	api.Post("/user/change-password", r.ChangePasswordHandler)
	api.Put("/user/update", r.UpdateUserProfile)
	api.Delete("/user/delete", r.DeleteUserProfile)
	api.Put("/errand-user/update", r.UpdateErrandRunnerProfile)
	api.Delete("/errand-runner/delete", r.DeleteErrandRunnerProfile)
	api.Post("/task/create", r.CreateTask)
	api.Post("/rating/:user_id/create", r.RateUser)
	api.Post("/rating/errand-runner/:errand_runner_id", r.RateErrandRunner)
	api.Get("/tasks", r.GetAllTasks)
	api.Get("/user/tasks", r.GetAllUserTasks)
	api.Put("/task/:task_id/update", r.UpdateTask)
	api.Delete("/task/:task_id/delete", r.DeleteTask)

    api.Post("/applications/apply/:task_id", r.CreateApplication)
    api.Get("/applications/errand/:task_id", r.GetApplicationsByErrandID)
    api.Patch("/applications/:appID", r.UpdateApplicationStatus)
}
