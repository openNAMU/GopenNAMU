package tool

import (
	"os"
	"path/filepath"
)

func Main_init() {
	version_file_path := filepath.Join("..", "data", "version.json")
	if _, err := os.Stat(version_file_path); err == nil {
		// 파일 있음
	} else {
		First_init()
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

func Get_now_version() {

}

func First_init() {

}

func Make_set_json() {

}