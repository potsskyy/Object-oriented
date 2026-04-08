package repository

import (
	"testing"
	"time"

	"todo/internal/models"
)

func TestMemoryRepoSessionLifecycle(t *testing.T) {
	repo := NewMemoryRepo()

	if _, err := repo.CreateSession("missing"); err == nil {
		t.Fatal("expected error when creating session for unknown user")
	}

	if err := repo.Register(&models.User{
		Username: "admin",
		Password: "hashed-password",
	}); err != nil {
		t.Fatalf("register user: %v", err)
	}

	sessionID, err := repo.CreateSession("admin")
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if len(sessionID) != 64 {
		t.Fatalf("session id length = %d, want 64", len(sessionID))
	}
	if sessionID == "admin" {
		t.Fatal("session id should be opaque, not equal to username")
	}

	username, err := repo.GetUsernameBySession(sessionID)
	if err != nil {
		t.Fatalf("get username by session: %v", err)
	}
	if username != "admin" {
		t.Fatalf("username = %q, want %q", username, "admin")
	}

	savedSession := repo.sessions[sessionID]
	savedSession.ExpiresAt = time.Now().Add(-time.Minute)
	repo.sessions[sessionID] = savedSession

	if _, err := repo.GetUsernameBySession(sessionID); err != ErrInvalidSession {
		t.Fatalf("expired session error = %v, want %v", err, ErrInvalidSession)
	}

	sessionID, err = repo.CreateSession("admin")
	if err != nil {
		t.Fatalf("create session after expiration: %v", err)
	}

	repo.DeleteSession(sessionID)
	if _, err := repo.GetUsernameBySession(sessionID); err != ErrInvalidSession {
		t.Fatalf("deleted session error = %v, want %v", err, ErrInvalidSession)
	}
}

func TestMemoryRepoTaskLifecycle(t *testing.T) {
	repo := NewMemoryRepo()

	if err := repo.Register(&models.User{
		Username: "student",
		Password: "hashed-password",
	}); err != nil {
		t.Fatalf("register user: %v", err)
	}

	taskID, err := repo.Add("student", &models.Task{
		Headline: "write tests",
		Details:  "cover main scenarios",
	})
	if err != nil {
		t.Fatalf("add task: %v", err)
	}
	if taskID != 1 {
		t.Fatalf("task id = %d, want 1", taskID)
	}

	activeTasks, err := repo.Get("student")
	if err != nil {
		t.Fatalf("get active tasks: %v", err)
	}
	if len(activeTasks) != 1 {
		t.Fatalf("active tasks len = %d, want 1", len(activeTasks))
	}

	if err := repo.Update("student", &models.Task{
		ID:       taskID,
		Headline: "write more tests",
		Details:  "cover auth middleware too",
	}); err != nil {
		t.Fatalf("update task: %v", err)
	}

	savedTask, err := repo.GetByID("student", taskID)
	if err != nil {
		t.Fatalf("get task by id: %v", err)
	}
	if savedTask.Headline != "write more tests" {
		t.Fatalf("headline = %q, want %q", savedTask.Headline, "write more tests")
	}

	if err := repo.Resolve("student", taskID); err != nil {
		t.Fatalf("resolve task: %v", err)
	}

	activeTasks, err = repo.Get("student")
	if err != nil {
		t.Fatalf("get active tasks after resolve: %v", err)
	}
	if len(activeTasks) != 0 {
		t.Fatalf("active tasks after resolve len = %d, want 0", len(activeTasks))
	}

	archivedTasks, err := repo.GetArchive("student")
	if err != nil {
		t.Fatalf("get archive: %v", err)
	}
	if len(archivedTasks) != 1 {
		t.Fatalf("archived tasks len = %d, want 1", len(archivedTasks))
	}
	if !archivedTasks[0].Completed || !archivedTasks[0].Archived {
		t.Fatal("resolved task must be archived and completed")
	}
}
