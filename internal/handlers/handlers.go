package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func GetHTML(w http.ResponseWriter, r *http.Request) {
	dataHTML, err := os.ReadFile("../index.html")
	if err != nil {
		http.Error(w, "couldn't read the html file", http.StatusNotFound)
		return
	}

	w.Header().Set("Contyyent-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(dataHTML)
}

func PostHTML(w http.ResponseWriter, r *http.Request) {
	// выделение памяти под содержимое файла
	r.ParseMultipartForm(10 << 20)

	// получение данных файла и метаданных
	file, h, err := r.FormFile("myFile")
	if err != nil {
		http.Error(w, "error receiving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "error read file data", http.StatusInternalServerError)
		return
	}

	convertedData, err := service.CodeDetection(string(data))
	if err != nil {
		http.Error(w, "error convertion data", http.StatusInternalServerError)
		return
	}

	// проверка наличия директории uploads, если ее нет, то создается
	_, err = os.Stat("../uploads")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir("../uploads", 0755)
			if err != nil {
				http.Error(w, "error creating dyrectory", http.StatusInternalServerError)
				return
			}
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	// открытие директории через root для сохранения файла
	root, err := os.OpenRoot("../uploads")
	if err != nil {
		http.Error(w, "open root directory error", http.StatusInternalServerError)
		return
	}
	defer root.Close()

	// time.Now().UTC().String()
	dst, err := root.Create(h.Filename)
	if err != nil {
		http.Error(w, "error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil || written == 0 {
		http.Error(w, "error write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Accept-Language", "ru,en;q=0.9")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(convertedData))
}
