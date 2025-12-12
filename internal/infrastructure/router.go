package infrastructure

import (
	"log"

	"github.com/ChangerzaryX1602/sysware-test/internal/handlers"
	"github.com/ChangerzaryX1602/sysware-test/pkg/auth"
	"github.com/ChangerzaryX1602/sysware-test/pkg/models"
	"github.com/ChangerzaryX1602/sysware-test/pkg/permission"
	"github.com/ChangerzaryX1602/sysware-test/pkg/role"
	"github.com/ChangerzaryX1602/sysware-test/pkg/role_permission"
	"github.com/ChangerzaryX1602/sysware-test/pkg/todo"
	"github.com/ChangerzaryX1602/sysware-test/pkg/user"
	"github.com/ChangerzaryX1602/sysware-test/pkg/user_role"
	"github.com/ChangerzaryX1602/sysware-test/pkg/utils"

	_ "github.com/ChangerzaryX1602/sysware-test/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// SetupRoutes is the Router for GoFiber App
func (s *Server) SetupRoutes(app *fiber.App) {

	// Prepare a static middleware to serve the built React files.
	app.Static("/", "./web/build")

	app.Get("/swagger/*", swagger.HandlerDefault)

	// API routes group
	groupApiV1 := app.Group("/api/v:version?", handlers.ApiLimiter)
	{
		groupApiV1.Get("/", handlers.Index())
	}

	// App Repository
	userRepository := user.NewUserRepository(s.Resources)
	authRepository := auth.NewAuthRepository(s.Resources)
	permissionRepository := permission.NewPermissionRepository(s.Resources)
	roleRepository := role.NewRoleRepository(s.Resources)
	rolePermissionRepository := role_permission.NewRolePermissionRepository(s.Resources)
	userRoleRepository := user_role.NewUserRoleRepository(s.Resources)
	todoRepository := todo.NewTodoRepository(s.Resources)
	// auto migrate DB only on main process
	if !fiber.IsChild() {
		if err := userRepository.Migrate(); err != nil {
			log.Panicf("Failed to migrate user repository: %v", err)
		}
		if err := models.MigrateTablePermissions(s.MainDbConn); err != nil {
			log.Panicf("Failed to migrate permissions: %v", err)
		}
		if err := utils.MigratePermission(s.MainDbConn); err != nil {
			log.Panicf("Failed to migrate permission data: %v", err)
		}
		if err := todoRepository.Migrate(); err != nil {
			log.Panicf("Failed to migrate todo repository: %v", err)
		}
	}
	routerResources := handlers.NewRouterResources(s.JwtResources.JwtKeyfunc, s.MainDbConn)
	// App Services
	userService := user.NewUserService(userRepository)
	authService := auth.NewAuthService(authRepository, userService)
	permissionService := permission.NewPermissionService(permissionRepository)
	roleService := role.NewRoleService(roleRepository)
	rolePermissionService := role_permission.NewRolePermissionService(rolePermissionRepository)
	userRoleService := user_role.NewUserRoleService(userRoleRepository)
	todoService := todo.NewTodoService(todoRepository)

	// App Handlers
	user.NewUserHandler(groupApiV1.Group("/users"), routerResources, userService)
	auth.NewAuthHandler(groupApiV1.Group("/auth"), routerResources, authService)
	permission.NewPermissionHandler(groupApiV1.Group("/permissions"), routerResources, permissionService)
	role.NewRoleHandler(groupApiV1.Group("/roles"), routerResources, roleService)
	role_permission.NewRolePermissionHandler(groupApiV1.Group("/role/permissions"), routerResources, rolePermissionService)
	user_role.NewUserRoleHandler(groupApiV1.Group("/user/roles"), routerResources, userRoleService)
	todo.NewTodoHandler(groupApiV1.Group("/todos"), routerResources, todoService)
	// App Routes

	// Prepare a fallback route to always serve the 'index.html', had there not be any matching routes.
	app.Static("*", "./web/build/index.html")
}
