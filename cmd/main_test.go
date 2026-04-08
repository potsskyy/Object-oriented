package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"todo/internal/handler"
	"todo/internal/middleware"
	"todo/internal/models"
	"todo/internal/repository"
	"todo/internal/session"
)

func newTestRouter() http.Handler {
	store := repository.NewMemoryRepo()

	authHandler := handler.NewAuthHandler(store)
	taskHandler := handler.NewTaskHandler(store)

	router := mux.NewRouter()
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	todoRouter := router.PathPrefix("/").Subrouter()
	todoRouter.Use(middleware.WithAuth(store))
	todoRouter.HandleFunc("/add", taskHandler.Add).Methods("POST")
	todoRouter.HandleFunc("/update", taskHandler.Update).Methods("POST")
	todoRouter.HandleFunc("/resolve/{id}", taskHandler.Resolve).Methods("POST")
	todoRouter.HandleFunc("/delete/{id}", taskHandler.Delete).Methods("POST")
	todoRouter.HandleFunc("/get", taskHandler.GetAll).Methods("GET")
	todoRouter.HandleFunc("/get/{id}", taskHandler.GetByID).Methods("GET")
	todoRouter.HandleFunc("/archive", taskHandler.GetArchive).Methods("GET")

	return router
}

func doJSONRequest(t *testing.T, router http.Handler, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var bodyReader *bytes.Reader
	if body == nil {
		bodyReader = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		bodyReader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func sessionCookieFromResponse(t *testing.T, rr *httptest.ResponseRecorder) *http.Cookie {
	t.Helper()

	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == session.CookieName {
			return cookie
		}
	}

	t.Fatal("session cookie not found in response")
	return nil
}

func TestLoginStoresOpaqueSessionAndLogoutInvalidatesIt(t *testing.T) {
	router := newTestRouter()

	registerResponse := doJSONRequest(t, router, http.MethodPost, "/register", models.User{
		Username: "admin",
		Password: "1234",
	})
	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerResponse.Code, http.StatusCreated)
	}

	loginResponse := doJSONRequest(t, router, http.MethodPost, "/login", models.User{
		Username: "admin",
		Password: "1234",
	})
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.Code, http.StatusOK)
	}

	sessionCookie := sessionCookieFromResponse(t, loginResponse)
	if sessionCookie.Value == "" {
		t.Fatal("session cookie is empty")
	}
	if sessionCookie.Value == "admin" {
		t.Fatal("session cookie stores raw username instead of opaque session token")
	}

	authorizedResponse := doJSONRequest(t, router, http.MethodGet, "/get", nil, sessionCookie)
	if authorizedResponse.Code != http.StatusOK {
		t.Fatalf("authorized request status = %d, want %d", authorizedResponse.Code, http.StatusOK)
	}

	forgedCookie := &http.Cookie{Name: session.CookieName, Value: "admin"}
	forgedResponse := doJSONRequest(t, router, http.MethodGet, "/get", nil, forgedCookie)
	if forgedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("forged cookie status = %d, want %d", forgedResponse.Code, http.StatusUnauthorized)
	}

	logoutResponse := doJSONRequest(t, router, http.MethodPost, "/logout", nil, sessionCookie)
	if logoutResponse.Code != http.StatusOK {
		t.Fatalf("logout status = %d, want %d", logoutResponse.Code, http.StatusOK)
	}

	afterLogoutResponse := doJSONRequest(t, router, http.MethodGet, "/get", nil, sessionCookie)
	if afterLogoutResponse.Code != http.StatusUnauthorized {
		t.Fatalf("request after logout status = %d, want %d", afterLogoutResponse.Code, http.StatusUnauthorized)
	}
}

func TestTaskLifecycleThroughHTTPAPI(t *testing.T) {
	router := newTestRouter()

	doJSONRequest(t, router, http.MethodPost, "/register", models.User{
		Username: "student",
		Password: "1234",
	})

	loginResponse := doJSONRequest(t, router, http.MethodPost, "/login", models.User{
		Username: "student",
		Password: "1234",
	})
	sessionCookie := sessionCookieFromResponse(t, loginResponse)

	addResponse := doJSONRequest(t, router, http.MethodPost, "/add", models.Task{
		Headline: "first task",
		Details:  "write tests",
	}, sessionCookie)
	if addResponse.Code != http.StatusCreated {
		t.Fatalf("add status = %d, want %d", addResponse.Code, http.StatusCreated)
	}

	var created map[string]int64
	if err := json.Unmarshal(addResponse.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode add response: %v", err)
	}
	if created["id"] != 1 {
		t.Fatalf("created id = %d, want 1", created["id"])
	}

	listResponse := doJSONRequest(t, router, http.MethodGet, "/get", nil, sessionCookie)
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
	}

	var activeTasks []models.Task
	if err := json.Unmarshal(listResponse.Body.Bytes(), &activeTasks); err != nil {
		t.Fatalf("decode list response: %v", err)
	}
	if len(activeTasks) != 1 {
		t.Fatalf("active tasks len = %d, want 1", len(activeTasks))
	}

	updateResponse := doJSONRequest(t, router, http.MethodPost, "/update", models.Task{
		ID:       1,
		Headline: "updated task",
		Details:  "ready for review",
	}, sessionCookie)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updateResponse.Code, http.StatusOK)
	}

	getByIDResponse := doJSONRequest(t, router, http.MethodGet, "/get/1", nil, sessionCookie)
	if getByIDResponse.Code != http.StatusOK {
		t.Fatalf("get by id status = %d, want %d", getByIDResponse.Code, http.StatusOK)
	}

	var savedTask models.Task
	if err := json.Unmarshal(getByIDResponse.Body.Bytes(), &savedTask); err != nil {
		t.Fatalf("decode task response: %v", err)
	}
	if savedTask.Headline != "updated task" {
		t.Fatalf("headline = %q, want %q", savedTask.Headline, "updated task")
	}

	resolveResponse := doJSONRequest(t, router, http.MethodPost, "/resolve/1", nil, sessionCookie)
	if resolveResponse.Code != http.StatusOK {
		t.Fatalf("resolve status = %d, want %d", resolveResponse.Code, http.StatusOK)
	}

	activeAfterResolve := doJSONRequest(t, router, http.MethodGet, "/get", nil, sessionCookie)
	if activeAfterResolve.Code != http.StatusOK {
		t.Fatalf("active after resolve status = %d, want %d", activeAfterResolve.Code, http.StatusOK)
	}

	activeTasks = nil
	if err := json.Unmarshal(activeAfterResolve.Body.Bytes(), &activeTasks); err != nil {
		t.Fatalf("decode active tasks after resolve: %v", err)
	}
	if len(activeTasks) != 0 {
		t.Fatalf("active tasks after resolve len = %d, want 0", len(activeTasks))
	}

	archiveResponse := doJSONRequest(t, router, http.MethodGet, "/archive", nil, sessionCookie)
	if archiveResponse.Code != http.StatusOK {
		t.Fatalf("archive status = %d, want %d", archiveResponse.Code, http.StatusOK)
	}

	var archivedTasks []models.Task
	if err := json.Unmarshal(archiveResponse.Body.Bytes(), &archivedTasks); err != nil {
		t.Fatalf("decode archive response: %v", err)
	}
	if len(archivedTasks) != 1 {
		t.Fatalf("archived tasks len = %d, want 1", len(archivedTasks))
	}
	if !archivedTasks[0].Archived || !archivedTasks[0].Completed {
		t.Fatal("resolved task must be archived and completed")
	}
}
