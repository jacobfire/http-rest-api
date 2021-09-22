package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jacobfire/http-rest-api/app/cache"
	"github.com/jacobfire/http-rest-api/app/store"
	"github.com/jacobfire/http-rest-api/configs"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type FileData struct {
	Name string
	CreationDate string `json:"CreatedAt,omitempty"`
	//Status int `json:"Status,omitempty"`
	Status int `json:"Status"`
}

type ByDate []FileData

func (a ByDate) Len() int {
	return len(a)
}

func (a ByDate) Less(i, j int) bool {
	return a[i].CreationDate > a[j].CreationDate
}

func (a ByDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type APIServer struct {
	logger *logrus.Logger
	router *mux.Router
	store *store.Store
	CacheManager cache.Cache
}

func New() *APIServer {
	return &APIServer{
		logger: logrus.New(),
		router: mux.NewRouter(),
		CacheManager: *cache.New(10 * time.Minute, 25 * time.Minute),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Starting server...")

	server := s.configureServer()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return nil
}

// Configure Store
func (s *APIServer) configureStore() error {
	config := configs.NewConfig()
	st := store.New(config.Store)

	log.Println("Configuring store")
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st

	return nil
}

// Handle Hello function
func (s * APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) configureLogger() error {
	config := configs.NewConfig()
	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureServer() http.Server {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/analyse", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Analyse"))
	})
	s.router.HandleFunc("/catalog/{category}/{id:[0-9]+}", s.categoryHandler())

	s.router.Handle("/migrations", middleware(http.HandlerFunc(s.fetchMigration))).Methods(http.MethodGet)
	s.router.Handle("/migrations", http.HandlerFunc(s.CreateMigration)).Methods(http.MethodPost)
	s.router.Handle("/filesContent", http.HandlerFunc(s.filesContent)).Methods(http.MethodGet)
	//s.router.Handle("/signin", http.HandlerFunc(s.Signin)).Methods(http.MethodPost)

	config := configs.NewConfig()
	server := http.Server {
		Addr: config.BindAddr,
		Handler: s.router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server
}

func (s * APIServer) categoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		io.WriteString(w, "Article")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v ID:%v \n", vars["category"], vars["id"])
	}
}

func (s *APIServer) Migrate() error {
	config := configs.NewConfig()
	log.Println(config.Store.DatabaseURL)
	dirPath, err := os.Getwd()
	dirPath = dirPath + "/" + "migrations/"

	m, err := migrate.New(
		"file:" + dirPath,
		config.Store.DatabaseURL,
		)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}

func (s *APIServer) CreateMigration(w http.ResponseWriter, r *http.Request) {
	/**
	1. get body from request
	2. unmarshal JSON
	3. validate data
	3.1 validate if empty field
	3.2 bad request in response
	3.3 if marshal with error we need to return "unprocessable entity" status
	4. create file according to name in JSON
	4.1 create func for creation of a file

	//TODO
	5. check already existing files
	6. write content to file
	 */

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		s.sendResponse(w, "Data can not be parsed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	fmt.Println(string(body))

	type Migration struct {
		Name string `json:"name"`
	}

	input := Migration {}
	err = json.Unmarshal(body, &input)

	if err != nil {
		s.sendResponse(w, "Data can not be parsed", http.StatusBadRequest)
		log.Println(err)
		return
	}

	log.Println(input)

	if len(input.Name) == 0 {
		s.sendResponse(w, "ERROR: empty name in request", http.StatusUnprocessableEntity)
		return
	}

	prefixes := []string {
		"up",
		"down",
	}

	for _, prefix := range prefixes {
		if _, err := s.createMigrationFile(input.Name, prefix); err != nil{
			s.sendResponse(w, "File not created", 0)
			return
		}
	}

	s.sendResponse(w, "Successfully processed", http.StatusCreated)
}

//prepare response
func (s *APIServer) sendResponse(w http.ResponseWriter, message string, status int) http.ResponseWriter {
	w.WriteHeader(status)
	if status == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Fatal(err)
	}

	return w
}

