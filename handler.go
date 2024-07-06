package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-sql-driver/mysql"
)

var sqlDB *sql.DB

type Tasks struct {
	ID       *string `json:"id"`
	Name     string  `json:"username"`
	Verified bool    `json:"is_verified"`
	Salary   int64   `json:"salary"`
}

func initiateDB() error {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASSWORD"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDR"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Println("error while connecting to db: ", err)
		return err
	}
	if err := db.Ping(); err != nil {
		log.Println("error in ping call", err)
		return err
	}
	sqlDB = db
	log.Println("db Connected")
	return nil
}

func GetAllRoutes() http.Handler {
	router := chi.NewRouter()
	err := initiateDB()
	if err != nil {
		log.Fatal("db error: ", err)
		return nil
	}
	router.Get("/", listAllTasks)
	router.Get("/{id}", getSingleTask)

	return router
}

func listAllTasks(rw http.ResponseWriter, r *http.Request) {
	sqlRows, err := sqlDB.Query("select * from user.testTable")
	if err != nil {
		log.Println("error while getting values from DB: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer sqlRows.Close()
	// resp := []Tasks{}
	var resp []Tasks
	for sqlRows.Next() {

		var task Tasks
		// fmt.Println("sqlRows.Scan: ", sqlRows.Scan())
		if err := sqlRows.Scan(&task.ID, &task.Name, &task.Verified, &task.Salary); err != nil {
			log.Println("error while reading from sql ", err)
			// continue
			return
		}
		resp = append(resp, task)
	}
	if err := sqlRows.Err(); err != nil {
		log.Println("Error during rows iteration:", err)
		// Handle the error appropriately
	}

	if err := json.NewEncoder(rw).Encode(resp); err != nil {
		log.Println("error while encoding result")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("error while marshalling"))
	}

}

func getSingleTask(rw http.ResponseWriter, r *http.Request) {
	// Single Record from sql
	id := chi.URLParam(r, "id")
	sqlRow := sqlDB.QueryRow("select * from user.testTable where id = ?", id)
	var task Tasks
	if err := sqlRow.Scan(&task.ID, &task.Name, &task.Verified, &task.Salary); err != nil {
		log.Println("error while reading from sql ", err)
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte("Not Found"))
		// continue
		return
	}
	if err := json.NewEncoder(rw).Encode(task); err != nil {
		log.Println("error while encoding result")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("error while marshalling"))
	}
}
