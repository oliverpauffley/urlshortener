package main

// define routes using DefaultServeMux
func (env *Env) routes() {
	env.router.HandleFunc("/url", env.urlHandler())
}
