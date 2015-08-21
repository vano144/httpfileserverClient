package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
}

func showError(err error) {
	if err != nil {
		log.Println("Internal Error", err)
	}
}

func show(Url, name, userPassword string, client *http.Client) error {
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(name, userPassword)
	req.Header.Set("Accept", "application/json")
	resp, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return err2
	}
	defer resp.Body.Close()
	res := []FileInfo{}
	err3 := json.Unmarshal(data, &res)
	if err3 != nil {
		return err3
	}
	for i, f := range res {
		fmt.Println(i, " ", f.Name, " size: ", f.Size)
	}
	return nil
}

func delete(Url, name, userPassword string, client *http.Client) error {
	fmt.Println("input name of file")
	var nameOfFile string
	fmt.Scanln(&nameOfFile)
	req, err := http.NewRequest("GET", Url, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(name, userPassword)
	req.Header.Set("Action", "delete "+nameOfFile)
	_, err = client.Do(req)
	if err != nil {
		return err
	}
	return nil
}

func uploadRequest(url string, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(file.Name(), filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", url, body)
}
func upload(Url, userName, userPassword string, client *http.Client) error {
	fmt.Println("input path to file")
	var path string
	fmt.Scanln(&path)
	fmt.Println("input name, which you see on server")
	var name string
	fmt.Scanln(&name)
	req, err := uploadRequest(Url, path)
	if err != nil {
		return err
	}
	req.Header.Set("Action", "upload "+name)
	req.SetBasicAuth(userName, userPassword)
	_, err1 := client.Do(req)
	if err1 != nil {
		return err1
	}
	return nil
}

func download(Url, name, userPassword string, client *http.Client) error {
	fmt.Println("input name of file")
	var nameOfFile string
	fmt.Scanln(&nameOfFile)
	req, err := http.NewRequest("GET", Url+"usersStorage/"+name+"/"+nameOfFile, nil)
	req.SetBasicAuth(name, userPassword)
	if err == nil {
		resp, err1 := client.Do(req)
		if err1 != nil {
			return err1
		}
		out, err2 := os.Create(nameOfFile)
		if err2 != nil {
			return err2
		}
		defer out.Close()
		_, err3 := io.Copy(out, resp.Body)
		if err3 != nil {
			return err3
		}
		return nil
	}
	return err
}

func main() {
	var name string
	var userPassword string
	var input string
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	Url := "https://127.0.0.1" + *port + "/cloud/"
	fmt.Println("URL:>", Url)
	client := &http.Client{Transport: tr}
	for {
		var input string
		fmt.Println("input name of user, input exit to exit")
		fmt.Scanln(&input)
		if input == "exit" {
			fmt.Println("Goodbye")
			os.Exit(0)
		}
		name = input
		fmt.Println("input password of user")
		fmt.Scanln(&input)
		userPassword = input
		req, err := http.NewRequest("GET", Url, nil)
		req.SetBasicAuth(name, userPassword)
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			fmt.Println("Successful auth")
			break
		}
		log.Println(err, "Bad Auth")
	}
	for {
		fmt.Println("input command: show, delete,download,upload,exit")
		fmt.Scanln(&input)
		switch {
		case input == "show":
			err := show(Url, name, userPassword, client)
			showError(err)
		case input == "delete":
			err := delete(Url, name, userPassword, client)
			showError(err)
		case input == "download":
			err := download(Url, name, userPassword, client)
			showError(err)
		case input == "upload":
			err := upload(Url, name, userPassword, client)
			showError(err)
		case input == "exit":
			fmt.Println("Goodbye")
			os.Exit(0)
		default:
			fmt.Println("unknown command")
		}
	}
}
