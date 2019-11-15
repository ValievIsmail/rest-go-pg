package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	log "github.com/sirupsen/logrus"
)

func createHTTPHandler(db *sql.DB) (http.Handler, error) {
	mux := chi.NewMux()

	mux.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(200)))
	})

	mux.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.SetHeader("Content-Type", "application/json; charset=utf-8;"))
		r.Use(cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			AllowCredentials: true,
			MaxAge:           30,
		}).Handler)

		r.Mount("/comment", commentHandler(db))
		r.Mount("/user", userHandler(db))
	})

	return mux, nil
}

func userHandler(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/", getAllUsersHandler(db))
		r.Post("/", createUserHandler(db))
		r.Put("/", updateUserHandler(db))
		r.Delete("/{uid}", deleteUserHandler(db))
	})

	return r
}

func commentHandler(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/{cid}", getCommentHandler(db))
		r.Get("/", getAllCommentsHandler(db))
		r.Post("/", createCommentHandler(db))
		r.Put("/", updateCommentHandler(db))
		r.Delete("/{cid}", deleteCommentHandler(db))
	})

	return r
}

func getAllUsersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := dbGetAllUsers(db)
		if err != nil {
			log.Errorf("getCommentHandler db: %v", err)
			http.Error(w, "getCommentHandler err", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(users); err != nil {
			log.Errorf("getCommentHandler write: %v", err)
			http.Error(w, "getCommentHandler write err", http.StatusInternalServerError)
			return
		}
	}
}

func createUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := Params{}

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "updateCommentHandler read invalid params", http.StatusBadRequest)
			return
		}

		id, err := dbCreateUser(p.Name, db)
		if err != nil {
			log.Errorf("createUserHandler db: %v", err)
			http.Error(w, "createUserHandler err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("updateUserHandler write id: %v", err)
			http.Error(w, "updateUserHandler write id", http.StatusInternalServerError)
			return
		}
	}
}

func updateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := Params{}

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "updateUserHandler read invalid params", http.StatusBadRequest)
			return
		}

		id, err := dbUpdateUser(p.UserID, p.Name, db)
		if err != nil {
			log.Errorf("dbUpdateUser: %v", err)
			http.Error(w, "dbUpdateUser err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("updateUserHandler write id: %v", err)
			http.Error(w, "updateUserHandler write id", http.StatusInternalServerError)
			return
		}
	}
}

func deleteUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := strconv.Atoi(chi.URLParam(r, "uid"))
		if err != nil {
			http.Error(w, "enter valid ID", http.StatusBadRequest)
			return
		}

		id, err := dbDeleteUser(uid, db)
		if err != nil {
			log.Errorf("deleteUserHandler db: %v", err)
			http.Error(w, "deleteUserHandler err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("updateCommentHandler write id: %v", err)
			http.Error(w, "updateCommentHandler write id", http.StatusInternalServerError)
			return
		}
	}
}

func getCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "cid"))
		if err != nil {
			http.Error(w, "enter valid ID", http.StatusBadRequest)
			return
		}

		comment, err := dbGetCommentByID(id, db)
		if err != nil {
			log.Errorf("getCommentHandler db: %v", err)
			http.Error(w, "getCommentHandler err", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(comment); err != nil {
			log.Errorf("getCommentHandler write: %v", err)
			http.Error(w, "getCommentHandler write err", http.StatusInternalServerError)
			return
		}
	}
}

func getAllCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comments, err := dbGetAllComments(db)
		if err != nil {
			log.Errorf("getAllCommentsHandler db: %v", err)
			http.Error(w, "getAllCommentsHandler err", http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(comments); err != nil {
			log.Errorf("getAllCommentsHandler write: %v", err)
			http.Error(w, "getAllCommentsHandler write err", http.StatusInternalServerError)
			return
		}
	}
}

func updateCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := Params{}

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "updateCommentHandler read invalid params", http.StatusBadRequest)
			return
		}

		id, err := dbUpdateComment(p.ID, p.Msg, db)
		if err != nil {
			log.Errorf("updateCommentHandler db: %v", err)
			http.Error(w, "updateCommentHandler err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("updateCommentHandler write id: %v", err)
			http.Error(w, "updateCommentHandler write id", http.StatusInternalServerError)
			return
		}
	}
}

func createCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := Params{}

		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "createCommentHandler read invalid params", http.StatusBadRequest)
			return
		}

		id, err := dbCreateComment(p.UserID, p.Msg, db)
		if err != nil {
			log.Errorf("createCommentHandler db: %v", err)
			http.Error(w, "createCommentHandler err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("createCommentHandler write id: %v", err)
			http.Error(w, "createCommentHandler write id", http.StatusInternalServerError)
			return
		}
	}
}

func deleteCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cid, err := strconv.Atoi(chi.URLParam(r, "cid"))
		if err != nil {
			http.Error(w, "enter valid ID", http.StatusBadRequest)
			return
		}

		id, err := dbDeleteComment(cid, db)
		if err != nil {
			log.Errorf("dbDeleteComment db: %v", err)
			http.Error(w, "dbDeleteComment err", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(&id); err != nil {
			log.Errorf("dbDeleteComment write id: %v", err)
			http.Error(w, "dbDeleteComment write id", http.StatusInternalServerError)
			return
		}
	}
}
