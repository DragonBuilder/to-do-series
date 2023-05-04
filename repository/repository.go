package repository

import (
	"context"

	"github.com/DragonBuilder/to-do-series/domain"
	"github.com/sirupsen/logrus"
)

var dbClient client

func init() {
	dbClient = new(firestoreClient)
	if err := dbClient.Connect(); err != nil {
		logrus.WithError(err).Fatalln("error connecting to database")
	}
	logrus.Infoln("Firestore client init successful")
}

type client interface {
	Connect() error
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	Read(ctx context.Context, uid string) (*domain.Task, error)
}

func NewTaskRepository() TaskRepository {
	return &taskFirestoreRepository{
		client: dbClient.(*firestoreClient),
	}
}
