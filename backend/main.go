package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"context"
	"errors"

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

	motherRepo := mothermysql.NewRepo(db)
	motherService := mother.NewService(motherRepo)
	motherHandler := motherapi.NewMotherHandler(motherService)

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

	filesDir := os.Getenv("FILES_DIR")
	if filesDir == "" {
		filesDir = "./uploads"
	}
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

	log.Println("сервер запущен на", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}