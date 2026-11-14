package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"backend/internal/help"
	helpapi "backend/internal/help/api"
	helpmysql "backend/internal/help/mysql"

	"backend/internal/mother"
	motherapi "backend/internal/mother/api"
	mothermysql "backend/internal/mother/mysql"

	"backend/internal/child"
	childapi "backend/internal/child/api"
	childmysql "backend/internal/child/mysql"

	"backend/internal/secretInfo"
	secretinfoapi "backend/internal/secretInfo/api"
	secretinfomysql "backend/internal/secretInfo/mysql"

	watchingmysql "backend/internal/watching/mysql"

	"backend/internal/event"
	eventapi "backend/internal/event/api"
	eventmysql "backend/internal/event/mysql"

	"backend/internal/profile"
	profileapi "backend/internal/profile/api"
	profilemysql "backend/internal/profile/mysql"

	"backend/internal/lifecycle"
	lifecycleapi "backend/internal/lifecycle/api"
	lifecyclemysql "backend/internal/lifecycle/mysql"

	"backend/internal/request"
	requestapi "backend/internal/request/api"
	requestmysql "backend/internal/request/mysql"

	"backend/internal/document"
	documentapi "backend/internal/document/api"
	documentmysql "backend/internal/document/mysql"
	documentstorage "backend/internal/document/storage"

	"backend/internal/donor"
	donorapi "backend/internal/donor/api"
	donormysql "backend/internal/donor/mysql"

	"backend/internal/report"
	reportapi "backend/internal/report/api"
	reportmysql "backend/internal/report/mysql"

	"backend/internal/learn"
	learnclient "backend/internal/learn/infrastructure/client"
	learnmysql "backend/internal/learn/infrastructure/mysql"

	"backend/internal/user"
	userapi "backend/internal/user/api"
	usermysql "backend/internal/user/mysql"
)

type helpAdapter struct {
	s *help.Service
}

func (a helpAdapter) MotherOf(helpId int) (int, bool, error) {
	item, err := a.s.ById(context.Background(), helpId)
	if errors.Is(err, help.ErrorNotFound) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return item.MotherId, true, nil
}

type motherReaderAdapter struct {
	repo mother.Repository
}

func (a motherReaderAdapter) GetAll() ([]learn.MotherRef, error) {
	mothers, err := a.repo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]learn.MotherRef, len(mothers))
	for i, m := range mothers {
		result[i] = learn.MotherRef{
			Id:        m.Id,
			FirstName: m.FirstName,
			LastName:  m.LastName,
			Phone:     m.Phone,
			Email:     m.Email,
		}
	}
	return result, nil
}

func (a motherReaderAdapter) UpdateLearnUserID(id int, learnUserID string) error {
	return a.repo.UpdateLearnUserID(id, learnUserID)
}

type learnerReaderAdapter struct {
	repo *learnmysql.Repo
}

func (a learnerReaderAdapter) GetAllLearners() ([]learn.Learner, error) {
	learners, err := a.repo.GetAllLearners()
	if err != nil {
		return nil, err
	}
	result := make([]learn.Learner, len(learners))
	for i, l := range learners {
		result[i] = learn.Learner{
			UserId:    l.UserId,
			Login:     l.Login,
			Email:     l.Email,
			FirstName: l.FirstName,
			LastName:  l.LastName,
		}
	}
	return result, nil
}

