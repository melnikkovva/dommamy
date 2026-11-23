package app

import (
	"database/sql"
	"log"
	"time"

	"backend/internal/child"
	childapi "backend/internal/child/api"
	childmysql "backend/internal/child/mysql"
	documentadapter "backend/internal/document/infrastructure/adapter"
	documentcommand "backend/internal/document/application/command"
	documentmysql "backend/internal/document/infrastructure/mysql"
	documenthttp "backend/internal/document/interfaces/httpapi"
	documentstorage "backend/internal/document/infrastructure/storage"
	"backend/internal/donor"
	donorapi "backend/internal/donor/api"
	donormysql "backend/internal/donor/mysql"
	eventcommand "backend/internal/event/application/command"
	eventmysql "backend/internal/event/infrastructure/mysql"
	eventhttp "backend/internal/event/interfaces/httpapi"
	"backend/internal/help"
	helpapi "backend/internal/help/api"
	helpmysql "backend/internal/help/mysql"
	learndomain "backend/internal/learn/domain"
	learnadapter "backend/internal/learn/infrastructure/adapter"
	learnclient "backend/internal/learn/infrastructure/client"
	learncommand "backend/internal/learn/application/command"
	learnmysql "backend/internal/learn/infrastructure/mysql"
	lifecyclecommand "backend/internal/lifecycle/application/command"
	lifecyclemysql "backend/internal/lifecycle/infrastructure/mysql"
	lifecyclehttp "backend/internal/lifecycle/interfaces/httpapi"
	motheradapter "backend/internal/mother/infrastructure/adapter"
	mothercommand "backend/internal/mother/application/command"
	mothermysql "backend/internal/mother/infrastructure/mysql"
	motherhttp "backend/internal/mother/interfaces/httpapi"
	motherquery "backend/internal/mother/application/query"
	motherstorage "backend/internal/mother/infrastructure/storage"
	profilecommand "backend/internal/profile/application/command"
	profilemysql "backend/internal/profile/infrastructure/mysql"
	profilehttp "backend/internal/profile/interfaces/httpapi"
	"backend/internal/report"
	reportapi "backend/internal/report/api"
	reportmysql "backend/internal/report/mysql"
	"backend/internal/secretInfo"
	secretinfoapi "backend/internal/secretInfo/api"
	secretinfomysql "backend/internal/secretInfo/mysql"
	"backend/internal/shared/httpx"
	"backend/internal/user"
	userapi "backend/internal/user/api"
	usermysql "backend/internal/user/mysql"
	watchingmysql "backend/internal/watching/mysql"

	requestadapter "backend/internal/request/infrastructure/adapter"
	requestmysql "backend/internal/request/infrastructure/mysql"
	requestcommand "backend/internal/request/application/command"
	requesthttp "backend/internal/request/interfaces/httpapi"
)

type DependencyContainer struct {
	routers     []httpx.Router
	filesDir    string
	learnWorker *LearnSyncWorker
}

func NewDependencyContainer(cfg Config, conns *ConnectionContainer) (*DependencyContainer, error) {
	db := conns.DB()
	now := time.Now

	motherDirectory, motherHandler := wireMother(db, cfg.FilesDir)

	secretHandler := wireSecretInfo(db, motherDirectory)

	helpHandler, helpService := wireHelp(db, motherDirectory)

	childHandler := wireChild(db, motherDirectory)

	eventHandler := wireEvent(db, motherDirectory)

	profileHandler := wireProfile(db, motherDirectory, now)

	lifecycleHandler := wireLifecycle(db, motherDirectory)

	requestHandler, err := wireRequest(db, motherDirectory, helpService, now)
	if err != nil {
		return nil, err
	}

	documentHandler, err := wireDocument(db, cfg.FilesDir, motherDirectory, helpService)
	if err != nil {
		return nil, err
	}

	donorHandler := wireDonor(db)
	reportHandler := wireReport(db)
	userHandler := wireUser(db, cfg)

	syncHandler, linkHandler := wireLearn(db, motherDirectory, cfg)

	return &DependencyContainer{
		routers: []httpx.Router{
			helpHandler, motherHandler, secretHandler, childHandler, eventHandler,
			profileHandler, lifecycleHandler, requestHandler, documentHandler,
			donorHandler, reportHandler, userHandler,
		},
		filesDir:    cfg.FilesDir,
		learnWorker: NewLearnSyncWorker(syncHandler, linkHandler, isLearnEnabled(cfg), 5*time.Minute),
	}, nil
}

func wireMother(db *sql.DB, filesDir string) (motherDirectory *motheradapter.MotherDirectory, handler httpx.Router) {
	motherRepo := mothermysql.NewRepository(db)
	learnReadModel := learnmysql.NewReadModel(db)
	learnStatsAdapter := motheradapter.NewLearnStatsAdapter(learnReadModel)
	motherReads := mothermysql.NewReadModel(db, learnStatsAdapter)
	photoStorage := motherstorage.NewLocalStorage(filesDir)

	handler = motherhttp.NewHandler(
		mothercommand.NewCreateHandler(motherRepo),
		mothercommand.NewUploadPhotoHandler(motherRepo, photoStorage, motherReads),
		motherquery.NewGetCoursesHandler(motherReads, learnStatsAdapter),
		motherquery.NewGetIPRHandler(motherReads, learnStatsAdapter),
		motherReads,
	)
	motherDirectory = motheradapter.NewMotherDirectory(motherReads)
	return
}

