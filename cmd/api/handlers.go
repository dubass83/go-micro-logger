package api

import (
	"net/http"
)

// api Handler
// func (s *Server) Authenticate(w http.ResponseWriter, r *http.Request) {
// 	var requestPayload struct {
// 		Email    string `json:"email"`
// 		Password string `json:"password"`
// 	}

// 	err := readJSON(w, r, &requestPayload)
// 	if err != nil {
// 		errorJSON(w, err, http.StatusBadRequest)
// 		return
// 	}
// 	// get user from database by email
// 	user, err := s.Db.GetUserByEmail(context.Background(), requestPayload.Email)
// 	if err != nil {
// 		errorJSON(w, errors.New("invalid username or password"), http.StatusBadRequest)
// 		return
// 	}
// 	// compare password and hash from database
// 	err = util.CheckPassword(requestPayload.Password, user.Password)
// 	if err != nil {
// 		errorJSON(w, errors.New("invalid username or password"), http.StatusBadRequest)
// 		return
// 	}

// 	resultStr := fmt.Sprintf("password is valid for user: %s", user.Email)

// 	payload := &jsonResponse{
// 		Error:   false,
// 		Massage: resultStr,
// 		Data:    user,
// 	}
// 	log.Debug().Msgf("payload: %+v", payload)
// 	_ = writeJSON(w, http.StatusAccepted, payload)
// }

// Test api Handler
func Test(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Massage: "Hello from Broker!",
	}

	_ = writeJSON(w, http.StatusAccepted, payload)
}
