package web

import (
	"boards-merger/internal/core"
	"boards-merger/internal/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func StartWebServer(port string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("POST /processPath", handleProcessPath)

	// Setup the file server based on the executable path to avoid relative path and CWD problems
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to setup static file server: %v", err.Error())
	}

	basePath := filepath.Dir(execPath)
	staticPath := filepath.Join(basePath, "../internal/web")
	mux.Handle("/static/", http.FileServer(http.Dir(staticPath)))

	return http.ListenAndServe(":"+port, mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl := GetTemplate()
	if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Error generating html.index", http.StatusInternalServerError)
	}
}

func handleProcessPath(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Error  string
		Result *model.BoardsInfo
	}

	r.ParseForm()
	path := r.FormValue("path")
	recursive := r.FormValue("recursive") == "on"
	depth := 10
	fmt.Sscanf(r.FormValue("depth"), "%d", &depth)

	defer func() {
		tmpl := GetTemplate()
		if err := tmpl.ExecuteTemplate(w, "boards_table.html", data); err != nil {
			http.Error(w, "Error generating boards result", http.StatusInternalServerError)
		}
	}()

	jsonList, err := core.ReadDirectory(path, recursive, depth)
	if err != nil {
		data.Error = err.Error()
		return
	}

	boards, err := core.ProcessJsonFiles(jsonList)
	if err != nil {
		data.Error = err.Error()
		return
	}

	data.Result = boards
}
