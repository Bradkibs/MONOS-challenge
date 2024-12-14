package routes

import (
	"github.com/Bradkibs/MONOS-challenge/controllers"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupBranchRoutes(app *fiber.App, db *pgxpool.Pool) {

	branchController := controllers.BranchController{DB: db}

	branchGroup := app.Group("/branches")

	branchGroup.Post("/add", branchController.AddBranch)
	branchGroup.Get("/", branchController.GetBranches)
	branchGroup.Put("/update", branchController.UpdateBranch)
	branchGroup.Delete("/delete", branchController.DeleteBranch)
	branchGroup.Put("/update-for-subscription", branchController.UpdateBranchesForSubscription)
}
