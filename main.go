package main

import (
	"fmt"
	log "github.com/cihub/seelog"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"flag"
	"github.com/tokenbankteam/sync_node_peers/config"
	"runtime"
	"github.com/tokenbankteam/tb_common/perf"
	"github.com/tokenbankteam/tb_common/health"
	"github.com/tokenbankteam/sync_node_peers/task"
	"github.com/tokenbankteam/sync_node_peers/service"
	"github.com/tokenbankteam/sync_node_peers/model/blockchain"
	"github.com/gin-gonic/gin"
	prom "github.com/tokenbankteam/tb_common/middleware"
	mlog "github.com/tokenbankteam/tb_common/middleware"
	"net/http"
	"context"
	"github.com/tokenbankteam/sync_node_peers/controller"
)

var (
	version = ""
)

func main() {
	showVersion := flag.Bool("version", false, "Show version")
	configFlag := flag.String("config", "config/config.yml", "configuration file")
	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	var err error
	var appConfig config.AppConfig
	if err = configor.Load(&appConfig, *configFlag); err != nil {
		log.Criticalf("load config error: %v", err)
		os.Exit(1)
	}

	var logger log.LoggerInterface
	logger, err = log.LoggerFromConfigAsFile(appConfig.Logger)
	if err != nil {
		log.Errorf("init logger from %s error: %v", appConfig.Logger, err)
	} else {
		log.ReplaceLogger(logger)
	}
	defer log.Flush()
	log.Infof("Started Application at %v", time.Now().Format("January 2, 2006 at 3:04pm (MST)"))
	log.Infof("Version: %v", version)

	runtime.GOMAXPROCS(runtime.NumCPU())

	// start prof
	perf.Init(appConfig.PprofAddrs)
	if appConfig.Monitor {
		health.InitMonitor(appConfig.MonitorAddrs)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// Global middleware
	router.Use(mlog.Logger())
	router.Use(gin.Recovery())

	p := prom.NewPrometheus("gin")
	p.Use(router)

	blockChainModel, err := blockchain.NewModel(&appConfig)
	if err != nil {
		log.Errorf("init blockChainModel error, %v", err)
		return
	}
	appContext := &service.AppContext{
		Config:   &appConfig,
		Services: map[string]interface{}{},
		Models: map[string]interface{}{
			"blockchainModel": blockChainModel,
		},
	}

	appContext.Services["blockchainService"], _ = service.NewBlockChainService(appContext)

	versionController, _ := controller.NewVersionController(appContext)
	router.GET("/v1/version", versionController.GetVersion)

	//业务接口　结束

	//启动定时任务管理器
	taskManager, err := task.NewTaskManager(&appConfig, appContext)
	if err != nil {
		log.Errorf("new task basePlugin error, %v", err)
		return
	}
	taskManager.Start()

	srv := &http.Server{
		Addr:           appConfig.Addr, // listen and serve on 0.0.0.0:8080
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil {
			log.Error("listen: ", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with  a timeout of 5 seconds.
	InitSignal()
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server Shutdown:", err)
	}
	log.Info("Server exist")
}
