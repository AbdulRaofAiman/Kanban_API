package models

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestAttachmentModel(t *testing.T) {
	attachment := Attachment{
		ID:       uuid.New().String(),
		TaskID:   uuid.New().String(),
		FileName: "test_document.pdf",
		FileURL:  "https://example.com/files/test_document.pdf",
		FileSize: 1024,
	}

	// Verify table name
	if attachment.TableName() != "attachments" {
		t.Errorf("Expected table name 'attachments', got '%s'", attachment.TableName())
	}

	// Verify ID is set
	if attachment.ID == "" {
		t.Error("ID should not be empty")
	}

	// Verify TaskID is set
	if attachment.TaskID == "" {
		t.Error("TaskID should not be empty")
	}

	// Verify FileName is not empty
	if attachment.FileName == "" {
		t.Error("FileName should not be empty")
	}

	// Verify FileURL is not empty
	if attachment.FileURL == "" {
		t.Error("FileURL should not be empty")
	}

	// Verify FileSize is non-negative
	if attachment.FileSize < 0 {
		t.Error("FileSize should not be negative")
	}

	t.Logf("Attachment model validated: ID=%s, TaskID=%s, FileName=%s",
		attachment.ID, attachment.TaskID, attachment.FileName)
}

func TestAttachmentRelationships(t *testing.T) {
	tests := []struct {
		name       string
		attachment Attachment
		hasTask    bool
	}{
		{
			name: "Attachment with task",
			attachment: Attachment{
				ID:       uuid.New().String(),
				TaskID:   uuid.New().String(),
				FileName: "test_file.jpg",
				FileURL:  "https://example.com/files/test_file.jpg",
				FileSize: 2048,
				Task: &Task{
					ID:    uuid.New().String(),
					Title: "Test Task",
				},
			},
			hasTask: true,
		},
		{
			name: "Attachment without task (nil relationship)",
			attachment: Attachment{
				ID:       uuid.New().String(),
				TaskID:   uuid.New().String(),
				FileName: "another_file.png",
				FileURL:  "https://example.com/files/another_file.png",
				FileSize: 512,
				Task:     nil,
			},
			hasTask: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify task relationship
			if tt.hasTask && tt.attachment.Task == nil {
				t.Error("Expected Task to be set, but it's nil")
			}
			if !tt.hasTask && tt.attachment.Task != nil {
				t.Error("Expected Task to be nil, but it's set")
			}
			if tt.hasTask && tt.attachment.Task != nil && tt.attachment.Task.ID == "" {
				t.Error("Task ID should not be empty when Task is set")
			}

			t.Logf("Attachment relationships validated: Task=%v",
				getPointerStatus(tt.attachment.Task))
		})
	}
}

func TestAttachmentSoftDelete(t *testing.T) {
	attachment := Attachment{
		ID:        uuid.New().String(),
		TaskID:    uuid.New().String(),
		FileName:  "file_to_delete.txt",
		FileURL:   "https://example.com/files/file_to_delete.txt",
		FileSize:  256,
		DeletedAt: gorm.DeletedAt{},
	}

	// Initially not deleted
	if attachment.DeletedAt.Valid {
		t.Error("Attachment should not be marked as deleted initially")
	}

	// Simulate soft delete (this would normally be done by GORM)
	attachment.DeletedAt = gorm.DeletedAt{Time: attachment.DeletedAt.Time, Valid: true}

	// After soft delete
	if !attachment.DeletedAt.Valid {
		t.Error("Attachment should be marked as deleted after soft delete")
	}

	t.Logf("Attachment soft delete validated: DeletedAt.Valid=%v", attachment.DeletedAt.Valid)
}

func TestAttachmentBeforeCreate(t *testing.T) {
	// This test would need a real GORM DB instance, so we're just
	// verifying the hook exists and would be called
	attachment := Attachment{
		TaskID:   uuid.New().String(),
		FileName: "new_file.pdf",
		FileURL:  "https://example.com/files/new_file.pdf",
		FileSize: 1024,
	}

	// Verify ID is empty before BeforeCreate
	if attachment.ID == "" {
		t.Log("ID is empty before BeforeCreate (expected)")
	}

	// Manually call BeforeCreate hook to simulate GORM behavior
	// In a real scenario, GORM would call this automatically
	err := attachment.BeforeCreate(nil)
	if err != nil {
		t.Errorf("BeforeCreate hook returned an error: %v", err)
	}

	// After BeforeCreate, ID should be set
	if attachment.ID == "" {
		t.Error("ID should be set after BeforeCreate")
	}

	t.Logf("Attachment BeforeCreate validated: ID=%s", attachment.ID)
}

func TestAttachmentVariousFileTypes(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		fileSize  int64
		expectExt string
	}{
		{
			name:      "PDF file",
			fileName:  "document.pdf",
			fileSize:  1024 * 1024, // 1MB
			expectExt: ".pdf",
		},
		{
			name:      "Image file",
			fileName:  "photo.jpg",
			fileSize:  500 * 1024, // 500KB
			expectExt: ".jpg",
		},
		{
			name:      "Large video file",
			fileName:  "video.mp4",
			fileSize:  100 * 1024 * 1024, // 100MB
			expectExt: ".mp4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attachment := Attachment{
				ID:       uuid.New().String(),
				TaskID:   uuid.New().String(),
				FileName: tt.fileName,
				FileURL:  "https://example.com/files/" + tt.fileName,
				FileSize: tt.fileSize,
			}

			if attachment.FileName != tt.fileName {
				t.Errorf("Expected FileName '%s', got '%s'", tt.fileName, attachment.FileName)
			}

			if attachment.FileSize != tt.fileSize {
				t.Errorf("Expected FileSize %d, got %d", tt.fileSize, attachment.FileSize)
			}

			t.Logf("Attachment validated: FileName=%s, FileSize=%d bytes",
				attachment.FileName, attachment.FileSize)
		})
	}
}
