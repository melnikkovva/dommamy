package app

import (
	"log"
	"time"

	"backend/internal/child"
	childapi "backend/internal/child/api"
	childmysql "backend/internal/child/mysql"
	"backend/internal/document"
	documentapi "backend/internal/document/api"
	documentmysql "backend/internal/document/mysql"
	"backend/internal/donor"
	donorapi "backend/internal/donor/api"
	donormysql "backend/internal/donor/mysql"
	"backend/internal/event"
	eventapi "backend/internal/event/api"
	eventmysql "backend/internal/event/mysql"
	"backend/internal/help"
	helpapi "backend/internal/help/api"
	helpmysql "backend/internal/help/mysql"
	"backend/internal/learn"
	learnclient "backend/internal/learn/infrastructure/client"
	learnmysql "backend/internal/learn/infrastructure/mysql"
	"backend/internal/lifecycle"
	lifecycleapi "backend/internal/lifecycle/api"
	lifecyclemysql "backend/internal/lifecycle/mysql"
	"backend/internal/mother"
	motherapi "backend/internal/mother/api"
	mothermysql "backend/internal/mother/mysql"
	"backend/internal/profile"
	profileapi "backend/internal/profile/api"
	profilemysql "backend/internal/profile/mysql"
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
	eventHandler := eventapi.New(event.NewService(eventmysql.NewRepo(db), motherService))
	profileHandler := profileapi.New(profile.NewService(profilemysql.NewRepo(db), motherService))
	lifecycleHandler := lifecycleapi.New(lifecycle.NewService(lifecyclemysql.NewRepo(db), motherService))

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

	documentService := document.NewService(documentmysql.NewRepo(db), conns.Files(), motherService, helpAdapter{s: helpService})
	documentHandler := documentapi.New(documentService)

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