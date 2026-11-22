package app

import (
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
	eventadapter "backend/internal/event/infrastructure/adapter"
	eventcommand "backend/internal/event/application/command"
	eventmysql "backend/internal/event/infrastructure/mysql"
	eventhttp "backend/internal/event/interfaces/httpapi"
	"backend/internal/help"
	helpapi "backend/internal/help/api"
	helpmysql "backend/internal/help/mysql"
	"backend/internal/learn"
	learnclient "backend/internal/learn/infrastructure/client"
	learnmysql "backend/internal/learn/infrastructure/mysql"
	lifecycleadapter "backend/internal/lifecycle/infrastructure/adapter"
	lifecyclecommand "backend/internal/lifecycle/application/command"
	lifecyclemysql "backend/internal/lifecycle/infrastructure/mysql"
	lifecyclehttp "backend/internal/lifecycle/interfaces/httpapi"
	"backend/internal/mother"
	motherapi "backend/internal/mother/api"
	mothermysql "backend/internal/mother/mysql"
	profileadapter "backend/internal/profile/infrastructure/adapter"
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

func NewDependencyContainer(cfg Config, conns *ConnectionContainer) *DependencyContainer {
	db := conns.DB()
	now := time.Now

	motherRepo := mothermysql.NewRepo(db)
	learnRepo := learnmysql.NewRepo(db)
	motherService := mother.NewService(motherRepo, learnRepo)
	motherHandler := motherapi.NewMotherHandler(motherService, cfg.FilesDir)

	watchingRepo := watchingmysql.NewRepo(db)

	secretService := secretInfo.NewService(secretinfomysql.NewRepo(db), motherService, watchingRepo)
	secretHandler := secretinfoapi.New(secretService)

	helpService := help.NewService(helpmysql.NewRepo(db), motherService)
	helpHandler := helpapi.New(helpService)

	childHandler := childapi.New(child.NewService(childmysql.NewRepo(db), motherService))

	eventRepo := eventmysql.NewRepository(db)
	eventUoW := eventmysql.NewUnitOfWork(db)
	eventReads := eventmysql.NewReadModel(db)
	eventMothers := eventadapter.NewMotherDirectory(motherService)
	eventHandler := eventhttp.NewHandler(
		eventcommand.NewCreateEventHandler(eventRepo),
		eventcommand.NewRegisterHandler(eventUoW, eventMothers),
		eventcommand.NewMarkHandler(eventUoW),
		eventReads,
	)

	profileRepo := profilemysql.NewRepository(db)
	profileReads := profilemysql.NewReadModel(db)
	profileMothers := profileadapter.NewMotherDirectory(motherService)
	profileHandler := profilehttp.NewHandler(
		profilecommand.NewSaveHandler(profileRepo, profileMothers, now),
		profilecommand.NewSetStatusHandler(profileRepo),
		profileReads,
	)
	lifecycleRepo := lifecyclemysql.NewRepository(db)
	lifecycleReads := lifecyclemysql.NewReadModel(db)
	lifecycleMothers := lifecycleadapter.NewMotherDirectory(motherService)
	lifecycleHandler := lifecyclehttp.NewHandler(
		lifecyclecommand.NewChangeHandler(lifecycleRepo, lifecycleMothers),
		lifecycleReads,
	)

	requestRepo := requestmysql.NewRepository(db)
	requestUoW := requestmysql.NewUnitOfWork(db)
	requestReads := requestmysql.NewReadModel(db)
	requestMothers := requestadapter.NewMotherDirectory(motherService)
	requestHelps := requestadapter.NewHelpDirectory(helpService)
	requestHandler := requesthttp.NewHandler(
		requestcommand.NewCreateRequestHandler(requestRepo, requestMothers, now),
		requestcommand.NewChangeStatusHandler(requestUoW, now),
		requestcommand.NewFulfillRequestHandler(requestUoW, requestHelps, now),
		requestReads,
	)

	documentStorage, err := documentstorage.NewLocal(cfg.FilesDir)
	if err != nil {
		log.Fatalf("file storage: %v", err)
	}
	documentRepo := documentmysql.NewRepository(db)
	documentReads := documentmysql.NewReadModel(db)
	documentMothers := documentadapter.NewMotherDirectory(motherService)
	documentHelps := documentadapter.NewHelpDirectory(helpService)
	documentHandler := documenthttp.NewHandler(
		documentcommand.NewCreateHandler(documentRepo, documentMothers, documentHelps),
		documentcommand.NewAttachFileHandler(documentRepo, documentStorage),
		documentcommand.NewDownloadHandler(documentRepo, documentStorage),
		documentcommand.NewSignHandler(documentRepo),
		documentReads,
	)

	donorHandler := donorapi.New(donor.NewService(donormysql.NewRepo(db)))
	reportHandler := reportapi.New(report.NewService(reportmysql.NewRepo(db)))
	userHandler := userapi.New(user.NewService(usermysql.NewRepo(db), cfg.AuthSecret, cfg.AuthTTL))

	var learnAPI learn.LearnApiClient
	if cfg.Learn.Enabled() {
		learnAPI = learnclient.NewClient(cfg.Learn.BaseURL, cfg.Learn.ClientID, cfg.Learn.ClientSecret)
		log.Println("Learn API клиент настроен (OAuth 2.0)")
	} else {
		log.Println("Learn API не настроен: укажите LEARN_API_BASE_URL, LEARN_API_CLIENT_ID и LEARN_API_CLIENT_SECRET в .env")
	}
	syncService := learn.NewSyncService(learnRepo, learnAPI)
	syncService.SetLinker(learnerReaderAdapter{repo: learnRepo}, motherReaderAdapter{repo: motherRepo})

	return &DependencyContainer{
		routers: []httpx.Router{
			helpHandler, motherHandler, secretHandler, childHandler, eventHandler,
			profileHandler, lifecycleHandler, requestHandler, documentHandler,
			donorHandler, reportHandler, userHandler,
		},
		filesDir:    cfg.FilesDir,
		learnWorker: NewLearnSyncWorker(syncService, learnAPI != nil, 5*time.Minute),
	}
}

func (d *DependencyContainer) Routers() []httpx.Router        { return d.routers }
func (d *DependencyContainer) FilesDir() string                { return d.filesDir }
func (d *DependencyContainer) LearnSyncWorker() *LearnSyncWorker { return d.learnWorker }