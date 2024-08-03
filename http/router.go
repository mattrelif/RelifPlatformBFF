package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"relif/bff/http/handlers"
	"relif/bff/http/middlewares"
)

func NewRouter(
	environment,
	routerContext string,
	authenticateByCookieMiddleware *middlewares.AuthenticateByCookie,
	authHandler *handlers.Auth,
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
	passwordHandler *handlers.Password,
	productTypesHandler *handlers.ProductTypes,
	updateOrganizationTypeRequestsHandler *handlers.UpdateOrganizationTypeRequests,
	usersHandler *handlers.Users,
	voluntaryPeopleHandler *handlers.VoluntaryPeople,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	if environment == "development" {
		router.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:3000"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"*"},
			AllowCredentials: true,
			MaxAge:           300,
		}))
	}

	router.Route(routerContext, func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", authHandler.SignUp)
			r.Post("/org-sign-up", authHandler.OrganizationSignUp)
			r.Post("/sign-in", authHandler.SignIn)
			r.With(authenticateByCookieMiddleware.Handle).Get("/me", authHandler.Me)
			r.With(authenticateByCookieMiddleware.Handle).Delete("/sign-out", authHandler.SignOut)
		})

		r.Route("/password", func(r chi.Router) {
			r.Post("/request-change", passwordHandler.RequestChange)
			r.Put("/{code}", passwordHandler.Update)
		})

		r.Route("/join-platform-invites", func(r chi.Router) {
			r.With(authenticateByCookieMiddleware.Handle).Post("/", joinPlatformInvitesHandler.Create)
			r.Delete("/{code}/consume", joinPlatformInvitesHandler.Consume)
		})

		r.Group(func(r chi.Router) {
			r.Use(authenticateByCookieMiddleware.Handle)

			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", usersHandler.FindOne)
				r.Put("/{id}", usersHandler.UpdateOne)
				r.Delete("/{id}", usersHandler.DeleteOne)

				r.Get("/{id}/join-organization-invites", joinOrganizationInvitesHandler.FindManyByUserId)
			})

			r.Route("/organizations", func(r chi.Router) {
				r.Post("/", organizationsHandler.Create)
				r.Get("/", organizationsHandler.FindMany)
				r.Get("/{id}", organizationsHandler.FindOne)
				r.Put("/{id}", organizationsHandler.UpdateOne)
				r.Get("/{id}/users", usersHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-invites", joinOrganizationInvitesHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-requests", joinOrganizationRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/requested-data-access-requests", organizationDataAccessRequestsHandler.FindManyByRequesterOrganizationId)
				r.Get("/{id}/targeted-data-access-requests", organizationDataAccessRequestsHandler.FindManyByTargetOrganizationId)
				r.Get("/{id}/data-access-grants", organizationDataAccessGrantHandler.FindManyByOrganizationId)
				r.Get("/{id}/targeted-data-access-grants", organizationDataAccessGrantHandler.FindManyByOrganizationId)
				r.Get("/{id}/update-organization-type-requests", updateOrganizationTypeRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/housings", housingsHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-platform-invites", joinPlatformInvitesHandler.FindManyByOrganizationId)
				r.Get("/{id}/voluntary-people", voluntaryPeopleHandler.FindManyByOrganizationId)
				r.Get("/{id}/product-types", productTypesHandler.FindManyByOrganizationId)
				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByOrganizationId)

				r.Post("/{id}/join-organization-requests", joinOrganizationRequestsHandler.Create)
				r.Post("/{id}/request-organization-data-access", organizationDataAccessRequestsHandler.Create)
				r.Post("/{id}/beneficiaries", beneficiariesHandler.Create)
				r.Post("/{id}/voluntary-people", voluntaryPeopleHandler.Create)
				r.Post("/{id}/product-types", productTypesHandler.Create)
			})

			r.Route("/join-organization-invites", func(r chi.Router) {
				r.Post("/", joinOrganizationInvitesHandler.Create)
				r.Delete("/{id}/accept", joinOrganizationInvitesHandler.Accept)
				r.Delete("/{id}/reject", joinOrganizationInvitesHandler.Reject)
			})

			r.Route("/join-organization-requests", func(r chi.Router) {
				r.Delete("/{id}/accept", joinOrganizationRequestsHandler.Accept)
				r.Delete("/{id}/reject", joinOrganizationRequestsHandler.Reject)
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
				r.Get("/{id}", housingsHandler.FindOneById)
				r.Put("/{id}", housingsHandler.Update)
				r.Delete("/{id}", housingsHandler.Delete)
				r.Get("/{id}/rooms", housingRoomsHandler.FindManyByHousingId)
				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByHousingId)
				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByHousingId)

				r.Post("/{id}/rooms", housingRoomsHandler.Create)
			})

			r.Route("/housing-rooms", func(r chi.Router) {
				r.Get("/{id}", housingRoomsHandler.FindOneById)
				r.Put("/{id}", housingRoomsHandler.Update)
				r.Delete("/{id}", housingRoomsHandler.Delete)
				r.Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByRoomId)
				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByRoomId)
			})

			r.Route("/beneficiaries", func(r chi.Router) {
				r.Get("/{id}", beneficiariesHandler.FindOneById)
				r.Put("/{id}", beneficiariesHandler.Update)
				r.Delete("/{id}", beneficiariesHandler.Delete)
				r.Post("/{id}/allocate", beneficiaryAllocationsHandler.Allocate)
				r.Post("/{id}/reallocate", beneficiaryAllocationsHandler.Reallocate)
				r.Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByBeneficiaryId)
			})

			r.Route("/voluntary-people", func(r chi.Router) {
				r.Get("/{id}", voluntaryPeopleHandler.FindOneById)
				r.Put("/{id}", voluntaryPeopleHandler.Update)
				r.Delete("/{id}", voluntaryPeopleHandler.DeleteById)
			})

			r.Route("/product-types", func(r chi.Router) {
				r.Get("/{id}", productTypesHandler.FindOneById)
				r.Put("/{id}", productTypesHandler.Update)
				r.Delete("/{id}", productTypesHandler.Delete)
			})
		})
	})

	return router
}
