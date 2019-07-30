package main

// define routes using DefaultServeMux
func (env *Env) routes() {
	env.router.HandleFunc("/url", env.urlHandler())
	env.router.HandleFunc("/", env.RedirectHandler())
}
