## Url Shortener

A web api/cli program to shorten urls.

### Usage

##### Web Api

To build the web api run `go run main.go` in the `cmd/webserver` folder.

Then then hit the server with a json post request at `8000/url` to generate a shorten url.
Json should be in the form:

```json
{
 "long_url": "http://<your_url_here>" 
}
```

The web api will then redirect when you direct a browser to the given url.


##### Cli

To build the cli version run `go run main.go` in the `cmd/cli` folder.

Then run the program created.

### TODO:

* Better error handling. Could DRY in lots of areas. See: [Behaviour type assertion] 
* System to remove old urls. Would be sensible to remove urls over a certain age (maybe check if they have been used recently?)


[Behaviour type assertion]: https://medium.com/ki-labs-engineering/rest-api-error-handling-in-go-behavioral-type-assertion-509d93636afd
