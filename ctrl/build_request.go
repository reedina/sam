package ctrl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

	"github.com/gorilla/mux"
	"github.com/reedina/sam/model"
)

//GetBuildRequestProfile (GET)
func GetBuildRequestProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email, err := mail.ParseAddress(vars["email"])

	//Get User and Team
	getUser := fmt.Sprintf("http://spm:4040/api/user/%s", email)
	user, err := getClientJSON(getUser)

	if err != nil {
		fmt.Println(err)

	}

	//Get Projects
	teamID := getClientTeamID(user)
	//teamID := 0
	getProjects := fmt.Sprintf("http://spm:4040/api/projects/team/id/%d", teamID)
	projects, err := getClientJSON(getProjects)

	if err != nil {
		fmt.Println(err)
	}

	//Get Environments
	getEnvironments := "http://sbm:4050/api/environments"
	environments, err := getClientJSON(getEnvironments)

	if err != nil {
		fmt.Println(err)
	}

	// if TeamID is 0 then override environment values and make empty
	if teamID == 0 {
		environments = "[]"
	}
	// Return UserProfile
	userProfile := `{"user": ` + user + `, "projects": ` + projects + `, "environments": ` + environments + `}`

	respondWithJSONBytes(w, http.StatusOK, []byte(userProfile))
}

// CreateBuild (POST)
func CreateBuild(w http.ResponseWriter, r *http.Request) {
	jsonData, err := ioutil.ReadAll(r.Body)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := model.CreateBuild(string(jsonData)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSONBytes(w, http.StatusOK, []byte(jsonData))
}

func getClientJSON(url string) (string, error) {

	// Build request
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", err
	}

	// Create Client object
	client := &http.Client{}

	// Execute Request
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	jsonData, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	// convert json bytes to string
	return string(jsonData), nil

}

func getClientTeamID(userJSON string) int {

	var userData map[string]interface{}

	if err := json.Unmarshal([]byte(userJSON), &userData); err != nil {
		return 0
	}
	team, ok := userData["team"].(map[string]interface{})

	if !ok {
		return 0
	}

	return int(team["id"].(float64))
}
