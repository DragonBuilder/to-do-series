package repository

import (
	"context"
	"testing"

	"github.com/DragonBuilder/to-do-series/domain"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func deleteCollection(ctx context.Context, c *firestoreClient, name string, batchSize int) error {
	col := c.Collection(name)
	bulkwriter := c.BulkWriter(ctx)

	for {
		// Get a batch of documents
		iter := col.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to the BulkWriter.
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			bulkwriter.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			bulkwriter.End()
			break
		}

		bulkwriter.Flush()
	}
	// fmt.Fprintf(w, "Deleted collection \"%s\"", name)
	return nil
}

func Test_TaskRepository_Create(t *testing.T) {
	type args struct {
		ctx  context.Context
		task *domain.Task
	}
	tests := []struct {
		name     string
		taskRepo TaskRepository
		args     args
		wantErr  bool
	}{
		{
			name:     "should create the task in database",
			taskRepo: NewTaskRepository(),
			args: args{
				ctx: context.Background(),
				task: &domain.Task{
					Content:  "first task",
					Status:   domain.TaskIncomplete,
					Priority: 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.taskRepo.Create(tt.args.ctx, tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_TaskRepository_Read(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		taskRepo TaskRepository
		args     args
		want     *domain.Task
		wantErr  bool
	}{
		{
			name:     "Should fetch the Task with the given UID",
			taskRepo: NewTaskRepository(),
			args: args{
				ctx: context.Background(),
			},
			want: &domain.Task{
				Content:  "Read check task 1",
				Status:   domain.TaskCompleted,
				Priority: 2,
			},
		},
		{
			name:     "Should error when Task with given UID doesn't exist",
			taskRepo: NewTaskRepository(),
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.want != nil {
				if err := tt.taskRepo.Create(context.Background(), tt.want); err != nil {
					t.Errorf("error saving mock data : %v", err)
					return
				}
			}

			if tt.want == nil {
				tt.want = &domain.Task{Model: domain.Model{UID: uuid.New().String()}}
			}

			got, err := tt.taskRepo.Read(tt.args.ctx, tt.want.UID)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskRepository.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && (err != nil) {
				return
			}
			if got.UID != tt.want.UID {
				t.Errorf("TaskRepository.Read() got = %v, want %v", got.UID, tt.want.UID)
			}
		})
	}
}
