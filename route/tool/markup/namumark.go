package markup

import (
	"database/sql"
	"log"
	"strings"
)

type namumark struct {
	db   *sql.DB
	data map[string]string

	render_data string
}

func Namumark_new(db *sql.DB, data map[string]string) *namumark {
	data_string := data["data"]
	data_string = "\n" + data_string + "\n"
	data_string = strings.ReplaceAll(data_string, "\r", "")

	return &namumark{
		db,
		data,

		data_string,
	}
}

func (class *namumark) line_render(line string) string {
	return line
}

func (class namumark) main() map[string]any {
	lines := strings.Split(class.render_data, "\n")
	for i, line := range lines {
		lines[i] = class.line_render(line)
	}

	class.render_data = strings.Join(lines, "\n")
	log.Default().Println(class.render_data)

	class.data["data"] = class.render_data

	render_data_class := Macromark_new(class.db, class.data)
	render_data := render_data_class.main()

	return render_data
}
