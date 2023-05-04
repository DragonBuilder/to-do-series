package repository

import (
	"context"
	"testing"

	"github.com/DragonBuilder/to-do-series/domain"
)

func Test_taskFirestoreRepository_Save(t *testing.T) {
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
			name:     "should save the to database",
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
			if err := tt.taskRepo.Save(tt.args.ctx, tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("taskFirestoreRepository.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
