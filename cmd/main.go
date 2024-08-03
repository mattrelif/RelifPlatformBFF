package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"relif/bff/clients"
	"relif/bff/http"
	"relif/bff/http/handlers"
	"relif/bff/http/middlewares"
	"relif/bff/repositories"
	"relif/bff/services"
	"relif/bff/settings"
	"relif/bff/utils"
	"syscall"
)

var (
	logger      *zap.Logger
	environment *settings.Environment
	awsConfig   aws.Config
)

func init() {
	var err error

	logger, err = zap.NewProduction()

	if err != nil {
		panic(err)
	}

	environment, err = settings.NewEnvironment()

	if err != nil {
		logger.Fatal("could not initialize environment settings", zap.Error(err))
	}

	awsConfig, err = settings.NewAWSConfig(environment.AWS.Region)

	if err != nil {
		logger.Fatal("could not initialize AWS config", zap.Error(err))
	}
}

func main() {
	defer logger.Sync()

	mongo, err := clients.NewMongoClient(environment.Mongo.URI, environment.Mongo.ConnectionTimeout)

	if err != nil {
		logger.Fatal("could not initialize mongo client", zap.Error(err))
	}

	database := mongo.Database(environment.Mongo.Database)

	sesClient := clients.NewSESClient(awsConfig)

	sessionsRepository := repositories.NewSessionsMongo(database)
	usersRepository := repositories.NewUsersMongo(database)
	passwordChangeRequestsRepository := repositories.NewMongoPasswordChangeRequests(database)
	organizationsRepository := repositories.NewMongoOrganizations(database)
	joinOrganizationRequestsRepository := repositories.NewMongoJoinOrganizationRequests(database)
	joinOrganizationInvitesRepository := repositories.NewMongoJoinOrganizationInvites(database)
	housingsRepository := repositories.NewMongoHousings(database)
	updateOrganizationTypeRequestsRepository := repositories.NewMongoUpdateOrganizationTypeRequests(database)
	organizationDataAccessGrantsRepository := repositories.NewMongoOrganizationDataAccessGrants(database)
	organizationDataAccessRepository := repositories.NewMongoOrganizationDataAccessRequests(database)
	joinPlatformInvitesRepository := repositories.NewMongoJoinPlatformInvites(database)
	beneficiariesRepository := repositories.NewMongoBeneficiaries(database)
	housingRoomsRepository := repositories.NewMongoHousingRooms(database)
	beneficiaryAllocationsRepository := repositories.NewMongoBeneficiaryAllocations(database)
	voluntaryPeopleRepository := repositories.NewMongoVoluntaryPeople(database)
	productTypesRepository := repositories.NewMongoProductTypesRepository(database)

	sesEmailService := services.NewSesEmail(sesClient)

	usersService := services.NewUsers(usersRepository)
	sessionsService := services.NewSessions(sessionsRepository, utils.GenerateUuid)
	organizationsService := services.NewOrganizations(organizationsRepository, usersService)
	passwordService := services.NewPassword(sesEmailService, usersService, passwordChangeRequestsRepository, utils.BcryptHash, utils.GenerateUuid)
	authService := services.NewAuth(usersService, sessionsService, utils.BcryptHash, utils.BcryptCompare)
	joinOrganizationRequestsService := services.NewJoinOrganizationRequests(usersService, joinOrganizationRequestsRepository)
	joinOrganizationInvitesService := services.NewJoinOrganizationInvites(usersService, joinOrganizationInvitesRepository)
	updateOrganizationTypeRequestsService := services.NewUpdateOrganizationTypeRequests(organizationsService, updateOrganizationTypeRequestsRepository)
	organizationsDataAccessGrantsService := services.NewOrganizationDataAccessGrants(organizationDataAccessGrantsRepository)
	housingsService := services.NewHousings(housingsRepository, organizationsDataAccessGrantsService)
	organizationDataAccessService := services.NewOrganizationDataAccessRequests(organizationDataAccessRepository, organizationsService, organizationsDataAccessGrantsService)
	joinPlatformInvitesService := services.NewJoinPlatformInvites(joinPlatformInvitesRepository, sesEmailService, usersService, utils.GenerateUuid)
	housingRoomsService := services.NewHousingRooms(housingRoomsRepository, housingsService)
	beneficiariesService := services.NewBeneficiaries(beneficiariesRepository, housingsService, housingRoomsService, organizationsDataAccessGrantsService)
	beneficiaryAllocationsService := services.NewBeneficiaryAllocations(beneficiaryAllocationsRepository, beneficiariesService, housingRoomsService)
	voluntaryPeopleService := services.NewVoluntaryPeople(voluntaryPeopleRepository)
	productTypesService := services.NewProductTypes(productTypesRepository)

	authenticateByCookieMiddleware := middlewares.NewAuthenticateByCookie(authService)

	authHandler := handlers.NewAuth(authService)
	passwordHandler := handlers.NewPassword(passwordService)
	usersHandler := handlers.NewUsers(usersService)
	organizationsHandler := handlers.NewOrganizations(organizationsService)
	joinOrganizationRequestsHandler := handlers.NewJoinOrganizationRequests(joinOrganizationRequestsService)
	joinOrganizationInvitesHandler := handlers.NewJoinOrganizationInvites(joinOrganizationInvitesService)
	housingsHandler := handlers.NewHousings(housingsService)
	updateOrganizationTypeRequestsHandler := handlers.NewUpdateOrganizationTypeRequests(updateOrganizationTypeRequestsService)
	organizationsDataAccessRequestsHandler := handlers.NewOrganizationDataAccessRequests(organizationDataAccessService)
	joinPlatformInvitesHandler := handlers.NewJoinPlatformInvites(joinPlatformInvitesService)
	beneficiariesHandler := handlers.NewBeneficiaries(beneficiariesService)
	housingRoomsHandler := handlers.NewHousingRooms(housingRoomsService)
	beneficiaryAllocationsHandler := handlers.NewBeneficiaryAllocations(beneficiaryAllocationsService)
	voluntaryPeopleHandler := handlers.NewVoluntaryPeople(voluntaryPeopleService)
	productTypesHandler := handlers.NewProductTypes(productTypesService)
	organizationsDataAccessGrantsHandler := handlers.NewOrganizationDataAccessGrants(organizationsDataAccessGrantsService)

	router := http.NewRouter(
		environment.Environment,
		environment.Router.Context,
		authenticateByCookieMiddleware,
		authHandler,
		beneficiariesHandler,
		beneficiaryAllocationsHandler,
		housingsHandler,
		housingRoomsHandler,
		joinOrganizationRequestsHandler,
		joinOrganizationInvitesHandler,
		joinPlatformInvitesHandler,
		organizationsHandler,
		organizationsDataAccessGrantsHandler,
		organizationsDataAccessRequestsHandler,
		passwordHandler,
		productTypesHandler,
		updateOrganizationTypeRequestsHandler,
		usersHandler,
		voluntaryPeopleHandler,
	)
	server := http.NewServer(router, environment.Server.Port, environment.Server.ReadTimeout, environment.Server.WriteTimeout)

	go func() {
		logger.Info("starting server", zap.Int("port", environment.Server.Port))
		if err = server.Start(); err != nil {
			logger.Fatal("could not start server", zap.Error(err))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	if err = clients.DisconnectMongoClient(mongo); err != nil {
		logger.Fatal("could not disconnect mongo client", zap.Error(err))
	}

	if err = server.Stop(); err != nil {
		logger.Fatal("could not stop server", zap.Error(err))
	}
}
