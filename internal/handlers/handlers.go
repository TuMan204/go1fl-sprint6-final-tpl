package handlers

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func GetHTML(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat("../index.html")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.Error(w, "file not exists", http.StatusNotFound)
		}
	}
	dataHTML, err := os.ReadFile("../index.html")
	if err != nil {
		http.Error(w, "couldn't read the html file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(dataHTML)
}

func PostHTML(w http.ResponseWriter, r *http.Request) {
	// выделение памяти под содержимое файла
	r.ParseMultipartForm(10 << 20)

	// получение данных файла и метаданных
	file, headers, err := r.FormFile("myFile")
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
	// file.Seek(0, 0)
	// if err != nil {
	// 	http.Error(w, "error read file data", http.StatusInternalServerError)
	// 	return
	// }

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
				http.Error(w, "error creating directory", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	// открытие директории через root для сохранения файла
	root, err := os.OpenRoot("../uploads")
	if err != nil {
		http.Error(w, "open root directory error", http.StatusInternalServerError)
		return
	}
	defer root.Close()

	// time.Now().UTC().String()
	dst, err := root.Create(time.Now().UTC().Format("20060102_150405") + filepath.Ext(headers.Filename))
	if err != nil {
		http.Error(w, "error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// _, err = io.Copy(dst, file)
	// if err != nil {
	// 	http.Error(w, "error write file", http.StatusInternalServerError)
	// 	return
	// }
	_, err = dst.WriteString(convertedData)
	if err != nil {
		http.Error(w, "error write file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Accept-Language", "ru")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(convertedData))
}
