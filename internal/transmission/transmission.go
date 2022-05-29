package transmission

import (
	"context"

	"github.com/addilafzal/rss-relay/internal/rss"
	transmissionrpc "github.com/hekmon/transmissionrpc/v2"
	"github.com/sirupsen/logrus"
)

type TransmisionConfig struct {
	Host     string
	Port     int
	HTTPs    bool
	Username string
	Password string
}

func NewTransmissionClient(config *TransmisionConfig) *transmissionrpc.Client {
	client, err := transmissionrpc.New(config.Host, config.Username, config.Password,
		&transmissionrpc.AdvancedConfig{
			HTTPS: config.HTTPs,
			Port:  uint16(config.Port),
		})
	if err != nil {
		logrus.Fatal("Failed to create transmission client: %s", err)
	}

	_, serverVersion, serverMinimumVersion, err := client.RPCVersion(context.Background())
	if err != nil {
		logrus.Errorf("Client failed to get RPC version from tranismission server: %s", err)
		return client
	}

	logrus.Debugf("Server version: %v", serverVersion)
	logrus.Debugf("Server minimum version: %v", serverMinimumVersion)

	return client
}

func AddMagnetLinkToDownloads(client *transmissionrpc.Client, downloadItem rss.DownloadItem) {
	torrent, err := client.TorrentAdd(context.Background(),
		transmissionrpc.TorrentAddPayload{
			Filename:    &downloadItem.Item.MagnetURI,
			DownloadDir: &downloadItem.Source.DownloadDirectory,
		})
	if err != nil {
		logrus.Warningf("Failed to add torrent to server: %s", err)
	} else {
		logrus.Infof("Added '%s'. Downloading to directory %s", *torrent.Name, &downloadItem.Source.DownloadDirectory)
	}
}
