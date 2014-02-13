package main

import (
	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Extensions: []string{".html"},
	}))

	m.Get("/", func(render render.Render, log *log.Logger) {
		render.HTML(200, "index", nil)
	})

	m.Post("/upload", func(r *http.Request) (int, string) {
		log.Println("parsing form")
		err := r.ParseMultipartForm(100000)
		if err != nil {
			return http.StatusInternalServerError, err.Error()
		}

		files := r.MultipartForm.File["files"]
		for i, _ := range files {
			log.Println("getting handle to file")
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}

			log.Println("creating destination file")
			dst, err := os.Create("./uploads/" + files[i].Filename)
			defer dst.Close()
			if err != nil {
				return http.StatusInternalServerError, err.Error()
			}

			log.Println("copying the uploaded file to the destination file")
			if _, err := io.Copy(dst, file); err != nil {
				return http.StatusInternalServerError, err.Error()
			}
		}

		return 200, "ok"
	})

	m.Run()
}