func main() {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("база недоступна: ", err)
	}

	filesDir := os.Getenv("FILES_DIR")
	if filesDir == "" {
		filesDir = "./uploads"
	}
	if err := os.MkdirAll(filesDir, 0o750); err != nil {
		log.Fatal("папка для файлов недоступна: ", err)
	}

	motherRepo := mothermysql.NewRepo(db)
	learnRepo := learnmysql.NewRepo(db)
	motherService := mother.NewService(motherRepo, learnRepo)
	motherHandler := motherapi.NewMotherHandler(motherService, filesDir)

	watchingRepo := watchingmysql.NewRepo(db)

	secretRepo := secretinfomysql.NewRepo(db)
	secretService := secretInfo.NewService(secretRepo, motherService, watchingRepo)
	secretHandler := secretinfoapi.New(secretService)

	helpRepo := helpmysql.NewRepo(db)
	helpService := help.NewService(helpRepo, motherService)
	helpHandler := helpapi.New(helpService)

	childRepo := childmysql.NewRepo(db)
	childService := child.NewService(childRepo, motherService)
	childHandler := childapi.New(childService)

	eventRepo := eventmysql.NewRepo(db)
	eventService := event.NewService(eventRepo, motherService)
	eventHandler := eventapi.New(eventService)

	profileRepo := profilemysql.NewRepo(db)
	profileService := profile.NewService(profileRepo, motherService)
	profileHandler := profileapi.New(profileService)

	lifecycleRepo := lifecyclemysql.NewRepo(db)
	lifecycleService := lifecycle.NewService(lifecycleRepo, motherService)
	lifecycleHandler := lifecycleapi.New(lifecycleService)

	requestRepo := requestmysql.NewRepo(db)
	requestService := request.NewService(requestRepo, motherService, helpAdapter{s: helpService})
	requestHandler := requestapi.New(requestService)

	files, err := documentstorage.NewLocal(filesDir)
	if err != nil {
		log.Fatal("папка для файлов недоступна: ", err)
	}

	documentRepo := documentmysql.NewRepo(db)
	documentService := document.NewService(documentRepo, files, motherService, helpAdapter{s: helpService})
	documentHandler := documentapi.New(documentService)

	donorRepo := donormysql.NewRepo(db)
	donorService := donor.NewService(donorRepo)
	donorHandler := donorapi.New(donorService)

	reportRepo := reportmysql.NewRepo(db)
	reportService := report.NewService(reportRepo)
	reportHandler := reportapi.New(reportService)

	userRepo := usermysql.NewRepo(db)
	userService := user.NewService(userRepo, os.Getenv("AUTH_SECRET"), 24*time.Hour)
	userHandler := userapi.New(userService)

	mux := http.NewServeMux()
	helpHandler.Register(mux)
	motherHandler.Register(mux)
	secretHandler.Register(mux)
	childHandler.Register(mux)
	eventHandler.Register(mux)
	profileHandler.Register(mux)
	lifecycleHandler.Register(mux)
	requestHandler.Register(mux)
	documentHandler.Register(mux)
	donorHandler.Register(mux)
	reportHandler.Register(mux)
	userHandler.Register(mux)

	mux.Handle("GET /api/photos/{path...}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fileName := r.PathValue("path")
		if fileName == "" {
			http.NotFound(w, nil)
			return
		}
		filePath := filepath.Join(filesDir, fileName)
		http.ServeFile(w, r, filePath)
	}))

	learnAPIBaseURL := os.Getenv("LEARN_API_BASE_URL")
	learnAPIClientID := os.Getenv("LEARN_API_CLIENT_ID")
	learnAPIClientSecret := os.Getenv("LEARN_API_CLIENT_SECRET")

	var learnAPI learn.LearnApiClient
	if learnAPIBaseURL != "" && learnAPIClientID != "" && learnAPIClientSecret != "" {
		learnAPI = learnclient.NewClient(learnAPIBaseURL, learnAPIClientID, learnAPIClientSecret)
		log.Println("Learn API клиент настроен (OAuth 2.0)")
	} else {
		log.Println("Learn API не настроен: укажите LEARN_API_BASE_URL, LEARN_API_CLIENT_ID и LEARN_API_CLIENT_SECRET в .env")
	}

	learnSyncService := learn.NewSyncService(learnRepo, learnAPI)
	learnSyncService.SetLinker(learnerReaderAdapter{repo: learnRepo}, motherReaderAdapter{repo: motherRepo})

	go func() {
		for {
			if learnAPI != nil {
				if err := learnSyncService.SyncAll(); err != nil {
					log.Printf("SyncAll ошибка: %v", err)
				}
			}
			if err := learnSyncService.SyncUsers(); err != nil {
				log.Printf("SyncUsers ошибка: %v", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()

	log.Println("сервер запущен на", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}