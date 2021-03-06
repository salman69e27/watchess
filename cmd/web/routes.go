package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New()
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)
	authDynamicMiddleware := dynamicMiddleware.Append(app.requireAdminUser)
	privilegeMiddleware := authDynamicMiddleware.Append(app.requireTournamentPrivilege)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/tournaments/create", authDynamicMiddleware.ThenFunc(app.createTournamentForm))
	mux.Post("/tournaments/create", authDynamicMiddleware.ThenFunc(app.createTournament))

	mux.Get("/tournaments/:id", dynamicMiddleware.ThenFunc(app.getTournamentPage))

	mux.Get("/tournaments/:id/rounds/create", privilegeMiddleware.ThenFunc(app.createRoundForm))
	mux.Post("/tournaments/:id/rounds/create", privilegeMiddleware.ThenFunc(app.createRound))

	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", authDynamicMiddleware.ThenFunc(app.logoutUser))

	mux.Get("/api/tournaments/:id", dynamicMiddleware.ThenFunc(app.getTournamentJSON))
	mux.Get("/api/rounds/:id", dynamicMiddleware.ThenFunc(app.getRoundGamesJSON))

	mux.Get("/ping", http.HandlerFunc(ping))

	fileServer := http.FileServer(http.Dir("./ui/"))
	mux.Get("/ui/static/css/", http.StripPrefix("/ui", fileServer))
	mux.Get("/ui/static/js/", http.StripPrefix("/ui", fileServer))
	return standardMiddleware.Then(mux)
}
