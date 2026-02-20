package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestTaskModel(t *testing.T) {
	task := Task{
		ID:          uuid.NewString(),
		ColumnID:    uuid.NewString(),
		Title:       "Test Task",
		Description: "This is a test task description",
	}

	// Verify table name
	if task.TableName() != "tasks" {
		t.Errorf("Expected table name 'tasks', got '%s'", task.TableName())
	}

	// Verify ID is set
	if task.ID == "" {
		t.Error("ID should not be empty")
	}

	// Verify ColumnID is set
	if task.ColumnID == "" {
		t.Error("ColumnID should not be empty")
	}

	// Verify Title is set
	if task.Title == "" {
		t.Error("Title should not be empty")
	}

	// Verify Description can be set
	if task.Description == "" {
		t.Error("Description can be set")
	}

	t.Logf("Task model validated: ID=%s, ColumnID=%s, Title=%s", task.ID, task.ColumnID, task.Title)
}

func TestTaskDeadline(t *testing.T) {
	tests := []struct {
		name     string
		deadline *time.Time
		isSet    bool
	}{
		{
			name:     "Task with deadline",
			deadline: func() *time.Time { t := time.Now().Add(24 * time.Hour); return &t }(),
			isSet:    true,
		},
		{
			name:     "Task without deadline",
			deadline: nil,
			isSet:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Task with deadline",
				Deadline: tt.deadline,
			}

			if tt.isSet && task.Deadline == nil {
				t.Error("Expected Deadline to be set, but it's nil")
			}
			if !tt.isSet && task.Deadline != nil {
				t.Error("Expected Deadline to be nil, but it's set")
			}
			if tt.isSet && task.Deadline != nil && task.Deadline.IsZero() {
				t.Error("Deadline should not be zero time when set")
			}

			t.Logf("Task deadline validated: Deadline=%v", task.Deadline)
		})
	}
}

func TestTaskRelationships(t *testing.T) {
	tests := []struct {
		name           string
		task           Task
		hasComments    bool
		hasLabels      bool
		hasAttachments bool
	}{
		{
			name: "Task with comments, labels, and attachments",
			task: Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Task with all relationships",
				Comments: []Comment{
					{
						ID:      uuid.NewString(),
						TaskID:  uuid.NewString(),
						UserID:  uuid.NewString(),
						Content: "First comment",
					},
					{
						ID:      uuid.NewString(),
						TaskID:  uuid.NewString(),
						UserID:  uuid.NewString(),
						Content: "Second comment",
					},
				},
				Labels: []Label{
					{
						ID:    uuid.NewString(),
						Name:  "Bug",
						Color: "#FF0000",
					},
					{
						ID:    uuid.NewString(),
						Name:  "High Priority",
						Color: "#FFA500",
					},
				},
				Attachments: []Attachment{
					{
						ID:       uuid.NewString(),
						TaskID:   uuid.NewString(),
						FileName: "screenshot.png",
						FileURL:  "https://example.com/screenshot.png",
					},
				},
			},
			hasComments:    true,
			hasLabels:      true,
			hasAttachments: true,
		},
		{
			name: "Task without comments, labels, and attachments (empty slices)",
			task: Task{
				ID:          uuid.NewString(),
				ColumnID:    uuid.NewString(),
				Title:       "Task without relationships",
				Comments:    []Comment{},
				Labels:      []Label{},
				Attachments: []Attachment{},
			},
			hasComments:    false,
			hasLabels:      false,
			hasAttachments: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify comments relationship
			if tt.hasComments && len(tt.task.Comments) == 0 {
				t.Error("Expected Comments to have items, but it's empty")
			}
			if !tt.hasComments && len(tt.task.Comments) > 0 {
				t.Error("Expected Comments to be empty, but it has items")
			}
			if tt.hasComments {
				for _, comment := range tt.task.Comments {
					if comment.ID == "" {
						t.Error("Comment ID should not be empty")
					}
				}
			}

			// Verify labels relationship (many-to-many)
			if tt.hasLabels && len(tt.task.Labels) == 0 {
				t.Error("Expected Labels to have items, but it's empty")
			}
			if !tt.hasLabels && len(tt.task.Labels) > 0 {
				t.Error("Expected Labels to be empty, but it has items")
			}
			if tt.hasLabels {
				for _, label := range tt.task.Labels {
					if label.ID == "" {
						t.Error("Label ID should not be empty")
					}
				}
			}

			// Verify attachments relationship (one-to-many)
			if tt.hasAttachments && len(tt.task.Attachments) == 0 {
				t.Error("Expected Attachments to have items, but it's empty")
			}
			if !tt.hasAttachments && len(tt.task.Attachments) > 0 {
				t.Error("Expected Attachments to be empty, but it has items")
			}
			if tt.hasAttachments {
				for _, attachment := range tt.task.Attachments {
					if attachment.ID == "" {
						t.Error("Attachment ID should not be empty")
					}
				}
			}

			t.Logf("Task relationships validated: Comments=%d, Labels=%d, Attachments=%d",
				len(tt.task.Comments),
				len(tt.task.Labels),
				len(tt.task.Attachments))
		})
	}
}

