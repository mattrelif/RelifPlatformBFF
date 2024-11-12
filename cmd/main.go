package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/joho/godotenv"
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
	authenticationUseCases "relif/platform-bff/usecases/authentication"
	beneficiariesUseCases "relif/platform-bff/usecases/beneficiaries"
	beneficiaryAllocationsUseCases "relif/platform-bff/usecases/beneficiary_allocations"
	donationsUseCases "relif/platform-bff/usecases/donations"
	filesUseCases "relif/platform-bff/usecases/files"
	housingRoomsUseCases "relif/platform-bff/usecases/housing_rooms"
	housingsUseCases "relif/platform-bff/usecases/housings"
	joinOrganizationInvitesUseCases "relif/platform-bff/usecases/join_organization_invites"
	joinOrganizationRequestsUseCases "relif/platform-bff/usecases/join_organization_requests"
	joinPlatformAdminInvitesUseCases "relif/platform-bff/usecases/join_platform_admin_invites"
	joinPlatformInvitesUseCases "relif/platform-bff/usecases/join_platform_invites"
	organizationDataAccessGrantsUseCases "relif/platform-bff/usecases/organization_data_access_grants"
	organizationDataAccessRequestsUseCases "relif/platform-bff/usecases/organization_data_access_requests"
	organizationsUseCases "relif/platform-bff/usecases/organizations"
	passwordRecoveryUseCases "relif/platform-bff/usecases/password_recovery"
	productTypeAllocationsUseCases "relif/platform-bff/usecases/product_type_alloctions"
	productTypesUseCases "relif/platform-bff/usecases/product_types"
	storageRecordsUseCases "relif/platform-bff/usecases/storage_records"
	updateOrganizationTypeRequestsUseCases "relif/platform-bff/usecases/update_organization_type_requets"
	usersUseCases "relif/platform-bff/usecases/users"
	voluntaryPeopleUseCases "relif/platform-bff/usecases/voluntary_people"
	"relif/platform-bff/utils"
	"syscall"
)

var (
	logger           *zap.Logger
	settingsInstance *settings.Settings
	awsConfig        aws.Config
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
	settingsInstance, err = settings.NewSettings(secretsManagerClient)

	if err != nil {
		logger.Fatal("could not initialize settings", zap.Error(err))
	}
}

