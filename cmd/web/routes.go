package main

import (
	"net/http"

	"cmd/web/main.go/ui"

	"github.com/justinas/alice"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.FileServerFS(ui.Files)) // Use the http.FileServerFS() function to create a HTTP handler which serves the embedded files in ui.Files. It's important to note that our static files are contained in the "static" folder of the ui.Files embedded filesystem. So, for example, our CSS stylesheet is located at static/css/main.css". This means that we no longer need to strip the prefix from the request URL -- any requests that start with /static/ can just be passed directly to the file server and the corresponding static file will be served (so long as it exists).
	mux.HandleFunc("GET /ping", ping)
	
	// Unprotected application routes using the "dynamic" middleware chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/UpdateSnippet/{id}", protected.ThenFunc(app.snippetGetIdEdit))
	mux.Handle("POST /snippet/UpdateSnippets/{id}", protected.ThenFunc(app.snippetUpdatePost))
	mux.Handle("GET /snippet/Delete/{id}", protected.ThenFunc(app.snippetDelete))
	//Users Route
	mux.Handle("GET /user/show", protected.ThenFunc(app.userShow))
	mux.Handle("GET /snippet/UsersEditPage/{id}", protected.ThenFunc(app.UsersEditPage))
	mux.Handle("POST /snippet/UsersEditPage/{id}", protected.ThenFunc(app.usersUpdatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	// Websocket
	mux.Handle("/websocketIndex", protected.ThenFunc(app.websocketIndex))
	mux.Handle("/ws", protected.ThenFunc(app.websocketMain))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders) // alice.New(mymiddleware1, mymiddleware2, mymiddleware3).then(myhandle)
	// Return the 'standard' middleware chain followed by the servemux.
	return standard.Then(mux)
}
