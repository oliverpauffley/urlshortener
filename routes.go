package urlshortener

// define routes for handlers
func (env *Env) Routes() {
	env.Router.HandleFunc("/url", env.urlHandler())
	env.Router.HandleFunc("/", env.RedirectHandler())
}
