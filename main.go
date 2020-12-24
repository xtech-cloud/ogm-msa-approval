package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"omo-msa-approval/config"
	"omo-msa-approval/handler"
	"omo-msa-approval/model"
	"omo-msa-approval/publisher"
	"os"
	"path/filepath"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	proto "github.com/xtech-cloud/omo-msp-approval/proto/approval"
)

func main() {
	config.Setup()
	model.Setup()
	defer model.Cancel()
	model.AutoMigrateDatabase()

	// New Service
	service := micro.NewService(
		micro.Name(config.Schema.Service.Name),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)

	// Initialise service
	service.Init()

	// Register publisher
	publisher.DefaultPublisher = micro.NewPublisher(config.Schema.Service.Name+".notification", service.Client())

	// Register Handler
	proto.RegisterWorkflowHandler(service.Server(), new(handler.Workflow))
	proto.RegisterOperatorHandler(service.Server(), new(handler.Operator))
	proto.RegisterTaskHandler(service.Server(), new(handler.Task))
	proto.RegisterActionHandler(service.Server(), new(handler.Action))

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Error(err)
	}
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
