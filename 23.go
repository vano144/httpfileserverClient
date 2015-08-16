package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

type FileInfo struct {
	Name string `json:"Name"`
	Size int64  `json:"Size"`
}
type mytype []map[string]string

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	port := flag.String("port", ":9111", "port in server")
	flag.Parse()
	Url := "https://127.0.0.1" + *port + "/cloud/"
	fmt.Println("URL:>", Url)
	client := &http.Client{Transport: tr}
	var name string
	var userPassword string
	var input string
	for {
		var input string
		fmt.Println("input name of user, input exit to exit")
		fmt.Scanln(&input)
		if input == "exit" {
			os.Exit(1)
		}
		name = input
		fmt.Println("input password of user")
		fmt.Scanln(&input)
		userPassword = input
		req, err := http.NewRequest("GET", Url, nil)
		req.SetBasicAuth(name, userPassword)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error : %s", err)
		} else {
			if resp.StatusCode == 200 {
				break
			}
			fmt.Println("Problem with auth")
		}
	}
	fmt.Println(name, ":", userPassword)
	for {
		fmt.Println("input command: show, delete,download,upload,exit")
		fmt.Scanln(&input)
		switch true {
		case input == "show":
			req, err := http.NewRequest("GET", Url, nil)
			if err != nil {
				fmt.Printf("Error : %s", err)
				continue
			}
			req.SetBasicAuth(name, userPassword)
			req.Header.Set("Accept", "application/json")
			resp, err1 := client.Do(req)
			if err1 != nil {
				fmt.Printf("Error : %s", err1)
				continue
			}
			data, err2 := ioutil.ReadAll(resp.Body)
			if err2 != nil {
				fmt.Printf("Error : %s", err2)
				continue
			}
			defer resp.Body.Close()
			res := &[]FileInfo{}
			err3 := json.Unmarshal(data, &res)
			if err3 != nil {
				fmt.Printf("Error : %s", err3)
				continue
			}
			for i, f := range *res {
				fmt.Println(i, " ", f.Name, " size: ", f.Size)
			}
		case input == "delete":
			fmt.Println("input name of file")
			var nameOfFile string
			fmt.Scanln(&nameOfFile)
			req, err := http.NewRequest("GET", Url, nil)
			if err != nil {
				fmt.Printf("Error : %s", err)
				continue
			}
			req.SetBasicAuth(name, userPassword)
			req.Header.Set("Action", "delete "+nameOfFile)
			_, err = client.Do(req)
			if err != nil {
				fmt.Printf("Error : %s", err)
				continue
			}
		case input == "upload":
			fmt.Println("input name of file")
			var nameOfFile string
			fmt.Scanln(&nameOfFile)
			req, err := http.NewRequest("GET", Url+"/usersStorage/"+name+"/"+nameOfFile, nil)
			req.SetBasicAuth(name, userPassword)
			if err == nil {
				resp, err := client.Do(req)
				if err != nil {
					fmt.Printf("Error : %s", err)
					continue
				}
				data, err2 := ioutil.ReadAll(resp.Body)
				if err2 != nil {
					fmt.Printf("Error : %s", err2)
					continue
				}
				defer resp.Body.Close()
				fmt.Println(string(data))
				out, err1 := os.Create(nameOfFile)
				if err1 != nil {
					fmt.Printf("Error : %s", err)
					continue
				}
				defer out.Close()
				_, err = io.Copy(out, resp.Body)
				if err != nil {
					fmt.Printf("Error : %s", err)
					continue
				}
			}
		case input == "download":
			//
		case input == "exit":
			os.Exit(1)
		default:
			fmt.Println("unknown command")
		}
	}

}
