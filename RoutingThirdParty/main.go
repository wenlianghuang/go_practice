package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// new router instance
	mux := mux.NewRouter()

	// define routes

	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Marhaba Misterr Wick 🧘🏽🧘🏽🧘🏽"))
	}))

	mux.Handle("/api/v1/user/", UserMux())

	mux.Handle("/api/v1/auth/", AuthMux())

	log.Printf("App running on port: %v\n", 8080)
	http.ListenAndServe(":8080", mux)

}

func AuthMux() http.Handler {
	authMux := http.NewServeMux()
	authMux.Handle("/signup", Post(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("You All Signed Up Misterr Wick 🧘🏽🧘🏽🧘🏽"))
	})))
	authMux.Handle("/resendVerificationEmail", Post(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Your Access Has Been Resent Misterr Wick 🧘🏽🧘🏽🧘🏽"))
	})))

	return http.StripPrefix("/api/v1/auth", authMux)
}

func UserMux() http.Handler {
	userMux := http.NewServeMux()
	userMux.Handle("/profile", Get(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Always a Pleasure To See You Misterr Wick 🧘🏽🧘🏽🧘🏽"))
	})))
	userMux.Handle("/weapons", Get(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("We Got The Finest Set Of Cutlery Here Just As You Like It Misterr Wick 🧘🏽🧘🏽🧘🏽"))
	})))

	return http.StripPrefix("/api/v1/user", userMux)
}

func Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Put(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Patch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPatch {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}

func Delete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodDelete {
			http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(res, req)
	})
}
