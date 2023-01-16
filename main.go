package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type File struct {
	Name string `json:"name"`
}


func main() {
	os.Mkdir("public", 0666)
	PORT := ":8080"

	Wg := sync.WaitGroup{}
	Mx := sync.Mutex{}

	for i := 0; i < 10; i++ {
		Wg.Add(1)
		go AddFile(&Wg, &Mx)
	}

	Wg.Wait()

	files := GetFiles()
	RenderFiles(&files)
	http.ListenAndServe(PORT, nil)
}

func AddFile(wg *sync.WaitGroup, mx *sync.Mutex) {

	fileName := UUID()
	CreateFile(fileName)

	mx.Lock()
	UpDataDB(fileName)
	mx.Unlock()

	wg.Done()
}

func CreateFile(fileName string)  {
	f, err := os.Create(fmt.Sprintf("./public/%s.txt", fileName))

	if err != nil {
		panic(err)
	}

	f.Write([]byte("siusbuviavsiua vsaiuvavsinoasvas avasvasv \n\n\n\n casuicbasu ccsauisacnasccas cbasucbauicbsauic asc ascasbucbsacu"))
	fmt.Println(f.Name())
	f.Close()
}

func UpDataDB(fileName string) {
	files := GetFiles()
	file := File{Name: fileName}
	files = append(files, file)
	data, err := json.Marshal(files)

	if err != nil {
		panic(err)
	}

	os.WriteFile("./files.json", data, 0666)
}

func RenderFiles(data *[]File) {
	var files []File = *data;

	for i := 0; i < len(files); i++ {
		var fileName string = files[i].Name
		http.HandleFunc(fmt.Sprintf("/public/%s", fileName), func(w http.ResponseWriter, r *http.Request) {
			bytes, err := os.ReadFile(fmt.Sprintf("./public/%s.txt", fileName))

			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		})
	}

	GetAllFilesName(data)
}


func GetAllFilesName(data *[]File) {
	var files []File = *data
	http.HandleFunc("/public", func(w http.ResponseWriter, r *http.Request) {

		bytes, err := json.Marshal(files)

		if err != nil {
			panic(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	})
}



func UUID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}


func GetFiles() []File {
	bytes, err := os.ReadFile("./files.json")
	if err != nil {
		panic(err)
	}
	var files []File
	json.Unmarshal(bytes, &files)
	return files
}
