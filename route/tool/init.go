package tool

import (
	"os"
	"path/filepath"
)

func Main_init() {
	if Get_now_version() == "" {
		First_init()
	} else {
		now_version := Get_now_version()
		last_version := Get_last_version()

		if now_version != last_version["c_ver"] {
			Make_set_json()
		}
	}

	DB_init_standalone()
}

func Get_last_version() map[string]string {
	version_file_path := filepath.Join("..", "version.json")
	if _, err := os.Stat(version_file_path); err == nil {
		data, err := os.ReadFile(version_file_path)
		if err != nil {
			panic(err)
		}

		version_json := map[string]string{}
    	json.Unmarshal([]byte(data), &version_json)

		return version_json
	} else {
		panic(err)
	}
}

func Get_now_version() string {
	version_file_path := filepath.Join("..", "data", "version.json")
	if _, err := os.Stat(version_file_path); err == nil {
		data, err := os.ReadFile(version_file_path)
		if err != nil {
			panic(err)
		}

		return string(data)
	} else {
		return ""
	}
}

func First_init() {

}

func Make_set_json() {

}