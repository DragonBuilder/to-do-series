package repository

import (
	"context"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	if err := deleteCollection(context.Background(), dbClient.(*firestoreClient), tasksTable, 20); err != nil {
		logrus.WithError(err).Errorf("error deleting the mock data from collection %s\n", tasksTable)
	}
	if dbClient != nil {
		if err := dbClient.Close(); err != nil {
			logrus.WithError(err).Errorln("error closing db connection")
		}
	}
	os.Exit((exitVal))
}
