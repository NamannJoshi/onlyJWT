package api

import (
	"demo/token"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error 
type ApiError struct {
	Error string
}

func makeHTTPhandler(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type ApiServer struct {
	listenAddr string
	store      Storage
	token *token.JWTMaker
}

func NewApiServer(listenAddr string, store Storage, secretKey string) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
		token: token.NewJWTMaker(secretKey),
	}
}

func (s *ApiServer) Run() {
	router := mux.NewRouter()
	
	router.HandleFunc("/user/login", makeHTTPhandler(s.LoginUser)).Methods("POST")
	router.HandleFunc("/user/{id}", makeHTTPhandler(s.AuthMiddleware(s.HandleAccountByID)))
	router.HandleFunc("/user", makeHTTPhandler(s.HandleAccount))

	http.ListenAndServe(s.listenAddr, router)
}

func (s *ApiServer) HandleAccount(w http.ResponseWriter, r *http.Request)error{
	if r.Method == http.MethodGet {
		return s.GetAllUsers(w, r)
	}
	if r.Method == http.MethodPost {
		return s.CreateUser(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "method not allowed"})
}

func (s *ApiServer) HandleAccountByID(w http.ResponseWriter, r *http.Request)error{
	if r.Method == http.MethodGet {
		return s.GetUserByID(w, r)
	}
	if r.Method == http.MethodPut {
		return s.UpdateUserByID(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "method not allowed"})
}

func (s *ApiServer) CreateUser(w http.ResponseWriter, r *http.Request) error {
	log.Println("kha agye aap")
	var req CreateUserReq
	json.NewDecoder(r.Body).Decode(&req)

	user := NewUser(req)

	s.store.CreateUserDb(user)
	return nil
}
func (s *ApiServer) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	u := CreateUserReq{
		FullName: "gibb",
		Email: "gialkfabb",
		IsAdmin: false,
		Number: 3242,
	}
	user := NewUser(u)
	return WriteJSON(w, http.StatusOK, user)
}

func (s *ApiServer) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, ApiError{Error: "you go girly"})
}

func (s *ApiServer) GetUserByEmail(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) UpdateUserByID(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *ApiServer) LoginUser(w http.ResponseWriter, r *http.Request) error {
	log.Println("login endpoint accesssed successfully")
	var req LoginUserReq
	json.NewDecoder(r.Body).Decode(&req)

	accessToken, err := s.token.CreateToken(23, "kdlafjfakfa", false, 15*time.Minute)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
	}
	log.Println("here is your access token",accessToken)

	return WriteJSON(w, http.StatusOK, accessToken)
}

func (s *ApiServer) AuthMiddleware(f ApiFunc) ApiFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		authStr := r.Header.Get("Authorization")
		if authStr == "" {
			fmt.Errorf("empty auth header")
		}

		filteredAuth := strings.TrimSpace(strings.Replace(authStr, "Bearer", "", 1))
		filteredAuth = strings.Trim(filteredAuth, "\"")
		log.Println(filteredAuth)


		// log.Println(authStr)
		authClaims, err := s.token.VerifyToken(filteredAuth)
		if err != nil {
			return fmt.Errorf("you have verifying non defying lode lagaing error: %v", err)
		}
		
		log.Println(authClaims)
		return WriteJSON(w, http.StatusOK, authClaims)
	}
}