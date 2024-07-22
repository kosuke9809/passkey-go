package auth

import "net/http"

func BeginRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}
	u := user.New(username)

}
