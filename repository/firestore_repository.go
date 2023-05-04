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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const (
	taskTable string = "task"
)

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

func (r *taskFirestoreRepository) Create(ctx context.Context, task *domain.Task) error {
	task.UID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if _, err := r.client.Collection(taskTable).Doc(task.UID).Set(ctx, task); err != nil {
		return fmt.Errorf("error saving task to table %s : %v", taskTable, err)
	}
	return nil
}

func (r *taskFirestoreRepository) Read(ctx context.Context, uid string) (*domain.Task, error) {
	iter := r.client.Collection(taskTable).Where("UID", "==", uid).Documents(ctx)
	var t domain.Task
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error fetching data from Tasks iterator : %v", err)
		}
		if err = doc.DataTo(&t); err != nil {
			return nil, fmt.Errorf("error unmarshalling data to Task struct : %v", err)
		}
		return &t, nil
	}
	return nil, fmt.Errorf("error fetching Task with UID %s", uid)
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
