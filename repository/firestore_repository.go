package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/DragonBuilder/to-do-series/domain"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

// var firestoreClient *firestore.Client

type firestoreClient struct {
	*firestore.Client
}

func (f *firestoreClient) Connect() error {
	ctx := context.Background()
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return fmt.Errorf("error creating new app for firestore init : %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error getting new firestore instance : %v", err)
	}
	f.Client = client
	return nil
}

type taskFirestoreRepository struct {
	client *firestoreClient
}

func (r *taskFirestoreRepository) Save(ctx context.Context, task *domain.Task) error {
	task.UID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if _, err := r.client.Collection("tasks").Doc(task.UID).Set(ctx, task); err != nil {
		return fmt.Errorf("error saving task : %v", err)
	}
	return nil
}

// func initFirestore() error {
// 	// Use a service account
// 	ctx := context.Background()
// 	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
// 	app, err := firebase.NewApp(ctx, nil, sa)
// 	if err != nil {
// 		return fmt.Errorf("error creating new app for firestore init : %v", err)
// 	}

// 	client, err := app.Firestore(ctx)
// 	if err != nil {
// 		return fmt.Errorf("error getting new firestore instance : %v", err)
// 	}
// 	firestoreClient = client
// 	logrus.Infoln("Firestore client init successful")
// 	return nil
// }
