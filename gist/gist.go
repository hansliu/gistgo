package gist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	grequests "github.com/levigross/grequests"
	gjson "github.com/tidwall/gjson"
)

// File is gistfile object
type File struct {
	Name    string `json:"filename"`
	Content string `json:"content"`
}

// Files is array of gistfile object
type Files map[string]*File

func check(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func getToken() (token string) {
	cmd := exec.Command("git", "config", "--get", "gist.token")
	// log.Println("Get gist token")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalln("The gist.token not found, please check your .gitconfig")
	}

	return strings.TrimRight(string(out), "\r\n")
}

func listGist() (resp *grequests.Response) {
	token := getToken()
	api := "https://api.github.com/gists"
	log.Println("ListGist")

	resp, err := grequests.Get(api,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": "token " + token,
				"Accept":        "application/vnd.github.v3+json",
			},
		},
	)
	if err != nil || resp.Ok != true {
		log.Println("Unable to make request: ", err)
		log.Fatalln("StatusCode: ", resp.StatusCode)
	}

	return resp
}

// GetGist is get gist to file by gistID
func GetGist(gistID string) (resp *grequests.Response) {
	token := getToken()
	api := "https://api.github.com/gists/" + gistID
	log.Println("GetGist: ", gistID)

	resp, err := grequests.Get(api,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": "token " + token,
				"Accept":        "application/vnd.github.v3+json",
			},
		},
	)
	if err != nil || resp.Ok != true {
		log.Println("Unable to make request: ", err)
		log.Fatalln("StatusCode: ", resp.StatusCode)
	}

	// create local folder
	filedir := gistID
	_ = os.Mkdir(filedir, 0700)
	// parse resp and donwload gist file to local folder
	results := gjson.Get(resp.String(), "files")
	results.ForEach(func(key, result gjson.Result) bool {
		filename := result.Get("filename").String()
		fileURL := result.Get("raw_url").String()
		downloadGist(filedir+"/"+filename, fileURL)
		return true
	})
	return resp
}

// UploadGist is upload file content to gist
func UploadGist(name string, path string, public bool) (resp *grequests.Response) {
	token := getToken()
	api := "https://api.github.com/gists"
	log.Println("UploadGist:", path)

	content, err := ioutil.ReadFile(path)
	check(err)

	// get filename
	filename := filepath.Base(path)

	files := make(Files)
	files[filename] = &File{
		Content: string(content),
	}
	if name != "" {
		files[filename].Name = name
	}

	obj := make(map[string]interface{})
	obj["files"] = files
	obj["public"] = public
	jsonObj, err := json.Marshal(obj)
	check(err)

	// This will upload the file as a multipart mime request
	resp, err = grequests.Post(api,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": "token " + token,
				"Accept":        "application/vnd.github.v3+json",
			},
			JSON: jsonObj,
		},
	)
	if err != nil || resp.Ok != true {
		log.Println("Unable to make request", err)
		log.Fatalln("StatusCode: ", resp.StatusCode)
	}
	fmt.Println("Upload successful:", gjson.Get(resp.String(), "html_url"))
	return resp
}

func downloadGist(path string, fileURL string) {
	token := getToken()
	log.Println("downloadGist:", fileURL)

	resp, err := grequests.Get(fileURL,
		&grequests.RequestOptions{
			Headers: map[string]string{
				"Authorization": "token " + token,
				"Accept":        "application/vnd.github.v3+json",
			},
		},
	)
	if err != nil || resp.Ok != true {
		log.Fatalln("Unable to make request", err)
		log.Fatalln("StatusCode: ", resp.StatusCode)
	}
	if err := resp.DownloadToFile(path); err != nil {
		log.Fatalln("Unable to download to file: ", err)
	}
}
