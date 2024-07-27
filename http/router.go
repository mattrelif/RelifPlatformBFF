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
	authHandler *handlers.Auth,
	housingsHandler *handlers.Housings,
	joinOrganizationRequestsHandler *handlers.JoinOrganizationRequests,
	joinOrganizationInvitesHandler *handlers.JoinOrganizationInvites,
	organizationsHandler *handlers.Organizations,
	organizationDataAccessRequestsHandler *handlers.OrganizationDataAccessRequests,
	passwordHandler *handlers.Password,
	updateOrganizationTypeRequestsHandler *handlers.UpdateOrganizationTypeRequests,
	usersHandler *handlers.Users,
) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	router.Route(routerContext, func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", authHandler.SignUp)
			r.Post("/sign-in", authHandler.SignIn)

			r.Group(func(r chi.Router) {
				r.Use(authenticateByCookieMiddleware.Handle)

				r.Get("/me", authHandler.Me)
				r.Delete("/sign-out", authHandler.SignOut)
			})
		})

		r.Route("/password", func(r chi.Router) {
			r.Post("/request-change", passwordHandler.RequestChange)
			r.Put("/{id}", passwordHandler.Update)
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
				r.Put("/{id}", organizationsHandler.UpdateOne)
				r.Get("/{id}/users", usersHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-invites", joinOrganizationInvitesHandler.FindManyByOrganizationId)
				r.Get("/{id}/join-requests", joinOrganizationRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/data-accesses-requests", organizationDataAccessRequestsHandler.FindManyByRequesterOrganizationId)
				r.Get("/{id}/update-organization-type-requests", updateOrganizationTypeRequestsHandler.FindManyByOrganizationId)
				r.Get("/{id}/housings", housingsHandler.FindManyByOrganizationId)
			})

			r.Route("/join-organization-invites", func(r chi.Router) {
				r.Post("/", joinOrganizationInvitesHandler.Create)
				r.Delete("/{id}/accept", joinOrganizationInvitesHandler.Accept)
				r.Delete("/{id}/reject", joinOrganizationInvitesHandler.Reject)
			})

			r.Route("/join-organization-requests", func(r chi.Router) {
				r.Post("/", joinOrganizationRequestsHandler.Create)
				r.Delete("/{id}/accept", joinOrganizationRequestsHandler.Accept)
				r.Delete("/{id}/reject", joinOrganizationRequestsHandler.Reject)
			})

			r.Route("/organization-data-access-requests", func(r chi.Router) {
				r.Post("/", organizationDataAccessRequestsHandler.Create)
				r.Get("/", organizationDataAccessRequestsHandler.FindMany)
				r.Put("/{id}/accept", organizationDataAccessRequestsHandler.Accept)
				r.Put("/{id}/reject", organizationDataAccessRequestsHandler.Reject)
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
			})
		})
	})

	return router
}
