package es

import (
	"github.com/loveyu233/gu/client"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type EsConfig struct {
	Endpoints string
	Username  string
	Password  string
}

func MustInitESClient(config ...EsConfig) *elastic.Client {
	var esConnConfig []elastic.ClientOptionFunc
	if len(config) > 0 {
		if config[0].Endpoints == "" {
			esConnConfig = append(esConnConfig, elastic.SetURL("http://127.0.0.1:9200"))
		} else {
			esConnConfig = append(esConnConfig, elastic.SetURL(config[0].Endpoints))
		}
		if config[0].Username != "" && config[0].Password != "" {
			esConnConfig = append(esConnConfig, elastic.SetBasicAuth(config[0].Username, config[0].Password))
		}
	}

	esConnConfig = append(esConnConfig, elastic.SetSniff(false))
	esClient, err := elastic.NewClient(esConnConfig...)
	if err != nil {
		logrus.Panicf("create es client failed: %v\n", err)
	}

	logrus.Info("successfully init es client")

	client.EsClient = esClient
	return esClient
}
