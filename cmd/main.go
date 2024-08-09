package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"relif/platform-bff/clients"
	"relif/platform-bff/http"
	"relif/platform-bff/http/handlers"
	"relif/platform-bff/http/middlewares"
	"relif/platform-bff/repositories"
	"relif/platform-bff/services"
	"relif/platform-bff/settings"
	"relif/platform-bff/utils"
	"syscall"
)

var (
	logger    *zap.Logger
	stgs      *settings.Settings
	awsConfig aws.Config
)

func init() {
	var err error

	logger, err = zap.NewProduction()

	if err != nil {
		panic(err)
	}

	logger.Info("initializing AWS configuration")
	awsConfig, err = settings.NewAWSConfig(settings.AWSRegion)

	if err != nil {
		logger.Fatal("could not initialize AWS config", zap.Error(err))
	}

	logger.Info("initializing Secrets Manager client")
	secretsManagerClient := clients.NewSecretsManager(awsConfig)

	logger.Info("initializing settings")
	stgs, err = settings.NewSettings(secretsManagerClient)

	if err != nil {
		logger.Fatal("could not initialize settings", zap.Error(err))
	}

	logger.Info("settings initialized", zap.String("settings", fmt.Sprintf("%+v", stgs)))
}

func main() {
	defer logger.Sync()

	mongo, err := clients.NewMongoClient(stgs.Mongo.URI, stgs.Mongo.ConnectionTimeout)

	if err != nil {
		logger.Fatal("could not initialize mongo client", zap.Error(err))
	}

	database := mongo.Database(stgs.Mongo.Database)

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

	sesEmailService := services.NewSesEmail(sesClient, stgs.Email.Domain)

	usersService := services.NewUsers(usersRepository)
	sessionsService := services.NewSessions(sessionsRepository, utils.GenerateUuid)
	organizationsService := services.NewOrganizations(organizationsRepository, usersService)
	passwordService := services.NewPassword(sesEmailService, usersService, passwordChangeRequestsRepository, utils.BcryptHash, utils.GenerateUuid)
	authService := services.NewAuth(usersService, sessionsService, organizationsService, utils.BcryptHash, utils.BcryptCompare)
	joinOrganizationRequestsService := services.NewJoinOrganizationRequests(usersService, joinOrganizationRequestsRepository)
	joinOrganizationInvitesService := services.NewJoinOrganizationInvites(usersService, joinOrganizationInvitesRepository)
	updateOrganizationTypeRequestsService := services.NewUpdateOrganizationTypeRequests(organizationsService, updateOrganizationTypeRequestsRepository)
	organizationsDataAccessGrantsService := services.NewOrganizationDataAccessGrants(organizationDataAccessGrantsRepository)
	housingsService := services.NewHousings(housingsRepository)
	organizationDataAccessService := services.NewOrganizationDataAccessRequests(organizationDataAccessRepository, organizationsService, organizationsDataAccessGrantsService)
	joinPlatformInvitesService := services.NewJoinPlatformInvites(joinPlatformInvitesRepository, sesEmailService, usersService, utils.GenerateUuid)
	housingRoomsService := services.NewHousingRooms(housingRoomsRepository)
	beneficiariesService := services.NewBeneficiaries(beneficiariesRepository)
	beneficiaryAllocationsService := services.NewBeneficiaryAllocations(beneficiaryAllocationsRepository, beneficiariesService, housingRoomsService, housingsService)
	voluntaryPeopleService := services.NewVoluntaryPeople(voluntaryPeopleRepository)
	productTypesService := services.NewProductTypes(productTypesRepository)

	authorizationService := services.NewAuthorization(
		usersService,
		organizationsService,
		housingsService,
		housingRoomsService,
		beneficiariesService,
		organizationDataAccessService,
		organizationsDataAccessGrantsService,
		joinOrganizationInvitesService,
		updateOrganizationTypeRequestsService,
		joinOrganizationRequestsService,
		voluntaryPeopleService,
		productTypesService,
	)

	authenticateByCookieMiddleware := middlewares.NewAuthenticateByCookie(authService)

	authHandler := handlers.NewAuth(authService)
	passwordHandler := handlers.NewPassword(passwordService)
	usersHandler := handlers.NewUsers(usersService, authorizationService)
	organizationsHandler := handlers.NewOrganizations(organizationsService, authorizationService)
	joinOrganizationRequestsHandler := handlers.NewJoinOrganizationRequests(joinOrganizationRequestsService, authorizationService)
	joinOrganizationInvitesHandler := handlers.NewJoinOrganizationInvites(joinOrganizationInvitesService, authorizationService)
	housingsHandler := handlers.NewHousings(housingsService, authorizationService)
	updateOrganizationTypeRequestsHandler := handlers.NewUpdateOrganizationTypeRequests(updateOrganizationTypeRequestsService, authorizationService)
	organizationsDataAccessRequestsHandler := handlers.NewOrganizationDataAccessRequests(organizationDataAccessService, authorizationService)
	joinPlatformInvitesHandler := handlers.NewJoinPlatformInvites(joinPlatformInvitesService, authorizationService)
	beneficiariesHandler := handlers.NewBeneficiaries(beneficiariesService, authorizationService)
	housingRoomsHandler := handlers.NewHousingRooms(housingRoomsService, authorizationService)
	beneficiaryAllocationsHandler := handlers.NewBeneficiaryAllocations(beneficiaryAllocationsService, authorizationService)
	voluntaryPeopleHandler := handlers.NewVoluntaryPeople(voluntaryPeopleService, authorizationService)
	productTypesHandler := handlers.NewProductTypes(productTypesService, authorizationService)
	organizationsDataAccessGrantsHandler := handlers.NewOrganizationDataAccessGrants(organizationsDataAccessGrantsService, authorizationService)

	healthHandler := handlers.NewHealth()

	router := http.NewRouter(
		stgs,
		authenticateByCookieMiddleware,
		healthHandler,
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
	server := http.NewServer(router, stgs.Server.Port, stgs.Server.ReadTimeout, stgs.Server.WriteTimeout)

	go func() {
		logger.Info("starting server", zap.Int("port", stgs.Server.Port))
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
