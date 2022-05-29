package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/addilafzal/rss-relay/internal/parser"
	"github.com/addilafzal/rss-relay/internal/rss"
	"github.com/addilafzal/rss-relay/internal/transmission"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	transmissionHost     string
	transmissionPort     int
	transmissionUsername string
	transmissionPassword string
	transmissionHTTPs    bool
	configFile           string
	cronSchedule         string
	sigs                 chan os.Signal
)

func init() {
	flag.StringVar(&transmissionHost, "host", "localhost", "transmission host address")
	flag.IntVar(&transmissionPort, "port", 9091, "transmission host port")
	flag.StringVar(&transmissionUsername, "username", "", "transmission username")
	flag.StringVar(&transmissionPassword, "password", "", "transmission password")
	flag.BoolVar(&transmissionHTTPs, "https", true, "whether to communicate with the rpc endpoint over https")
	flag.StringVar(&configFile, "config", "/etc/rss-relay/config.yaml", "path to the config file")
	flag.StringVar(&cronSchedule, "cron", "*/1 * * * *", "cron schedule")

	flag.Parse()
	sigs = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	logrus.Info("Starting")
	c := cron.New()
	c.AddFunc(cronSchedule, run)
	c.Start()
	<-sigs
	logrus.Info("Recieved SIGTERM")
}

func run() {
	transmissionClient := transmission.NewTransmissionClient(&transmission.TransmisionConfig{
		Host:     transmissionHost,
		Port:     transmissionPort,
		Username: transmissionUsername,
		Password: transmissionPassword,
		HTTPs:    transmissionHTTPs,
	})
	defer transmissionClient.SessionClose(context.Background())
	file, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Failed to open config file: ", err)
	}
	rssConfig, err := parser.ParseConfigFile(file)
	if err != nil {
		log.Printf("Failed to parse config file: %s", err)
		return
	}

	matchingItems := make([]rss.DownloadItem, 0)

	logrus.Infof("Found %d items", len(matchingItems))

	for _, source := range rssConfig.Source {
		matchingItems = append(matchingItems, source.FindMatchingItems()...)
	}

	for _, downloadItem := range matchingItems {
		transmission.AddMagnetLinkToDownloads(transmissionClient, downloadItem)
	}
}
