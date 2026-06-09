package wire

import (
	"flowforge-api/handler"
	"flowforge-api/handler/middleware"
	infraClerk "flowforge-api/infrastructure/clerk"
	"flowforge-api/infrastructure/config"
	"flowforge-api/infrastructure/messaging/rabbitmq"
	repoGorm "flowforge-api/repository/gorm"
	"flowforge-api/usecase/auth"
	"flowforge-api/usecase/clerk"
	"flowforge-api/usecase/connexion"
	"flowforge-api/usecase/endpoint"
	"flowforge-api/usecase/organization"
	"flowforge-api/usecase/step"
	"flowforge-api/usecase/step_run"
	"flowforge-api/usecase/user"
	"flowforge-api/usecase/workflow"
	"flowforge-api/usecase/workflow_run"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	ClerkMiddleware        *middleware.ClerkMiddleware
	ClerkHandler           *handler.ClerkHandler
	UserHandler            *handler.UserHandler
	OrganizationHandler    *handler.OrganizationHandler
	EndpointHandler        *handler.EndpointHandler
	ConnexionHandler       *handler.ConnexionHandler
	StepHandler            *handler.StepHandler
	WorkflowHandler        *handler.WorkflowHandler
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	jwksProvider, err := infraClerk.NewJWKSProvider(env)
	if err != nil {
		log.Fatalf("failed to create JWKS provider: %v", err)
	}

	stepRunPublisher, err := rabbitmq.NewPublisherFromEnv(env)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ publisher: %v", err)
	}

	userRepo := repoGorm.NewUserRepository(db)
	organizationRepo := repoGorm.NewOrganizationRepository(db)
	endpointRepo := repoGorm.NewEndpointRepository(db)
	connexionRepo := repoGorm.NewConnexionRepository(db)
	stepRepo := repoGorm.NewStepRepository(db)
	workflowRepo := repoGorm.NewWorkflowRepository(db)
	workflowRunRepo := repoGorm.NewWorkflowRunRepository(db)
	stepRunRepo := repoGorm.NewStepRunRepository(db)

	validateTokenUseCase := auth.NewValidateTokenUseCase(jwksProvider, &userRepo)
	fetchUserUseCase := clerk.NewFetchUserUseCase(env)
	getUserByClerkIDUseCase := user.NewGetUserByClerkIDUseCase(&userRepo)
	createUserUseCase := user.NewCreateUserUseCase(&userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(&userRepo)
	deleteUserByClerkIDUseCase := user.NewDeleteUserByClerkIDUseCase(&userRepo)

	createOrganizationUseCase := organization.NewCreateOrganizationUseCase(&organizationRepo)
	listOrganizationsUseCase := organization.NewListOrganizationsUseCase(&organizationRepo)
	getOrganizationByIDUseCase := organization.NewGetOrganizationByIDUseCase(&organizationRepo)
	updateOrganizationUseCase := organization.NewUpdateOrganizationUseCase(&organizationRepo)
	activateOrganizationUseCase := organization.NewActivateOrganizationUseCase(&organizationRepo)

	listEndpointsUseCase := endpoint.NewListEndpointsUseCase(&endpointRepo)
	createEndpointUseCase := endpoint.NewCreateEndpointUseCase(&endpointRepo)
	updateEndpointUseCase := endpoint.NewUpdateEndpointUseCase(&endpointRepo)
	getEndpointUseCase := endpoint.NewGetEndpointUseCase(&endpointRepo)
	importFromOpenAPIUseCase := endpoint.NewImportFromOpenAPIUseCase(&endpointRepo, createEndpointUseCase)
	endpointHasStepUseCase := endpoint.NewEndpointHasStepUseCase(&stepRepo)
	deleteEndpointUseCase := endpoint.NewDeleteEndpointUseCase(&endpointRepo)
	createStepUseCase := step.NewCreateStepUseCase(&stepRepo)
	assignTreeIndicesUseCase := step.NewAssignTreeIndicesUseCase(&stepRepo, &connexionRepo)

	createConnexionUseCase := connexion.NewCreateConnexionUseCase(&connexionRepo, assignTreeIndicesUseCase)
	deleteConnexionUseCase := connexion.NewDeleteConnexionUseCase(&connexionRepo, assignTreeIndicesUseCase)

	listWorkflowsUseCase := workflow.NewListWorkflowsUseCase(&workflowRepo)
	createWorkflowUseCase := workflow.NewCreateWorkflowUseCase(&workflowRepo)
	calculateExecutionOrderUseCase := step.NewCalculateExecutionOrderUseCase()
	getStepUseCase := step.NewGetStepUseCase(&stepRepo)
	updateStepUseCase := step.NewUpdateStepUseCase(&stepRepo)
	deleteStepUseCase := step.NewDeleteStepUseCase(&stepRepo, &connexionRepo, &workflowRepo)
	getWorkflowUseCase := workflow.NewGetWorkflowUseCase(&workflowRepo)
	updateWorkflowUseCase := workflow.NewUpdateWorkflowUseCase(&workflowRepo)
	activateWorkflowUseCase := workflow.NewActivateWorkflowUseCase(&workflowRepo)
	deactivateWorkflowUseCase := workflow.NewDeactivateWorkflowUseCase(&workflowRepo)
	getWorkflowAnalyticsUseCase := workflow_run.NewGetWorkflowAnalyticsUseCase(
		&workflowRepo,
		&workflowRunRepo,
	)
	upsertWorkflowUseCase := workflow.NewUpsertWorkflowUseCase(
		&workflowRepo,
		&stepRepo,
		&endpointRepo,
		&connexionRepo,
		calculateExecutionOrderUseCase,
		createStepUseCase,
		assignTreeIndicesUseCase,
	)
	getWorkflowRunsUseCase := workflow_run.NewGetWorkflowRunsUseCase(
		&workflowRepo,
		&workflowRunRepo,
	)

	createWorkflowRunUseCase := workflow_run.NewCreateWorkflowRunUseCase(&workflowRunRepo)
	hasStepRunUseCase := step_run.NewHasStepRunUseCase(&stepRunRepo)
	createStepRunUseCase := step_run.NewCreateStepRunUseCase(&stepRunRepo, &stepRepo)
	executeStepRunUseCase := step_run.NewExecuteStepRunUseCase(&stepRunRepo, &stepRepo)
	executeWorkflowRunUseCase := workflow_run.NewExecuteWorkflowRunUseCase(&workflowRunRepo)

	runWorkflowUseCase := workflow.NewRunWorkflowUseCase(&workflowRepo, &workflowRunRepo, &stepRepo, createWorkflowRunUseCase, hasStepRunUseCase, createStepRunUseCase, executeStepRunUseCase, executeWorkflowRunUseCase, env, stepRunPublisher)

	startWorkflowUseCase := workflow.NewStartWorkflowUseCase(&workflowRepo, &workflowRunRepo, runWorkflowUseCase)
	stopWorkflowUseCase := workflow.NewStopWorkflowUseCase(&workflowRepo, &workflowRunRepo, &stepRunRepo, runWorkflowUseCase)

	clerkMiddleware := middleware.NewClerkMiddleware(env.ClerkWebhookSecret)
	authenticateMiddleware := middleware.NewAuthenticateMiddleware(
		validateTokenUseCase,
		fetchUserUseCase,
		createUserUseCase,
		createOrganizationUseCase,
		updateUserUseCase,
	)

	return &Container{
		AuthenticateMiddleware: authenticateMiddleware,
		ClerkMiddleware:        clerkMiddleware,
		ClerkHandler: handler.NewClerkHandler(
			getUserByClerkIDUseCase,
			createUserUseCase,
			createOrganizationUseCase,
			updateUserUseCase,
			deleteUserByClerkIDUseCase,
		),
		UserHandler: handler.NewUserHandler(),
		OrganizationHandler: handler.NewOrganizationHandler(
			listOrganizationsUseCase,
			createOrganizationUseCase,
			getOrganizationByIDUseCase,
			updateOrganizationUseCase,
			activateOrganizationUseCase,
		),
		EndpointHandler: handler.NewEndpointHandler(
			listEndpointsUseCase,
			createEndpointUseCase,
			updateEndpointUseCase,
			getEndpointUseCase,
			importFromOpenAPIUseCase,
			endpointHasStepUseCase,
			deleteEndpointUseCase,
		),
		ConnexionHandler: handler.NewConnexionHandler(
			createConnexionUseCase,
			deleteConnexionUseCase,
		),
		StepHandler: handler.NewStepHandler(
			getStepUseCase,
			updateStepUseCase,
			deleteStepUseCase,
		),
		WorkflowHandler: handler.NewWorkflowHandler(
			listWorkflowsUseCase,
			createWorkflowUseCase,
			getWorkflowUseCase,
			updateWorkflowUseCase,
			activateWorkflowUseCase,
			deactivateWorkflowUseCase,
			upsertWorkflowUseCase,
			getWorkflowRunsUseCase,
			getWorkflowAnalyticsUseCase,
			startWorkflowUseCase,
			stopWorkflowUseCase,
		),
	}
}
