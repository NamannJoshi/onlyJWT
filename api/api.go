package api

import (
	"demo/token"
	"demo/utils"
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
			log.Println("errrrrorz", err)
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
	if r.Method == "GET" {
		return s.GetUserByID(w, r)
	}
	if r.Method == "PUT" {
		return s.UpdateUserByID(w, r)
	}
	return WriteJSON(w, http.StatusMethodNotAllowed, ApiError{Error: "method not allowed"})
}

func (s *ApiServer) CreateUser(w http.ResponseWriter, r *http.Request) error {
	log.Println("kha agye aap")
	var req CreateUserReq
	json.NewDecoder(r.Body).Decode(&req)

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("error in hashing password")
	}

	req.Password = passHash
	fmt.Println(req)
	user := NewUser(req)

	if err := s.store.CreateUserDb(user); err != nil {
		fmt.Println(err)
		return WriteJSON(w, http.StatusBadRequest, err)
	}
	return WriteJSON(w, http.StatusOK, &CreateUserRes{
		FullName: req.FullName,
		Email: req.Email,
		IsAdmin: req.IsAdmin,
		Number: req.Number,
	})
}

func (s *ApiServer) GetAllUsers(w http.ResponseWriter, r *http.Request) error {
	log.Println("apka swagat hai!!!!")
	users, err := s.store.GetUsersDb()
	log.Println(users)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}
	return WriteJSON(w, http.StatusOK, users)
}

func (s *ApiServer) GetUserByID(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("tommow is never")
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

	gu, err := s.store.GetUserByEmailDb(req.Email)
	fmt.Println(gu)
	if err != nil {
		return fmt.Errorf("err fetching user by email %w", err)
	}

	accessToken, err := s.token.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 15*time.Minute)
	if err != nil {
		return WriteJSON(w, http.StatusUnauthorized, ApiError{Error: err.Error()})
	}
	log.Println("here is your access token",accessToken)

	res := LoginUserRes{
		AccessToken: accessToken,
		UserRes: CreateUserRes{
			FullName: gu.FullName,
			Email: gu.Email,
			IsAdmin: gu.IsAdmin,
			Number: gu.Number,
		},
	}

	return WriteJSON(w, http.StatusOK, res)
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