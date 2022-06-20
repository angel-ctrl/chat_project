package middlewares

import (
	"net/http"

	"github.com/sql_chat/utils"
)

func ValidoJWT(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "POST" || r.URL.Path == "/API/addfriend/v1" {

			_, _, err := utils.ProcesoToken(r.Header.Get("Authorization"))

			if err != nil {
				http.Error(w, "Error en el token "+err.Error(), http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
		} else {
			next.ServeHTTP(w, r)
		}

	}
}
