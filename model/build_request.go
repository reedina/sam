package model

//CreateBuild (POST)
func CreateBuild(buildJSON string) error {

	var ID int

	err := db.QueryRow(
		"INSERT INTO sam_build_request(request) VALUES($1) RETURNING id", buildJSON).Scan(&ID)

	if err != nil {
		return err
	}

	return nil
}
