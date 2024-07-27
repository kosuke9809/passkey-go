package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"passkey-auth/user"
)

type RequestBody struct {
	Username string `json:"username"`
}

func parseJSONBody(r *http.Request) (RequestBody, error) {
	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	return body, err
}

func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := parseJSONBody(r)
	if err != nil {
		log.Printf("Error parsing JSON in BeginRegistration: %v", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if body.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	u := user.New(body.Username)
	user.Store(u)

	options, sessionData, err := Web.BeginRegistration(u)
	if err != nil {
		log.Printf("Error in BeginRegistration: %v", err)
		http.Error(w, fmt.Sprintf("Registration initialization failed: %v", err), http.StatusInternalServerError)
		return
	}
	u.SetSessionData(*sessionData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body in FinishRegistration: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received registration data: %s", string(bodyBytes))

	var body struct {
		Username string `json:"username"`
		ID       string `json:"id"`
		RawID    string `json:"rawId"`
		Type     string `json:"type"`
		Response struct {
			AttestationObject string `json:"attestationObject"`
			ClientDataJSON    string `json:"clientDataJSON"`
		} `json:"response"`
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		log.Printf("Error parsing JSON in FinishRegistration: %v", err)
		http.Error(w, fmt.Sprintf("Parse error for Registration: %v", err), http.StatusBadRequest)
		return
	}

	if body.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	u, err := user.Get(body.Username)
	if err != nil {
		log.Printf("Error getting user in FinishRegistration: %v", err)
		http.Error(w, fmt.Sprintf("Error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	sessionData := u.SessionData()
	log.Printf("User session data: %+v", sessionData)

	// Create a new http.Request with the parsed body data
	newBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newRequest, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewReader(newBody))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newRequest.Header = r.Header

	credential, err := Web.FinishRegistration(u, sessionData, newRequest)
	if err != nil {
		log.Printf("Error in FinishRegistration: %v", err)
		http.Error(w, fmt.Sprintf("Registration failed: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Credential created: %+v", credential)
	u.AddCredential(credential)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successful"))
}
func BeginLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := parseJSONBody(r)
	if err != nil {
		log.Printf("Error parsing JSON in BeginLogin: %v", err)
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if body.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	u, err := user.Get(body.Username)
	if err != nil {
		log.Printf("Error getting user in BeginLogin: %v", err)
		http.Error(w, fmt.Sprintf("Error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	options, sessionData, err := Web.BeginLogin(u)
	if err != nil {
		log.Printf("Error in BeginLogin: %v", err)
		http.Error(w, fmt.Sprintf("Login initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	u.SetSessionData(*sessionData)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

func FinishLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body in FinishLogin: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	log.Printf("Received login data: %s", string(bodyBytes))

	var body struct {
		Username string `json:"username"`
		ID       string `json:"id"`
		RawID    string `json:"rawId"`
		Type     string `json:"type"`
		Response struct {
			AuthenticatorData string `json:"authenticatorData"`
			ClientDataJSON    string `json:"clientDataJSON"`
			Signature         string `json:"signature"`
			UserHandle        string `json:"userHandle"`
		} `json:"response"`
	}
	if err := json.Unmarshal(bodyBytes, &body); err != nil {
		log.Printf("Error parsing JSON in FinishLogin: %v", err)
		http.Error(w, fmt.Sprintf("Parse error for Login: %v", err), http.StatusBadRequest)
		return
	}

	if body.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	u, err := user.Get(body.Username)
	if err != nil {
		log.Printf("Error getting user in FinishLogin: %v", err)
		http.Error(w, fmt.Sprintf("Error getting user: %v", err), http.StatusInternalServerError)
		return
	}

	sessionData := u.SessionData()
	log.Printf("User session data: %+v", sessionData)

	// Create a new http.Request with the parsed body data
	newBody, err := json.Marshal(body)
	if err != nil {
		log.Printf("Error marshaling body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newRequest, err := http.NewRequest(r.Method, r.URL.String(), bytes.NewReader(newBody))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newRequest.Header = r.Header

	_, err = Web.FinishLogin(u, sessionData, newRequest)
	if err != nil {
		log.Printf("Error in FinishLogin: %v", err)
		http.Error(w, fmt.Sprintf("Login failed: %v", err), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
