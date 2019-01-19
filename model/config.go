package model

type Config struct {
	Database struct {
		User     string `json:"user"`
		Password string `json:"password"`
		DBName string `json:"dbName"`
	} `json:"database"`
}