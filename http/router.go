package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"relif/bff/http/handlers"
	"relif/bff/http/middlewares"
)

func NewRouter(
	routerContext string,
	authenticateByCookieMiddleware *middlewares.AuthenticateByCookie,
	rbacMiddleware *middlewares.RoleBasedAccessControl,
	authHandler *handlers.Auth,
	beneficiariesHandler *handlers.Beneficiaries,
	beneficiaryAllocationsHandler *handlers.BeneficiaryAllocations,
	housingsHandler *handlers.Housings,
	housingRoomsHandler *handlers.HousingRooms,
	joinOrganizationRequestsHandler *handlers.JoinOrganizationRequests,
	joinOrganizationInvitesHandler *handlers.JoinOrganizationInvites,
	joinPlatformInvitesHandler *handlers.JoinPlatformInvites,
	organizationsHandler *handlers.Organizations,
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
			})

			r.Route("/organizations", func(r chi.Router) {
				r.Post("/", organizationsHandler.Create)
				r.Get("/", organizationsHandler.FindMany)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}", organizationsHandler.UpdateOne)
				r.Get("/{id}/users", usersHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-invites", joinOrganizationInvitesHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-requests", joinOrganizationRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/data-accesses-requests", organizationDataAccessRequestsHandler.FindManyByRequesterOrganizationId)
				r.Get("/{id}/update-organization-type-requests", updateOrganizationTypeRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/housings", housingsHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-platform-invites", joinPlatformInvitesHandler.FindManyByOrganizationId)
				r.Get("/{id}/voluntary-people", voluntaryPeopleHandler.FindManyByOrganizationId)
				r.Get("/{id}/product-types", productTypesHandler.FindManyByOrganizationId)
			})

			r.Route("/join-organization-invites", func(r chi.Router) {
				r.With(rbacMiddleware.Middleware([]string{})).Post("/", joinOrganizationInvitesHandler.Create)
				r.Delete("/{id}/accept", joinOrganizationInvitesHandler.Accept)
				r.Delete("/{id}/reject", joinOrganizationInvitesHandler.Reject)
			})

			r.Route("/join-organization-requests", func(r chi.Router) {
				r.Post("/", joinOrganizationRequestsHandler.Create)
				r.With(rbacMiddleware.Middleware([]string{})).Delete("/{id}/accept", joinOrganizationRequestsHandler.Accept)
				r.With(rbacMiddleware.Middleware([]string{})).Delete("/{id}/reject", joinOrganizationRequestsHandler.Reject)
			})

			r.Route("/organization-data-access-requests", func(r chi.Router) {
				r.With(rbacMiddleware.Middleware([]string{})).Post("/", organizationDataAccessRequestsHandler.Create)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/", organizationDataAccessRequestsHandler.FindMany)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}/accept", organizationDataAccessRequestsHandler.Accept)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}/reject", organizationDataAccessRequestsHandler.Reject)
			})

			r.Route("/update-organization-type-requests", func(r chi.Router) {
				r.With(rbacMiddleware.Middleware([]string{})).Post("/", updateOrganizationTypeRequestsHandler.Create)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/", updateOrganizationTypeRequestsHandler.FindMany)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}/accept", updateOrganizationTypeRequestsHandler.Accept)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}/reject", updateOrganizationTypeRequestsHandler.Reject)
			})

			r.Route("/housings", func(r chi.Router) {
				r.With(rbacMiddleware.Middleware([]string{})).Post("/", housingsHandler.Create)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}", housingsHandler.FindOneById)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}", housingsHandler.Update)
				r.With(rbacMiddleware.Middleware([]string{})).Delete("/{id}", housingsHandler.Delete)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/rooms", housingRoomsHandler.FindManyByHousingId)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByHousingId)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByHousingId)

			})

			r.Route("/housing-rooms", func(r chi.Router) {
				r.With(rbacMiddleware.Middleware([]string{})).Post("/", housingRoomsHandler.Create)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}", housingRoomsHandler.FindOneById)
				r.With(rbacMiddleware.Middleware([]string{})).Put("/{id}", housingRoomsHandler.Update)
				r.With(rbacMiddleware.Middleware([]string{})).Delete("/{id}", housingRoomsHandler.Delete)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/beneficiaries", beneficiariesHandler.FindManyByRoomId)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByRoomId)
			})

			r.Route("/beneficiaries", func(r chi.Router) {
				r.Post("/", beneficiariesHandler.Create)
				r.Get("/{id}", beneficiariesHandler.FindOneById)
				r.Put("/{id}", beneficiariesHandler.Update)
				r.Delete("/{id}", beneficiariesHandler.Delete)
				r.With(rbacMiddleware.Middleware([]string{})).Post("/{id}/allocate", beneficiaryAllocationsHandler.Allocate)
				r.With(rbacMiddleware.Middleware([]string{})).Post("/{id}/reallocate", beneficiaryAllocationsHandler.Reallocate)
				r.With(rbacMiddleware.Middleware([]string{})).Get("/{id}/allocations", beneficiaryAllocationsHandler.FindManyByBeneficiaryId)
			})

			r.Route("/voluntary-people", func(r chi.Router) {
				r.Post("/", voluntaryPeopleHandler.Create)
				r.Get("/{id}", voluntaryPeopleHandler.FindOneById)
				r.Put("/{id}", voluntaryPeopleHandler.Update)
				r.Delete("/{id}", voluntaryPeopleHandler.DeleteById)
			})

			r.Route("/product-types", func(r chi.Router) {
				r.Post("/", productTypesHandler.Create)
				r.Get("/{id}", productTypesHandler.FindOneById)
				r.Put("/{id}", productTypesHandler.Update)
				r.Delete("/{id}", productTypesHandler.Delete)
			})
		})
	})

	return router
}
