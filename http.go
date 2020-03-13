package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/colinmarc/hdfs"

	log "github.com/sirupsen/logrus"

	humanize "github.com/dustin/go-humanize"
)

var (
	tmplList *template.Template
	tmplInfo *template.Template
)

type listPageData struct {
	Error       string
	Message     string
	BasePath    string
	Breadcrumbs []string
	Parts       []string
	Files       []os.FileInfo
	Humanize    func(uint64) string
}

type infoPageData struct {
	Error    error
	FsInfo   hdfs.FsInfo
	Humanize func(uint64) string
}

func init() {
	tmplList = template.Must(template.ParseFiles("templates/list.html"))
	tmplInfo = template.Must(template.ParseFiles("templates/info.html"))
}

func browse(w http.ResponseWriter, r *http.Request) {
	data := listPageData{
		Humanize: humanize.Bytes,
	}

	queries := r.URL.Query()
	data.Error = queries.Get("err")
	data.Message = queries.Get("msg")

	path := strings.TrimPrefix(r.URL.Path, "/browse")

	// Check the path. If it is a file start a download, else return the dir content.
	fi, err := client.Stat(path)
	if err != nil {
		log.Error(err)
		data.Error = "Error while cheking the path on HDFS: " + err.Error()
		tmplList.Execute(w, data)
		return
	}

	// Serve the file (content).
	if !fi.IsDir() {
		fr, err := client.Open(path)
		if err != nil {
			log.Error(err)
			data.Error = "Error while cheking the path on HDFS: " + err.Error()
			tmplList.Execute(w, data) // TODO: better error msg.
			return
		}
		defer fr.Close()

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))

		http.ServeContent(w, r, fi.Name(), fi.ModTime(), fr)

		return
	}

	files, err := client.ReadDir(path)
	if err != nil {
		log.Error(err)
		data.Error = "Error while listing the directory: " + err.Error()
		tmplList.Execute(w, data) // TODO: better error msg.
		return
	}

	data.Files = files

	data.BasePath = strings.Trim(path, "/")
	data.Parts = strings.Split(data.BasePath, "/")

	// Build strings for the breadcrumb elements.
	var breadcrumbs []string
	for _, part := range data.Parts {
		bcLen := len(breadcrumbs)
		if bcLen > 0 {
			breadcrumbs = append(breadcrumbs, fmt.Sprintf("%s/%s", breadcrumbs[bcLen-1], part))
		} else {
			breadcrumbs = append(breadcrumbs, part)
		}
	}
	data.Breadcrumbs = breadcrumbs

	tmplList.Execute(w, data)
}

func upload(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	path = strings.Trim(path, "/")

	err := r.ParseMultipartForm(1 << 20)
	if err != nil {
		log.Error(err)
		redirURL := fmt.Sprintf("/browse/%s?err=Error while parsing form: %s", path, err)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)
		return
	}

	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		log.Error(err)
		redirURL := fmt.Sprintf("/browse/%s?err=Error while reading file: %s", path, err)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)
		return
	}
	defer file.Close()

	filepath := fmt.Sprintf("/%s/%s", path, header.Filename)

	fw, err := client.Create(filepath)
	if err != nil {
		log.Error(err)
		redirURL := fmt.Sprintf("/browse/%s?err=Error while creating %s on HDFS: %s", path, filepath, err)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)
		return
	}
	defer fw.Close()

	_, err = io.Copy(fw, file)
	if err != nil {
		log.Error(err)
		redirURL := fmt.Sprintf("/browse/%s?err=Error uploading %s to HDFS: %s", path, filepath, err)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)
		return
	}

	redirURL := fmt.Sprintf("/browse/%s?msg=Uploaded to %s", path, filepath)

	http.Redirect(w, r, redirURL, http.StatusSeeOther)
}

func delete(w http.ResponseWriter, r *http.Request) {
	path := r.FormValue("path")
	path = strings.Trim(path, "/")

	lastIdx := strings.LastIndex(path, "/")
	redirPath := path[:lastIdx]

	// Do not allow to delete root.
	if path == "" {
		log.Error("Trying to delete root? LOL")
		redirURL := fmt.Sprintf("/browse/%s?err=Deleting root is not permitted!", redirPath)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)

		return
	}

	err := client.Remove("/" + path)
	if err != nil {
		log.Error(err)
		redirURL := fmt.Sprintf("/browse/%s?err=Cannot delete /%s: %s", redirPath, path, err)
		http.Redirect(w, r, redirURL, http.StatusSeeOther)
		return
	}

	redirURL := fmt.Sprintf("/browse/%s?msg=Deleted: /%s", redirPath, path)

	http.Redirect(w, r, redirURL, http.StatusSeeOther)
}

func sysInfo(w http.ResponseWriter, r *http.Request) {
	data := infoPageData{
		Humanize: humanize.Bytes,
	}

	data.FsInfo, data.Error = client.StatFs()

	tmplInfo.Execute(w, data)
}