func wireSecretInfo(db *sql.DB, mothers *motheradapter.MotherDirectory) httpx.Router {
	watchingRepo := watchingmysql.NewRepo(db)
	secretService := secretInfo.NewService(secretinfomysql.NewRepo(db), mothers, watchingRepo)
	return secretinfoapi.New(secretService)
}

func wireHelp(db *sql.DB, mothers *motheradapter.MotherDirectory) (handler httpx.Router, service *help.Service) {
	service = help.NewService(helpmysql.NewRepo(db), mothers)
	return helpapi.New(service), service
}

func wireChild(db *sql.DB, mothers *motheradapter.MotherDirectory) httpx.Router {
	return childapi.New(child.NewService(childmysql.NewRepo(db), mothers))
}

func wireEvent(db *sql.DB, mothers *motheradapter.MotherDirectory) httpx.Router {
	eventRepo := eventmysql.NewRepository(db)
	eventUoW := eventmysql.NewUnitOfWork(db)
	eventReads := eventmysql.NewReadModel(db)
	return eventhttp.NewHandler(
		eventcommand.NewCreateEventHandler(eventRepo),
		eventcommand.NewRegisterHandler(eventUoW, mothers),
		eventcommand.NewMarkHandler(eventUoW),
		eventReads,
	)
}

func wireProfile(db *sql.DB, mothers *motheradapter.MotherDirectory, now func() time.Time) httpx.Router {
	profileRepo := profilemysql.NewRepository(db)
	profileReads := profilemysql.NewReadModel(db)
	return profilehttp.NewHandler(
		profilecommand.NewSaveHandler(profileRepo, mothers, now),
		profilecommand.NewSetStatusHandler(profileRepo),
		profileReads,
	)
}

func wireLifecycle(db *sql.DB, mothers *motheradapter.MotherDirectory) httpx.Router {
	lifecycleRepo := lifecyclemysql.NewRepository(db)
	lifecycleReads := lifecyclemysql.NewReadModel(db)
	return lifecyclehttp.NewHandler(
		lifecyclecommand.NewChangeHandler(lifecycleRepo, mothers),
		lifecycleReads,
	)
}

func wireRequest(db *sql.DB, mothers *motheradapter.MotherDirectory, helps *help.Service, now func() time.Time) (httpx.Router, error) {
	requestRepo := requestmysql.NewRepository(db)
	requestUoW := requestmysql.NewUnitOfWork(db)
	requestReads := requestmysql.NewReadModel(db)
	requestHelps := requestadapter.NewHelpDirectory(helps)
	return requesthttp.NewHandler(
		requestcommand.NewCreateRequestHandler(requestRepo, mothers, now),
		requestcommand.NewChangeStatusHandler(requestUoW, now),
		requestcommand.NewFulfillRequestHandler(requestUoW, requestHelps, now),
		requestReads,
	), nil
}

func wireDocument(db *sql.DB, filesDir string, mothers *motheradapter.MotherDirectory, helps *help.Service) (httpx.Router, error) {
	documentStorage, err := documentstorage.NewLocal(filesDir)
	if err != nil {
		return nil, err
	}
	documentRepo := documentmysql.NewRepository(db)
	documentReads := documentmysql.NewReadModel(db)
	documentHelps := documentadapter.NewHelpDirectory(helps)
	return documenthttp.NewHandler(
		documentcommand.NewCreateHandler(documentRepo, mothers, documentHelps),
		documentcommand.NewAttachFileHandler(documentRepo, documentStorage),
		documentcommand.NewDownloadHandler(documentRepo, documentStorage),
		documentcommand.NewSignHandler(documentRepo),
		documentReads,
	), nil
}

func wireDonor(db *sql.DB) httpx.Router {
	return donorapi.New(donor.NewService(donormysql.NewRepo(db)))
}

func wireReport(db *sql.DB) httpx.Router {
	return reportapi.New(report.NewService(reportmysql.NewRepo(db)))
}

func wireUser(db *sql.DB, cfg Config) httpx.Router {
	return userapi.New(user.NewService(usermysql.NewRepo(db), cfg.AuthSecret, cfg.AuthTTL))
}

func wireLearn(db *sql.DB, mothers *motheradapter.MotherDirectory, cfg Config) (syncHandler learncommand.SyncAll, linkHandler learncommand.LinkUsers) {
	learnWriteRepo := learnmysql.NewWriteRepo(db)
	learnReadModel := learnmysql.NewReadModel(db)

	var learnAPI learndomain.APIClient
	if cfg.Learn.Enabled() {
		learnAPI = learnclient.NewClient(cfg.Learn.BaseURL, cfg.Learn.ClientID, cfg.Learn.ClientSecret)
		log.Println("Learn API клиент настроен (OAuth 2.0)")
	} else {
		log.Println("Learn API не настроен: укажите LEARN_API_BASE_URL, LEARN_API_CLIENT_ID и LEARN_API_CLIENT_SECRET в .env")
	}

	syncHandler = learncommand.NewSyncHandler(learnWriteRepo, learnAPI)
	motherDir := learnadapter.NewMotherDirectoryAdapter(mothermysql.NewReadModel(db, motheradapter.NewLearnStatsAdapter(learnReadModel)))
	linkRepo := learnmysql.NewMotherLearnLinkRepo(db)
	linkHandler = learncommand.NewLinkUsersHandler(learnReadModel, motherDir, linkRepo)
	return
}

func isLearnEnabled(cfg Config) bool {
	return cfg.Learn.Enabled()
}

func (d *DependencyContainer) Routers() []httpx.Router        { return d.routers }
func (d *DependencyContainer) FilesDir() string                { return d.filesDir }
func (d *DependencyContainer) LearnSyncWorker() *LearnSyncWorker { return d.learnWorker }
