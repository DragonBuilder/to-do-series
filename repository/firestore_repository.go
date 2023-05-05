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
	tasksTable string = "tasks"
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

func (f *firestoreClient) Close() error {
	if err := f.Client.Close(); err != nil {
		return fmt.Errorf("error closing firestore client : %v", err)
	}
	return nil
}

type taskFirestoreRepository struct {
	client *firestoreClient
}

func (r *taskFirestoreRepository) Create(ctx context.Context, task *domain.Task) error {
	task.UID = uuid.New().String()
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if _, err := r.client.Collection(tasksTable).Doc(task.UID).Set(ctx, task); err != nil {
		return fmt.Errorf("error saving task to table %s : %v", tasksTable, err)
	}
	return nil
}

func (r *taskFirestoreRepository) Read(ctx context.Context, uid string) (*domain.Task, error) {
	iter := r.client.Collection(tasksTable).Where("UID", "==", uid).Documents(ctx)
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

func (r *taskFirestoreRepository) Update(ctx context.Context, new domain.Task) (*domain.Task, error) {
	if new.UID == "" {
		return nil, fmt.Errorf("cannot update Task without a UID")
	}
	existing, err := r.Read(ctx, new.UID)
	if err != nil {
		return nil, fmt.Errorf("error finding the task with UID %s : %v", new.UID, err)
	}
	existing.Content = new.Content
	existing.Priority = new.Priority
	existing.Status = new.Status
	existing.UpdatedAt = time.Now()
	if _, err := r.client.Collection(tasksTable).Doc(existing.UID).Set(ctx, existing); err != nil {
		return nil, fmt.Errorf("error updating task with UID %s to table %s : %v", existing.UID, tasksTable, err)
	}
	return existing, nil
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