func TestTaskColumnRelationship(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		hasCol   bool
		colTitle string
	}{
		{
			name: "Task with column",
			task: Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Task with column",
				Column: &Column{
					ID:    uuid.NewString(),
					Title: "To Do",
				},
			},
			hasCol:   true,
			colTitle: "To Do",
		},
		{
			name: "Task without column (nil)",
			task: Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Task without column",
				Column:   nil,
			},
			hasCol: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify column relationship
			if tt.hasCol && tt.task.Column == nil {
				t.Error("Expected Column to be set, but it's nil")
			}
			if !tt.hasCol && tt.task.Column != nil {
				t.Error("Expected Column to be nil, but it's set")
			}
			if tt.hasCol && tt.task.Column != nil && tt.task.Column.ID == "" {
				t.Error("Column ID should not be empty when Column is set")
			}
			if tt.hasCol && tt.task.Column != nil && tt.task.Column.Title != tt.colTitle {
				t.Errorf("Expected Column title '%s', got '%s'", tt.colTitle, tt.task.Column.Title)
			}

			if tt.task.Column != nil {
				t.Logf("Task column relationship validated: ColumnID=%s, ColumnTitle=%s",
					tt.task.Column.ID,
					tt.task.Column.Title)
			}
		})
	}
}

func TestTaskSoftDelete(t *testing.T) {
	task := Task{
		ID:          uuid.NewString(),
		ColumnID:    uuid.NewString(),
		Title:       "Task to be soft deleted",
		Description: "This task will be soft deleted",
		DeletedAt:   gorm.DeletedAt{},
	}

	// Initially not deleted
	if task.DeletedAt.Valid {
		t.Error("Task should not be marked as deleted initially")
	}

	// Simulate soft delete (this would normally be done by GORM)
	task.DeletedAt = gorm.DeletedAt{Time: task.DeletedAt.Time, Valid: true}

	// After soft delete
	if !task.DeletedAt.Valid {
		t.Error("Task should be marked as deleted after soft delete")
	}

	t.Logf("Task soft delete validated: DeletedAt.Valid=%v", task.DeletedAt.Valid)
}

func TestTaskBeforeCreateHook(t *testing.T) {
	tests := []struct {
		name      string
		task      Task
		expectGen bool
	}{
		{
			name: "Task with empty ID should generate UUID",
			task: Task{
				ID:       "",
				ColumnID: uuid.NewString(),
				Title:    "Task without ID",
			},
			expectGen: true,
		},
		{
			name: "Task with existing ID should not generate UUID",
			task: Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Task with existing ID",
			},
			expectGen: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalID := tt.task.ID

			// Call BeforeCreate hook
			err := tt.task.BeforeCreate(nil)
			if err != nil {
				t.Errorf("BeforeCreate hook returned error: %v", err)
			}

			if tt.expectGen {
				if tt.task.ID == "" {
					t.Error("BeforeCreate hook should have generated UUID")
				}
				if originalID != "" && tt.task.ID == originalID {
					t.Error("BeforeCreate hook should have generated new UUID")
				}
			} else {
				if tt.task.ID != originalID {
					t.Error("BeforeCreate hook should not have changed existing ID")
				}
			}

			t.Logf("Task BeforeCreate hook validated: ID=%s (original: %s)",
				tt.task.ID, originalID)
		})
	}
}

func TestTaskFields(t *testing.T) {
	task := Task{
		ID:          uuid.NewString(),
		ColumnID:    uuid.NewString(),
		Title:       "Complete task implementation",
		Description: "Implement the Task model with all required fields and relationships",
		Deadline:    func() *time.Time { t := time.Now().Add(7 * 24 * time.Hour); return &t }(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Verify all fields
	if task.ID == "" {
		t.Error("ID should not be empty")
	}
	if task.ColumnID == "" {
		t.Error("ColumnID should not be empty")
	}
	if task.Title == "" {
		t.Error("Title should not be empty")
	}
	if task.Description == "" {
		t.Error("Description should not be empty")
	}
	if task.Deadline == nil {
		t.Error("Deadline should be set")
	}
	if task.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if task.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	t.Logf("Task fields validated: ID=%s, Title=%s, Description=%s, Deadline=%v",
		task.ID, task.Title, task.Description, task.Deadline)
}

func TestTaskEmptyValidation(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		wantErr bool
	}{
		{
			name:    "Empty ID",
			field:   "ID",
			value:   "",
			wantErr: true,
		},
		{
			name:    "Empty ColumnID",
			field:   "ColumnID",
			value:   "",
			wantErr: true,
		},
		{
			name:    "Empty Title",
			field:   "Title",
			value:   "",
			wantErr: true,
		},
		{
			name:    "Valid Title",
			field:   "Title",
			value:   "Valid Title",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:       uuid.NewString(),
				ColumnID: uuid.NewString(),
				Title:    "Default Title",
			}

			// Set field based on test case
			switch tt.field {
			case "ID":
				task.ID = tt.value
			case "ColumnID":
				task.ColumnID = tt.value
			case "Title":
				task.Title = tt.value
			}

			// Simple validation
			hasEmptyField := task.ID == "" || task.ColumnID == "" || task.Title == ""

			if tt.wantErr && !hasEmptyField {
				t.Errorf("Expected empty field error, but all fields are valid")
			}
			if !tt.wantErr && hasEmptyField {
				t.Errorf("Expected all fields to be valid, but found empty field")
			}

			t.Logf("Task validation: %s='%s', hasEmpty=%v",
				tt.field, tt.value, hasEmptyField)
		})
	}
}
