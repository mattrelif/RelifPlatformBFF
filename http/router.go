package http

import (
	"fmt"
	"net/http"
	"relif/platform-bff/http/handlers"
	"relif/platform-bff/http/middlewares"
	"relif/platform-bff/settings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(
	settingsInstance *settings.Settings,
	authenticateByTokenMiddleware *middlewares.AuthenticateByToken,
	healthHandler *handlers.Health,
	authenticationHandler *handlers.Authentication,
	beneficiariesHandler *handlers.Beneficiaries,
	beneficiaryAllocationsHandler *handlers.BeneficiaryAllocations,
	housingsHandler *handlers.Housings,
	housingRoomsHandler *handlers.HousingRooms,
	joinOrganizationRequestsHandler *handlers.JoinOrganizationRequests,
	joinOrganizationInvitesHandler *handlers.JoinOrganizationInvites,
	joinPlatformInvitesHandler *handlers.JoinPlatformInvites,
	organizationsHandler *handlers.Organizations,
	organizationDataAccessGrantHandler *handlers.OrganizationDataAccessGrants,
	organizationDataAccessRequestsHandler *handlers.OrganizationDataAccessRequests,
	passwordRecoveryHandler *handlers.PasswordRecovery,
	productTypesHandler *handlers.ProductTypes,
	updateOrganizationTypeRequestsHandler *handlers.UpdateOrganizationTypeRequests,
	usersHandler *handlers.Users,
	voluntaryPeopleHandler *handlers.VoluntaryPeople,
	productTypeAllocationsHandler *handlers.ProductTypeAllocations,
	donationsHandler *handlers.Donations,
	storageRecordsHandler *handlers.StorageRecords,
	joinPlatformAdminInvites *handlers.JoinPlatformAdminInvites,
	casesHandler *handlers.Cases,
	caseNotesHandler *handlers.CaseNotes,
	caseDocumentsHandler *handlers.CaseDocuments,
) http.Handler {
	router := chi.NewRouter()

	// Simple health check for load balancers
	router.Get("/health", healthHandler.SimpleHealthCheck)

	router.Route(settingsInstance.RouterContext, func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.Logger)
		r.Use(middleware.SetHeader("Content-Type", "application/json"))

		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{
				fmt.Sprintf("http://%s", settingsInstance.FrontendDomain),
				fmt.Sprintf("https://%s", settingsInstance.FrontendDomain),
				"https://app.relifaid.org", // Production frontend
				"http://localhost:3000",    // Development frontend
				"https://localhost:3000",   // Development frontend with HTTPS
			},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		}))

		// Detailed health check with database status
		r.Get("/health", healthHandler.HealthCheck)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", authenticationHandler.SignUp)
			r.Post("/org-sign-up", authenticationHandler.OrganizationSignUp)
			r.Post("/admin-sign-up", authenticationHandler.AdminSignUp)
			r.Post("/sign-in", authenticationHandler.SignIn)
			r.Get("/verify-email", authenticationHandler.Verify)
			r.With(authenticateByTokenMiddleware.Handle).Get("/me", authenticationHandler.Me)
			r.With(authenticateByTokenMiddleware.Handle).Delete("/sign-out", authenticationHandler.SignOut)
		})

		r.Route("/password", func(r chi.Router) {
			r.Post("/request-change", passwordRecoveryHandler.RequestChange)
			r.Put("/{code}", passwordRecoveryHandler.Update)
		})

		r.Route("/join-platform-invites", func(r chi.Router) {
			r.With(authenticateByTokenMiddleware.Handle).Post("/", joinPlatformInvitesHandler.Create)
			r.Delete("/{code}/consume", joinPlatformInvitesHandler.Consume)
		})

		r.Route("/join-platform-admin-invites", func(r chi.Router) {
			r.With(authenticateByTokenMiddleware.Handle).Post("/", joinPlatformAdminInvites.Create)
			r.With(authenticateByTokenMiddleware.Handle).Get("/", joinPlatformAdminInvites.FindManyPaginated)
			r.Delete("/{code}/consume", joinPlatformAdminInvites.ConsumeByCode)
		})

		r.Group(func(r chi.Router) {
			r.Use(authenticateByTokenMiddleware.Handle)

			r.Route("/users", func(r chi.Router) {
				r.Get("/relif-members", usersHandler.FindManyRelifMembers)
				r.Get("/{id}", usersHandler.FindOne)
				r.Put("/{id}", usersHandler.UpdateOne)
				r.Delete("/{id}", usersHandler.InactivateOne)
				r.Put("/{id}/reactivate", usersHandler.ReactivateOne)

				r.Get("/{id}/join-organization-requests", joinOrganizationRequestsHandler.FindManyByUserID)
				r.Get("/{id}/join-organization-invites", joinOrganizationInvitesHandler.FindManyByUserID)
			})

			r.Route("/organizations", func(r chi.Router) {
				r.Post("/", organizationsHandler.Create)

				r.Get("/", organizationsHandler.FindMany)
				r.Get("/{id}", organizationsHandler.FindOne)
				r.Get("/{id}/users", usersHandler.FindManyByOrganizationID)
				r.Get("/{id}/join-invites", joinOrganizationInvitesHandler.FindManyByOrganizationID)
				r.Get("/{id}/join-requests", joinOrganizationRequestsHandler.FindManyByOrganizationID)
				r.Get("/{id}/requested-data-access-requests", organizationDataAccessRequestsHandler.FindManyByRequesterOrganizationID)
				r.Get("/{id}/targeted-data-access-requests", organizationDataAccessRequestsHandler.FindManyByTargetOrganizationID)
				r.Get("/{id}/data-access-grants", organizationDataAccessGrantHandler.FindManyByOrganizationID)
				r.Get("/{id}/targeted-data-access-grants", organizationDataAccessGrantHandler.FindManyByTargetOrganizationID)
				r.Get("/{id}/update-organization-type-requests", updateOrganizationTypeRequestsHandler.FindManyByOrganizationID)
				r.Get("/{id}/housings", housingsHandler.FindManyByOrganizationID)
				r.Get("/{id}/join-platform-invites", joinPlatformInvitesHandler.FindManyByOrganizationID)
				r.Get("/{id}/voluntary-people", voluntaryPeopleHandler.FindManyByOrganizationID)
				r.Get("/{id}/product-types", productTypesHandler.FindManyByOrganizationID)
				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByOrganizationID)
				r.Get("/{id}/storage-records", storageRecordsHandler.FindManyByOrganizationID)

				r.Put("/{id}", organizationsHandler.UpdateOne)
				r.Put("/{id}/reactivate", organizationsHandler.ReactivateOne)

				r.Post("/{id}/join-organization-requests", joinOrganizationRequestsHandler.Create)
				r.Post("/{id}/request-organization-data-access", organizationDataAccessRequestsHandler.Create)
				r.Post("/{id}/beneficiaries", beneficiariesHandler.Create)
				r.Post("/{id}/voluntary-people", voluntaryPeopleHandler.Create)
				r.Post("/{id}/product-types", productTypesHandler.Create)

				r.Delete("/{id}", organizationsHandler.Delete)
			})

			r.Route("/join-organization-invites", func(r chi.Router) {
				r.Post("/", joinOrganizationInvitesHandler.Create)
				r.Put("/{id}/accept", joinOrganizationInvitesHandler.Accept)
				r.Put("/{id}/reject", joinOrganizationInvitesHandler.Reject)
			})

			r.Route("/join-organization-requests", func(r chi.Router) {
				r.Put("/{id}/accept", joinOrganizationRequestsHandler.Accept)
				r.Put("/{id}/reject", joinOrganizationRequestsHandler.Reject)
			})

			r.Route("/organization-data-access-requests", func(r chi.Router) {
				r.Put("/{id}/accept", organizationDataAccessRequestsHandler.Accept)
				r.Put("/{id}/reject", organizationDataAccessRequestsHandler.Reject)
			})

			r.Route("/organization-data-access-grants", func(r chi.Router) {
				r.Delete("/{id}", organizationDataAccessGrantHandler.Delete)
			})

			r.Route("/update-organization-type-requests", func(r chi.Router) {
				r.Post("/", updateOrganizationTypeRequestsHandler.Create)
				r.Get("/", updateOrganizationTypeRequestsHandler.FindMany)
				r.Put("/{id}/accept", updateOrganizationTypeRequestsHandler.Accept)
				r.Put("/{id}/reject", updateOrganizationTypeRequestsHandler.Reject)
			})

			r.Route("/housings", func(r chi.Router) {
				r.Post("/", housingsHandler.Create)
				r.Get("/{id}", housingsHandler.FindOneByID)
				r.Put("/{id}", housingsHandler.Update)
				r.Delete("/{id}", housingsHandler.Delete)

				r.Get("/{id}/rooms", housingRoomsHandler.FindManyByHousingID)
				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByHousingID)
				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByHousingID)
				r.Get("/{id}/storage-records", storageRecordsHandler.FindManyByHousingID)

				r.Post("/{id}/rooms", housingRoomsHandler.Create)
			})

			r.Route("/housing-rooms", func(r chi.Router) {
				r.Get("/{id}", housingRoomsHandler.FindOneByID)
				r.Put("/{id}", housingRoomsHandler.Update)
				r.Delete("/{id}", housingRoomsHandler.Delete)

				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByRoomID)
				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByRoomID)
			})

			r.Route("/beneficiaries", func(r chi.Router) {
				r.Get("/{id}", beneficiariesHandler.FindOneByID)
				r.Put("/{id}", beneficiariesHandler.Update)
				r.Delete("/{id}", beneficiariesHandler.Delete)

				r.Post("/{id}/allocate", beneficiaryAllocationsHandler.Allocate)
				r.Post("/{id}/reallocate", beneficiaryAllocationsHandler.Reallocate)
				r.Post("/{id}/donations", donationsHandler.Create)

				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByBeneficiaryID)
				r.Get("/{id}/donations", donationsHandler.FindManyByBeneficiaryID)
				r.Get("/{id}/cases", casesHandler.FindManyByBeneficiaryID)

				r.Post("/generate-profile-image-upload-link", beneficiariesHandler.GenerateProfileImageUploadLink)
			})

			r.Route("/voluntary-people", func(r chi.Router) {
				r.Get("/{id}", voluntaryPeopleHandler.FindOneByID)
				r.Put("/{id}", voluntaryPeopleHandler.Update)
				r.Delete("/{id}", voluntaryPeopleHandler.Delete)
			})

			r.Route("/product-types", func(r chi.Router) {
				r.Get("/{id}", productTypesHandler.FindOneByID)
				r.Put("/{id}", productTypesHandler.Update)
				r.Delete("/{id}", productTypesHandler.Delete)

				r.Post("/{id}/allocate", productTypeAllocationsHandler.Allocate)
				r.Post("/{id}/reallocate", productTypeAllocationsHandler.Reallocate)

				r.Get("/{id}/donations", donationsHandler.FindManyByProductTypeID)
				r.Get("/{id}/allocations", productTypeAllocationsHandler.FindManyByProductTypeID)
				r.Get("/{id}/storage-records", storageRecordsHandler.FindManyByProductTypeID)
			})

			r.Route("/cases", func(r chi.Router) {
				r.Post("/", casesHandler.CreateCase)
				r.Get("/", casesHandler.FindManyByOrganization)
				r.Get("/stats", casesHandler.GetStats)
				r.Get("/{id}", casesHandler.FindOne)
				r.Put("/{id}", casesHandler.UpdateOne)
				r.Delete("/{id}", casesHandler.DeleteOne)

				r.Route("/{case_id}/notes", func(r chi.Router) {
					r.Get("/", caseNotesHandler.ListByCaseID)
					r.Post("/", caseNotesHandler.Create)
					r.Put("/{note_id}", caseNotesHandler.Update)
					r.Delete("/{note_id}", caseNotesHandler.Delete)
				})

				r.Route("/{case_id}/documents", func(r chi.Router) {
					r.Get("/", caseDocumentsHandler.ListByCaseID)
					r.Post("/generate-upload-link", caseDocumentsHandler.GenerateUploadLink)
					r.Post("/", caseDocumentsHandler.Create)
					r.Put("/{doc_id}", caseDocumentsHandler.Update)
					r.Delete("/{doc_id}", caseDocumentsHandler.Delete)
				})
			})
		})
	})

	return router
}
