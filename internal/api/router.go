package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(
	router *gin.Engine,
	authHandler *AuthHandler,
	projectHandler *ProjectHandler,
	testSuiteHandler *TestSuiteHandler,
	testCaseHandler *TestCaseHandler,
	tagHandler *TagHandler,
) {
	// Public routes
	public := router.Group("/api/v1")
	{
		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Temporarily make these endpoints public for testing
		testSuites := public.Group("/test-suites")
		{
			testSuites.GET("", testSuiteHandler.ListTestSuites)
			testSuites.GET("/:id", testSuiteHandler.GetTestSuite)
		}

		tags := public.Group("/tags")
		{
			tags.GET("", tagHandler.ListTags)
			tags.GET("/:id", tagHandler.GetTag)
		}

		testCases := public.Group("/test-cases")
		{
			testCases.GET("/:id", testCaseHandler.GetTestCase)
		}
	}

	// Protected routes
	protected := router.Group("/api/v1")
	protected.Use(authHandler.AuthMiddleware())
	{
		// Projects
		projects := protected.Group("/projects")
		{
			projects.GET("", projectHandler.ListProjects)
			projects.POST("", projectHandler.CreateProject)
			projects.GET("/:id", projectHandler.GetProject)
			projects.PUT("/:id", projectHandler.UpdateProject)
			projects.DELETE("/:id", projectHandler.DeleteProject)
		}

		// Project test suites
		protected.GET("/project-test-suites/:projectId", testSuiteHandler.ListTestSuitesByProject)

		// Project test cases
		protected.GET("/project-test-cases/:projectId", testCaseHandler.ListTestCasesByProject)

		// Test Suites (protected operations)
		testSuitesProtected := protected.Group("/test-suites")
		{
			testSuitesProtected.POST("", testSuiteHandler.CreateTestSuite)
			testSuitesProtected.PUT("/:id", testSuiteHandler.UpdateTestSuite)
			testSuitesProtected.DELETE("/:id", testSuiteHandler.DeleteTestSuite)
		}

		// Suite test cases
		protected.GET("/suite-test-cases/:suiteId", testCaseHandler.ListTestCasesBySuite)

		// Test Cases (protected operations)
		testCasesProtected := protected.Group("/test-cases")
		{
			testCasesProtected.POST("", testCaseHandler.CreateTestCase)
			testCasesProtected.PUT("/:id", testCaseHandler.UpdateTestCase)
			testCasesProtected.DELETE("/:id", testCaseHandler.DeleteTestCase)
		}

		// Test case steps
		protected.POST("/test-case-steps/:testCaseId", testCaseHandler.AddTestStep)

		// Test Steps
		protected.PUT("/test-steps/:stepId", testCaseHandler.UpdateTestStep)
		protected.DELETE("/test-steps/:stepId", testCaseHandler.DeleteTestStep)

		// Step notes
		protected.POST("/step-notes/:stepId", testCaseHandler.AddStepNote)
		protected.DELETE("/step-notes/:noteId", testCaseHandler.DeleteStepNote)

		// Step attachments
		protected.POST("/step-attachments/:stepId", testCaseHandler.UploadStepAttachment)
		protected.DELETE("/step-attachments/:attachmentId", testCaseHandler.DeleteStepAttachment)

		// Tags (protected operations)
		tagsProtected := protected.Group("/tags")
		{
			tagsProtected.POST("", tagHandler.CreateTag)
			tagsProtected.DELETE("/:id", tagHandler.DeleteTag)
		}
	}
}
