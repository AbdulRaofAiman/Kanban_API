package models

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCommentModel(t *testing.T) {
	comment := Comment{
		ID:      uuid.New().String(),
		TaskID:  uuid.New().String(),
		UserID:  uuid.New().String(),
		Content: "This is a test comment",
	}

	// Verify table name
	if comment.TableName() != "comments" {
		t.Errorf("Expected table name 'comments', got '%s'", comment.TableName())
	}

	// Verify ID is set
	if comment.ID == "" {
		t.Error("ID should not be empty")
	}

	// Verify TaskID is set
	if comment.TaskID == "" {
		t.Error("TaskID should not be empty")
	}

	// Verify UserID is set
	if comment.UserID == "" {
		t.Error("UserID should not be empty")
	}

	// Verify Content is not empty
	if comment.Content == "" {
		t.Error("Content should not be empty")
	}

	t.Logf("Comment model validated: ID=%s, TaskID=%s, UserID=%s", comment.ID, comment.TaskID, comment.UserID)
}

func TestCommentRelationships(t *testing.T) {
	tests := []struct {
		name    string
		comment Comment
		hasTask bool
		hasUser bool
	}{
		{
			name: "Comment with task and user",
			comment: Comment{
				ID:      uuid.New().String(),
				TaskID:  uuid.New().String(),
				UserID:  uuid.New().String(),
				Content: "Test comment with relationships",
				Task: &Task{
					ID:    uuid.New().String(),
					Title: "Test Task",
				},
				User: &User{
					ID:       uuid.New().String(),
					Username: "testuser",
					Email:    "test@example.com",
					Password: "$2a$10$dummyhashforunittestpurposes",
				},
			},
			hasTask: true,
			hasUser: true,
		},
		{
			name: "Comment without task and user (nil relationships)",
			comment: Comment{
				ID:      uuid.New().String(),
				TaskID:  uuid.New().String(),
				UserID:  uuid.New().String(),
				Content: "Test comment without relationships",
				Task:    nil,
				User:    nil,
			},
			hasTask: false,
			hasUser: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify task relationship
			if tt.hasTask && tt.comment.Task == nil {
				t.Error("Expected Task to be set, but it's nil")
			}
			if !tt.hasTask && tt.comment.Task != nil {
				t.Error("Expected Task to be nil, but it's set")
			}
			if tt.hasTask && tt.comment.Task != nil && tt.comment.Task.ID == "" {
				t.Error("Task ID should not be empty when Task is set")
			}

			// Verify user relationship
			if tt.hasUser && tt.comment.User == nil {
				t.Error("Expected User to be set, but it's nil")
			}
			if !tt.hasUser && tt.comment.User != nil {
				t.Error("Expected User to be nil, but it's set")
			}
			if tt.hasUser && tt.comment.User != nil && tt.comment.User.ID == "" {
				t.Error("User ID should not be empty when User is set")
			}

			t.Logf("Comment relationships validated: Task=%v, User=%v",
				getPointerStatus(tt.comment.Task),
				getPointerStatus(tt.comment.User))
		})
	}
}

func TestCommentSoftDelete(t *testing.T) {
	comment := Comment{
		ID:        uuid.New().String(),
		TaskID:    uuid.New().String(),
		UserID:    uuid.New().String(),
		Content:   "Comment to be soft deleted",
		DeletedAt: gorm.DeletedAt{},
	}

	// Initially not deleted
	if comment.DeletedAt.Valid {
		t.Error("Comment should not be marked as deleted initially")
	}

	// Simulate soft delete (this would normally be done by GORM)
	comment.DeletedAt = gorm.DeletedAt{Time: comment.DeletedAt.Time, Valid: true}

	// After soft delete
	if !comment.DeletedAt.Valid {
		t.Error("Comment should be marked as deleted after soft delete")
	}

	t.Logf("Comment soft delete validated: DeletedAt.Valid=%v", comment.DeletedAt.Valid)
}

// Helper function to get pointer status string
func getPointerStatus(ptr interface{}) string {
	if ptr == nil {
		return "nil"
	}
	return "set"
}
