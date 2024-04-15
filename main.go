package main

import (
	"encoding/json"
	"fmt"
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		panic(fmt.Sprintf("Usage: %s <api_token> <public_key> <frontend_api_url>", os.Args[0]))
	}
	apiToken := os.Args[1]
	publicKey := os.Args[2]
	frontendApiUrl := os.Args[3]
	client, err := clerk.NewClient(apiToken)
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.BasicAuth("Realm", map[string]string{"admin": "admin"}))
	r.Use(clerk.WithSessionV2(client))
	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`
		<!doctype html>
			<html lang="en">
			  <head>
				<meta charset="UTF-8" />
				<link rel="icon" type="image/svg+xml" href="/clerk.svg" />
				<meta name="viewport" content="width=device-width, initial-scale=1.0" />
				<title>Clerk + JavaScript Quickstart</title>
			  </head>
			  <body>
				<div id="app"></div>
				<script
				  async
				  crossorigin="anonymous"
				  data-clerk-publishable-key="` + publicKey + `"
				  src="` + frontendApiUrl + `/npm/@clerk/clerk-js@latest/dist/clerk.browser.js"
				  type="text/javascript"
				></script>
				<script>
				  window.addEventListener("load", async function () {
					await Clerk.load();
				 
					if (Clerk.user) {
					  document.getElementById("app").innerHTML = "<div id=\"user-button\"></div>";
				 
					  const userButtonDiv = document.getElementById("user-button");
				 
					  Clerk.mountUserButton(userButtonDiv);
					} else {
					  document.getElementById("app").innerHTML = "<div id=\"sign-in\"></div>";
				 
					  const signInDiv = document.getElementById("sign-in");
				 
					  Clerk.mountSignIn(signInDiv);
					}
				  });
				</script>
				<a href='/'>Home (should out put claims on home once signed in.)</a>
			  </body>
			</html>
				
		`))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		sessClaims, ok := clerk.SessionFromContext(r.Context())
		if ok {
			jsonStr, _ := json.Marshal(sessClaims)
			_, _ = w.Write(jsonStr)
		} else {
			_, _ = w.Write([]byte(`<!doctype html> <html lang=\"en\">
			<p>No session found in context.</p>
			<a href='/login'>Login</a>
			</html>`))
		}
	})
	fmt.Println("Server is running on port http://localhost:9000")
	err = http.ListenAndServe(":9000", r)
	if err != nil {
		panic(err)
	}
}