func main() {
	defer logger.Sync()

	mongo, err := clients.NewMongoClient(settingsInstance.MongoURI)

	if err != nil {
		logger.Fatal("could not initialize mongo client", zap.Error(err))
	}

	database := mongo.Database(settingsInstance.MongoDatabase)

	sesClient := clients.NewSESClient(awsConfig)
	s3Client := clients.NewS3(awsConfig)

	/** Repositories **/
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
	productTypesAllocationRepository := repositories.NewMongoProductTypeAllocations(database)
	storageRecordsRepository := repositories.NewMongoStorageRecords(database)
	donationsRepository := repositories.NewDonations(database)
	joinPlatformAdminInvitesRepository := repositories.NewJoinPlatformAdminInvites(database)

	/** Services **/
	sesEmailService := services.NewSesEmail(sesClient, settingsInstance.EmailDomain, settingsInstance.FrontendDomain)
	tokensService := services.NewTokens([]byte(settingsInstance.TokenSecret))
	s3FileUploadsService := services.NewS3FileUploads(s3Client, settingsInstance.S3BucketName)
	cognitoService, err := services.NewCognito(settingsInstance.AWS_REGION, settingsInstance.COGNITO_CLIENT_ID, usersRepository)

	/** Use Cases **/
	generateUploadLinkUseCase := filesUseCases.NewGenerateUploadLink(s3FileUploadsService, utils.GenerateUuid)

	authenticateTokenUseCase := authenticationUseCases.NewAuthenticateToken(usersRepository, sessionsRepository, tokensService)

	createUserUseCase := usersUseCases.NewCreate(usersRepository)

	signUpUseCase := authenticationUseCases.NewSignUp(sessionsRepository, tokensService, createUserUseCase, utils.BcryptHash, cognitoService)
	adminSignUpUseCase := authenticationUseCases.NewAdminSignUp(sessionsRepository, tokensService, createUserUseCase, utils.BcryptHash)
	organizationSignUpUseCase := authenticationUseCases.NewOrganizationSignUp(sessionsRepository, organizationsRepository, tokensService, createUserUseCase, utils.BcryptHash)
	signInUseCase := authenticationUseCases.NewSignIn(usersRepository, sessionsRepository, tokensService, utils.BcryptCompare, cognitoService)
	signOutUseCase := authenticationUseCases.NewSignOut(sessionsRepository)
	verifyUseCase := authenticationUseCases.NewVerify(cognitoService)

	requestPasswordChangeUseCase := passwordRecoveryUseCases.NewRequestChange(usersRepository, passwordChangeRequestsRepository, sesEmailService, utils.GenerateUuid)
	changePasswordUseCase := passwordRecoveryUseCases.NewChange(usersRepository, passwordChangeRequestsRepository, sesEmailService, utils.BcryptHash)

	findOneUserCompleteUseCase := usersUseCases.NewFindOneCompleteByID(usersRepository)
	findManyUsersByOrganizationIDUseCase := usersUseCases.NewFindManyByOrganizationIDPaginated(usersRepository, organizationsRepository)
	findManyRelifMembersUseCase := usersUseCases.NewFindManyRelifMembersPaginated(usersRepository)
	updateOneUserByIDUseCase := usersUseCases.NewUpdateOneByID(usersRepository)
	inactivateOneUserByIDUseCase := usersUseCases.NewInactivateOneByID(usersRepository)
	reactivateOneUseByIDUseCase := usersUseCases.NewReactivateOneByID(usersRepository)

	createOrganizationUseCase := organizationsUseCases.NewCreate(organizationsRepository, usersRepository)
	findManyOrganizationsUseCase := organizationsUseCases.NewFindManyPaginated(organizationsRepository)
	findOneOrganizationByIDUseCase := organizationsUseCases.NewFindOneByID(organizationsRepository)
	updateOneOrganizationByIDUseCase := organizationsUseCases.NewUpdateOneByID(organizationsRepository)
	inactivateOneOrganizationByIDUseCase := organizationsUseCases.NewInactivateOneByID(organizationsRepository)
	reactivateOneOrganizationByIDUseCase := organizationsUseCases.NewReactivateOneByID(organizationsRepository)

	createJoinOrganizationRequestUseCase := joinOrganizationRequestsUseCases.NewCreate(joinOrganizationRequestsRepository, organizationsRepository)
	findManyJoinOrganizationRequestsByOrganizationIDUseCase := joinOrganizationRequestsUseCases.NewFindManyByOrganizationIDPaginated(joinOrganizationRequestsRepository, organizationsRepository)
	findManyJoinOrganizationRequestsByUserIDUseCase := joinOrganizationRequestsUseCases.NewFindManyByUserIDPaginated(joinOrganizationRequestsRepository, usersRepository)
	acceptJoinOrganizationRequestUseCase := joinOrganizationRequestsUseCases.NewAccept(joinOrganizationRequestsRepository, usersRepository, organizationsRepository)
	rejectJoinOrganizationRequestUseCase := joinOrganizationRequestsUseCases.NewReject(joinOrganizationRequestsRepository, organizationsRepository)

	createJoinOrganizationInviteUseCase := joinOrganizationInvitesUseCases.NewCreate(joinOrganizationInvitesRepository, organizationsRepository)
	findManyJoinOrganizationInvitesByOrganizationIDUseCase := joinOrganizationInvitesUseCases.NewFindManyByOrganizationIDPaginated(joinOrganizationInvitesRepository, organizationsRepository)
	findManyJoinOrganizationInvitesByUserIDUseCase := joinOrganizationInvitesUseCases.NewFindManyByUserIDPaginated(joinOrganizationInvitesRepository, usersRepository)
	acceptJoinOrganizationInviteUseCase := joinOrganizationInvitesUseCases.NewAccept(joinOrganizationInvitesRepository, usersRepository)
	rejectJoinOrganizationInviteUseCase := joinOrganizationInvitesUseCases.NewReject(joinOrganizationInvitesRepository, usersRepository)

	createHousingUseCase := housingsUseCases.NewCreate(housingsRepository)
	findManyHousingsByOrganizationIDUseCase := housingsUseCases.NewFindManyByOrganizationIDPaginated(housingsRepository, organizationsRepository)
	findOneHousingCompleteByIDUseCase := housingsUseCases.NewFindOneCompleteByID(housingsRepository, organizationsRepository)
	updateOneHousingByIDUseCase := housingsUseCases.NewUpdateOneByID(housingsRepository, organizationsRepository)
	deleteOneHousingByIDUseCase := housingsUseCases.NewDeleteOneByID(housingsRepository, organizationsRepository)

	createUpdateOrganizationTypeRequestUseCase := updateOrganizationTypeRequestsUseCases.NewCreate(updateOrganizationTypeRequestsRepository, organizationsRepository)
	findManyUpdateOrganizationTypeRequestsUseCase := updateOrganizationTypeRequestsUseCases.NewFindManyPaginated(updateOrganizationTypeRequestsRepository)
	findManyUpdateOrganizationTypeRequestsByOrganizationIDUseCase := updateOrganizationTypeRequestsUseCases.NewFindManyByOrganizationIDPaginated(updateOrganizationTypeRequestsRepository, organizationsRepository)
	acceptUpdateOrganizationTypeRequestUseCase := updateOrganizationTypeRequestsUseCases.NewAccept(updateOrganizationTypeRequestsRepository, organizationsRepository)
	rejectUpdateOrganizationTypeRequestUseCase := updateOrganizationTypeRequestsUseCases.NewReject(updateOrganizationTypeRequestsRepository)

	createOrganizationDataAccessRequestsUseCase := organizationDataAccessRequestsUseCases.NewCreate(organizationDataAccessRepository, organizationsRepository)
	findManyOrganizationDataAccessRequestsByTargetUseCase := organizationDataAccessRequestsUseCases.NewFindManyByTargetOrganizationIDPaginated(organizationDataAccessRepository, organizationsRepository)
	findManyOrganizationDataAccessRequestsByRequesterUseCase := organizationDataAccessRequestsUseCases.NewFindManyByRequesterOrganizationIDPaginated(organizationDataAccessRepository, organizationsRepository)
	acceptOrganizationDataAccessRequestsUseCase := organizationDataAccessRequestsUseCases.NewAccept(organizationDataAccessRepository, organizationDataAccessGrantsRepository, organizationsRepository)
	rejectOrganizationDataAccessRequestsUseCase := organizationDataAccessRequestsUseCases.NewReject(organizationDataAccessRepository, organizationsRepository)

	createJoinPlatformInviteUseCase := joinPlatformInvitesUseCases.NewCreate(joinPlatformInvitesRepository, organizationsRepository, usersRepository, sesEmailService, utils.GenerateUuid)
	findManyJoinPlatformInvitesByOrganizationIDUseCase := joinPlatformInvitesUseCases.NewFindManyByOrganizationIDPaginated(joinPlatformInvitesRepository, organizationsRepository)
	consumeJoinPlatformInviteByCodeUseCase := joinPlatformInvitesUseCases.NewConsumeByCode(joinPlatformInvitesRepository)

	createBeneficiaryUseCase := beneficiariesUseCases.NewCreate(beneficiariesRepository, organizationsRepository)
	findManyBeneficiariesByOrganizationIDUseCase := beneficiariesUseCases.NewFindManyByOrganizationIDPaginated(beneficiariesRepository, organizationsRepository)
	findManyBeneficiariesByHousingIDUseCase := beneficiariesUseCases.NewFindManyByHousingIDPaginated(beneficiariesRepository, housingsRepository, organizationsRepository)
	findManyBeneficiariesByHousingRoomIDUseCase := beneficiariesUseCases.NewFindManyByHousingRoomIDPaginated(beneficiariesRepository, housingRoomsRepository, housingsRepository, organizationsRepository)
	findOneBeneficiaryCompleteByIDUseCase := beneficiariesUseCases.NewFindOneCompleteByID(beneficiariesRepository)
	updateBeneficiaryByIDUseCase := beneficiariesUseCases.NewUpdateOneByID(beneficiariesRepository, organizationsRepository)
	deleteBeneficiaryByIDUseCase := beneficiariesUseCases.NewDeleteOneByID(beneficiariesRepository, organizationsRepository)
	generateBeneficiaryProfileImageUploadLinkUseCase := beneficiariesUseCases.NewGenerateProfileImageUploadLink(generateUploadLinkUseCase)

	createHousingRoomUseCase := housingRoomsUseCases.NewCreateHousingRoom(housingRoomsRepository, housingsRepository, organizationsRepository)
	findOneHousingRoomCompleteByIDUseCase := housingRoomsUseCases.NewFindOneCompleteByID(housingRoomsRepository, organizationsRepository)
	findManyHousingRoomsByHousingIDUseCase := housingRoomsUseCases.NewFindManyByHousingIDPaginated(housingRoomsRepository, housingsRepository, organizationsRepository)
	updateHousingRoomByIDUseCase := housingRoomsUseCases.NewUpdateOneByID(housingRoomsRepository, housingsRepository, organizationsRepository)
	deleteHousingRoomByIDUseCase := housingRoomsUseCases.NewDeleteOneByID(housingRoomsRepository, housingsRepository, organizationsRepository)

	createEntranceBeneficiaryAllocationUseCase := beneficiaryAllocationsUseCases.NewCreateEntrance(beneficiaryAllocationsRepository, beneficiariesRepository, organizationsRepository, housingsRepository, housingRoomsRepository)
	createReallocationBeneficiaryAllocationUseCase := beneficiaryAllocationsUseCases.NewCreateReallocation(beneficiaryAllocationsRepository, beneficiariesRepository, organizationsRepository, housingsRepository, housingRoomsRepository)
	findManyBeneficiaryAllocationsByBeneficiaryIDUseCase := beneficiaryAllocationsUseCases.NewFindManyByBeneficiaryIDPaginated(beneficiaryAllocationsRepository, beneficiariesRepository, organizationsRepository)
	findManyBeneficiaryAllocationsByHousingIDUseCase := beneficiaryAllocationsUseCases.NewFindManyByHousingIDPaginated(beneficiaryAllocationsRepository, housingsRepository, organizationsRepository)
	findManyBeneficiaryAllocationsByHousingRoomIDUseCase := beneficiaryAllocationsUseCases.NewFindManyByHousingRoomIDPaginated(beneficiaryAllocationsRepository, housingsRepository, housingRoomsRepository, organizationsRepository)

	createVoluntaryPersonUseCase := voluntaryPeopleUseCases.NewCreate(voluntaryPeopleRepository, organizationsRepository)
	findManyVoluntaryPeopleByOrganizationIDUseCase := voluntaryPeopleUseCases.NewFindManyByOrganizationIDPaginated(voluntaryPeopleRepository, organizationsRepository)
	findOneVoluntaryPersonCompleteByIDUseCase := voluntaryPeopleUseCases.NewFindOneByID(voluntaryPeopleRepository, organizationsRepository)
	updateOneVoluntaryPersonByIDUseCase := voluntaryPeopleUseCases.NewUpdateOneByID(voluntaryPeopleRepository, organizationsRepository)
	deleteOneVoluntaryPersonByIDUseCase := voluntaryPeopleUseCases.NewDeleteOneByID(voluntaryPeopleRepository, organizationsRepository)

	createProductTypeUseCase := productTypesUseCases.NewCreate(organizationsRepository, productTypesRepository)
	findManyProductTypesByOrganizationIDUseCase := productTypesUseCases.NewFindManyByOrganizationIDPaginated(productTypesRepository, organizationsRepository)
	findOneProductTypeByIDUseCase := productTypesUseCases.NewFindOneCompleteByID(productTypesRepository)
	updateProductTypeByIDUseCase := productTypesUseCases.NewUpdateOneByID(productTypesRepository, organizationsRepository)
	deleteProductTypeByIDUseCase := productTypesUseCases.NewDeleteOneByID(productTypesRepository, organizationsRepository, donationsRepository, storageRecordsRepository, productTypesAllocationRepository)

	findManyGrantsByOrganizationIDUseCase := organizationDataAccessGrantsUseCases.NewFindManyByOrganizationIDPaginated(organizationDataAccessGrantsRepository, organizationsRepository)
	findManyGrantsByTargetOrganizationIDUseCase := organizationDataAccessGrantsUseCases.NewFindManyByTargetOrganizationIDPaginated(organizationDataAccessGrantsRepository, organizationsRepository)
	deleteGrantByIDUseCase := organizationDataAccessGrantsUseCases.NewDeleteOneByIDImpl(organizationDataAccessGrantsRepository, organizationsRepository)

	createProductTypeAllocationEntranceUseCase := productTypeAllocationsUseCases.NewCreateEntrance(productTypesAllocationRepository, productTypesRepository, organizationsRepository, storageRecordsRepository)
	createProductTypeAllocationReallocationUseCase := productTypeAllocationsUseCases.NewCreateReallocation(productTypesAllocationRepository, productTypesRepository, organizationsRepository, storageRecordsRepository)
	findManyProductTypeAllocationsByProductTypeIDPaginatedUseCase := productTypeAllocationsUseCases.NewFindManyByProductTypeIDPaginated(productTypesRepository, productTypesAllocationRepository, organizationsRepository)

	createDonationUseCase := donationsUseCases.NewCreate(donationsRepository, beneficiariesRepository, storageRecordsRepository, organizationsRepository, productTypesRepository)
	findManyDonationsByBeneficiaryIDUseCase := donationsUseCases.NewFindManyByBeneficiaryIDPaginated(donationsRepository, beneficiariesRepository, organizationsRepository)
	findManyDonationsByProductTypeIDUseCase := donationsUseCases.NewFindManyByProductTypeIDPaginated(donationsRepository, productTypesRepository, organizationsRepository)

	findManyStorageRecordsByOrganizationIDUseCase := storageRecordsUseCases.NewFindManyByOrganizationIDPaginated(storageRecordsRepository, organizationsRepository)
	findManyStorageRecordsByHousingIDUseCase := storageRecordsUseCases.NewFindManyByHousingIDPaginated(storageRecordsRepository, organizationsRepository, housingsRepository)
	findManyStorageRecordsByProductTypeIDUseCase := storageRecordsUseCases.NewFindManyByProductTypeID(productTypesRepository, organizationsRepository, storageRecordsRepository)

	createJoinPlatformAdminInvitesUseCase := joinPlatformAdminInvitesUseCases.NewCreate(joinPlatformAdminInvitesRepository, sesEmailService, utils.GenerateUuid)
	findManyJoinPlatformAdminInvitesPaginatedUseCase := joinPlatformAdminInvitesUseCases.NewFindManyPaginated(joinPlatformAdminInvitesRepository)
	consumeJoinPlatformAdminInviteByCodeUseCase := joinPlatformAdminInvitesUseCases.NewConsumeByCode(joinPlatformAdminInvitesRepository)

	/** Middlewares **/
	authenticateByCookieMiddleware := middlewares.NewAuthenticateByToken(authenticateTokenUseCase)

	/** Handlers **/
	authenticationHandler := handlers.NewAuthentication(signUpUseCase, organizationSignUpUseCase, adminSignUpUseCase, signInUseCase, signOutUseCase, verifyUseCase)
	passwordRecoveryHandler := handlers.NewPassword(requestPasswordChangeUseCase, changePasswordUseCase)
	usersHandler := handlers.NewUsers(findOneUserCompleteUseCase, findManyUsersByOrganizationIDUseCase, findManyRelifMembersUseCase, updateOneUserByIDUseCase, inactivateOneUserByIDUseCase, reactivateOneUseByIDUseCase)
	organizationsHandler := handlers.NewOrganizations(createOrganizationUseCase, findManyOrganizationsUseCase, findOneOrganizationByIDUseCase, updateOneOrganizationByIDUseCase, inactivateOneOrganizationByIDUseCase, reactivateOneOrganizationByIDUseCase)
	joinOrganizationRequestsHandler := handlers.NewJoinOrganizationRequests(createJoinOrganizationRequestUseCase, findManyJoinOrganizationRequestsByOrganizationIDUseCase, findManyJoinOrganizationRequestsByUserIDUseCase, acceptJoinOrganizationRequestUseCase, rejectJoinOrganizationRequestUseCase)
	joinOrganizationInvitesHandler := handlers.NewJoinOrganizationInvites(createJoinOrganizationInviteUseCase, findManyJoinOrganizationInvitesByOrganizationIDUseCase, findManyJoinOrganizationInvitesByUserIDUseCase, acceptJoinOrganizationInviteUseCase, rejectJoinOrganizationInviteUseCase)
	housingsHandler := handlers.NewHousings(createHousingUseCase, findManyHousingsByOrganizationIDUseCase, findOneHousingCompleteByIDUseCase, updateOneHousingByIDUseCase, deleteOneHousingByIDUseCase)
	updateOrganizationTypeRequestsHandler := handlers.NewUpdateOrganizationTypeRequests(createUpdateOrganizationTypeRequestUseCase, findManyUpdateOrganizationTypeRequestsUseCase, findManyUpdateOrganizationTypeRequestsByOrganizationIDUseCase, acceptUpdateOrganizationTypeRequestUseCase, rejectUpdateOrganizationTypeRequestUseCase)
	organizationsDataAccessRequestsHandler := handlers.NewOrganizationDataAccessRequests(createOrganizationDataAccessRequestsUseCase, findManyOrganizationDataAccessRequestsByRequesterUseCase, findManyOrganizationDataAccessRequestsByTargetUseCase, acceptOrganizationDataAccessRequestsUseCase, rejectOrganizationDataAccessRequestsUseCase)
	joinPlatformInvitesHandler := handlers.NewJoinPlatformInvites(createJoinPlatformInviteUseCase, findManyJoinPlatformInvitesByOrganizationIDUseCase, consumeJoinPlatformInviteByCodeUseCase)
	beneficiariesHandler := handlers.NewBeneficiaries(createBeneficiaryUseCase, findManyBeneficiariesByOrganizationIDUseCase, findManyBeneficiariesByHousingIDUseCase, findManyBeneficiariesByHousingRoomIDUseCase, findOneBeneficiaryCompleteByIDUseCase, updateBeneficiaryByIDUseCase, deleteBeneficiaryByIDUseCase, generateBeneficiaryProfileImageUploadLinkUseCase)
	housingRoomsHandler := handlers.NewHousingRooms(createHousingRoomUseCase, findOneHousingRoomCompleteByIDUseCase, findManyHousingRoomsByHousingIDUseCase, updateHousingRoomByIDUseCase, deleteHousingRoomByIDUseCase)
	beneficiaryAllocationsHandler := handlers.NewBeneficiaryAllocations(createEntranceBeneficiaryAllocationUseCase, createReallocationBeneficiaryAllocationUseCase, findManyBeneficiaryAllocationsByBeneficiaryIDUseCase, findManyBeneficiaryAllocationsByHousingIDUseCase, findManyBeneficiaryAllocationsByHousingRoomIDUseCase)
	voluntaryPeopleHandler := handlers.NewVoluntaryPeople(createVoluntaryPersonUseCase, findManyVoluntaryPeopleByOrganizationIDUseCase, findOneVoluntaryPersonCompleteByIDUseCase, updateOneVoluntaryPersonByIDUseCase, deleteOneVoluntaryPersonByIDUseCase)
	productTypesHandler := handlers.NewProductTypes(createProductTypeUseCase, findManyProductTypesByOrganizationIDUseCase, findOneProductTypeByIDUseCase, updateProductTypeByIDUseCase, deleteProductTypeByIDUseCase)
	organizationsDataAccessGrantsHandler := handlers.NewOrganizationDataAccessGrants(findManyGrantsByOrganizationIDUseCase, findManyGrantsByTargetOrganizationIDUseCase, deleteGrantByIDUseCase)
	productTypesAllocationHandler := handlers.NewProductTypeAllocations(createProductTypeAllocationEntranceUseCase, createProductTypeAllocationReallocationUseCase, findManyProductTypeAllocationsByProductTypeIDPaginatedUseCase)
	donationsHandler := handlers.NewDonations(createDonationUseCase, findManyDonationsByBeneficiaryIDUseCase, findManyDonationsByProductTypeIDUseCase)
	storageRecordsHandler := handlers.NewStorageRecords(findManyStorageRecordsByOrganizationIDUseCase, findManyStorageRecordsByHousingIDUseCase, findManyStorageRecordsByProductTypeIDUseCase)
	joinPlatformAdminInvitesHandler := handlers.NewJoinPlatformAdminInvites(createJoinPlatformAdminInvitesUseCase, findManyJoinPlatformAdminInvitesPaginatedUseCase, consumeJoinPlatformAdminInviteByCodeUseCase)

	healthHandler := handlers.NewHealth()

	router := http.NewRouter(
		settingsInstance,
		authenticateByCookieMiddleware,
		healthHandler,
		authenticationHandler,
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
		passwordRecoveryHandler,
		productTypesHandler,
		updateOrganizationTypeRequestsHandler,
		usersHandler,
		voluntaryPeopleHandler,
		productTypesAllocationHandler,
		donationsHandler,
		storageRecordsHandler,
		joinPlatformAdminInvitesHandler,
	)
	server := http.NewServer(router, settingsInstance.ServerPort)

	go func() {
		logger.Info("starting server", zap.String("port", settingsInstance.ServerPort))
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