// Create file
func (s *APIServer) createMigrationFile(name string, prefix string) (*os.File, error) {
	_, e := os.Stat(name)
	if e == nil {
		return nil, os.ErrExist
	}
	rootFolder, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	dateAsString := currentTime.Format("20060102150405")
	fileName := dateAsString + "_" + "create_" + strings.ToLower(name) + "." + strings.ToLower(prefix) + ".sql"

	fullPath := rootFolder + "/migrations/" + fileName
	return os.Create(fullPath)
}

func (s *APIServer) fetchMigration(w http.ResponseWriter, r *http.Request) {
	rootFolder, err := os.Getwd()
	if err != nil {
		s.sendResponse(w, "files not found", 0)
		return
	}
	migrationFolder := rootFolder + "/migrations"
	files, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		s.sendResponse(w, "files not found", 0)
		return
	}

	token, ok := r.Context().Value("authtoken").(string)
	if !ok {
		s.sendResponse(w, "unautharized", http.StatusUnauthorized)
		return
	}

	//var preparedData [][]FileData
	var fileNames []FileData
	for index, file := range files {

		fileData := FileData {
			Name: file.Name(),
			CreationDate: file.ModTime().Format("2006-01-02 15:04:05"),
		}

		if index == 2 {
			fileData.Status = index
		}
		fileNames = append(fileNames, fileData)
	}
	//preparedData = fileNames
	sort.Sort(ByDate(fileNames))
	resultString, e := json.Marshal(fileNames)
	if e != nil {
		s.sendResponse(w, "Internal error", 0)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Bearer", token)
	s.sendResponse(w, string(resultString), http.StatusOK)
}

func (s *APIServer) filesContent(w http.ResponseWriter, r *http.Request) {
	rootFolder, err := os.Getwd()
	if err != nil {
		s.sendResponse(w, "searching for files has been stopped", 0)
		return
	}
	migrationFolder := rootFolder + "/migrations"
	files, err := ioutil.ReadDir(migrationFolder)
	if err != nil {
		s.sendResponse(w, "files not found", 0)
		return
	}
	
	readFiles := func() <-chan []byte {
		contentChannel := make(chan []byte)
		for _, file := range files {
			var  translating = false
			//var cachedData interface{}
			var cachedFileContent []byte
			cachedData, unpacking := s.CacheManager.Get(file.Name())

			fmt.Println("UNPACKING ", unpacking)
			fmt.Println("DATA ", cachedData)
			if  unpacking {
				cachedFileContent, translating = cachedData.([]byte)
				fmt.Println("DATA CONVERTED ", string(cachedFileContent))

				fmt.Println("Translating => ", unpacking)
				if translating {
					go func(t []byte) {
						contentChannel <-t
					}(cachedFileContent)
				}
			}
			if !unpacking && !translating {
				fmt.Println("goroutine works")
				go func(fileName string) {
					content, err := ioutil.ReadFile(migrationFolder + "/" + fileName)
					if err != nil {
						log.Fatal(err)
					}
					s.CacheManager.Set(fileName, content, 10 * time.Minute)
					fmt.Println("FILENAME ", fileName)
					contentChannel <- content
				} (file.Name())
			}
		}

		return contentChannel
	}

	readContent := ""
	filesContent := readFiles()
	for i := 0; i < len(files); i++ {
		readContent += string(<-filesContent)
	}

	s.sendResponse(w, readContent, http.StatusOK)
}

//func (s *APIServer) Signin(w http.ResponseWriter, r *http.Request) {
//	token, err := GenerateJWT("test@test.com", "admin")
//	if err != nil {
//		s.sendResponse(w, "cant generate token", http.StatusUnauthorized)
//	}
//	s.sendResponse(w, token, http.StatusOK)
//
//}


func Paginate(total, page int, size int) (int, int) {
	page = page * size
	size = page + size
	if page > total {
		page = total
	}
	if size > total {
		size = total
	}
	return page, size
}

func Sum(a, b int) int {
	return a + b
}

