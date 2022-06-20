package middlewares

/*import (
	"fmt"
	"net/http"

	db "github.com/sql_study/database"
)

func CheckBD(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if err := db.PingDataBase(); err != nil {
			fmt.Println(err)
			http.Error(w, "Conexion perdida con la base de datos ", 500)
			return
		}

		next.ServeHTTP(w, r)

	}

}*/
