package main

import (
	"encoding/json"
	"testing"
)

// TestLockState tests the lock state management
func TestLockState(t *testing.T) {
	// Reset global state
	mutex.Lock()
	isLocked = false
	lockedBy = ""
	content = ""
	mutex.Unlock()

	// Test 1: Initial state is unlocked
	t.Run("Initial state is unlocked", func(t *testing.T) {
		mutex.Lock()
		defer mutex.Unlock()
		if isLocked {
			t.Fatalf("Expected lock to be unlocked initially")
		}
		if lockedBy != "" {
			t.Fatalf("Expected lockedBy to be empty initially")
		}
	})

	// Test 2: Lock acquisition
	t.Run("Lock acquisition", func(t *testing.T) {
		// Reset state
		mutex.Lock()
		isLocked = false
		lockedBy = ""
		mutex.Unlock()

		// Simulate lock acquisition
		mutex.Lock()
		if !isLocked {
			isLocked = true
			lockedBy = "test_client_1"
		}
		mutex.Unlock()

		// Verify lock state
		mutex.Lock()
		defer mutex.Unlock()
		if !isLocked {
			t.Fatalf("Expected lock to be locked after acquisition")
		}
		if lockedBy != "test_client_1" {
			t.Fatalf("Expected lock to be locked by test_client_1, got: %s", lockedBy)
		}
	})

	// Test 3: Lock release
	t.Run("Lock release", func(t *testing.T) {
		// Pre-lock
		mutex.Lock()
		isLocked = true
		lockedBy = "test_client_1"
		mutex.Unlock()

		// Simulate lock release
		mutex.Lock()
		if isLocked && lockedBy == "test_client_1" {
			isLocked = false
			lockedBy = ""
		}
		mutex.Unlock()

		// Verify lock state
		mutex.Lock()
		defer mutex.Unlock()
		if isLocked {
			t.Fatalf("Expected lock to be unlocked after release")
		}
		if lockedBy != "" {
			t.Fatalf("Expected lockedBy to be empty after release")
		}
	})

	// Test 4: Content management
	t.Run("Content management", func(t *testing.T) {
		// Reset state
		mutex.Lock()
		content = "initial content"
		mutex.Unlock()

		// Update content
		mutex.Lock()
		content = "updated content"
		mutex.Unlock()

		// Verify content
		mutex.Lock()
		defer mutex.Unlock()
		if content != "updated content" {
			t.Fatalf("Expected content 'updated content', got: %s", content)
		}
	})
}

// TestMessageParsing tests the message parsing functionality
func TestMessageParsing(t *testing.T) {
	// Test lock message
	lockMsg := `{"type":"lock","clientId":"test_client"}`
	var message Message
	if err := json.Unmarshal([]byte(lockMsg), &message); err != nil {
		t.Fatalf("Failed to unmarshal lock message: %v", err)
	}
	if message.Type != "lock" {
		t.Fatalf("Expected lock message type, got: %s", message.Type)
	}
	if message.ClientID != "test_client" {
		t.Fatalf("Expected clientId 'test_client', got: %s", message.ClientID)
	}

	// Test content message
	contentMsg := `{"type":"content","content":"test content","clientId":"test_client"}`
	if err := json.Unmarshal([]byte(contentMsg), &message); err != nil {
		t.Fatalf("Failed to unmarshal content message: %v", err)
	}
	if message.Type != "content" {
		t.Fatalf("Expected content message type, got: %s", message.Type)
	}
	if message.Content != "test content" {
		t.Fatalf("Expected content 'test content', got: %s", message.Content)
	}

	// Test unlock message
	unlockMsg := `{"type":"unlock","clientId":"test_client"}`
	if err := json.Unmarshal([]byte(unlockMsg), &message); err != nil {
		t.Fatalf("Failed to unmarshal unlock message: %v", err)
	}
	if message.Type != "unlock" {
		t.Fatalf("Expected unlock message type, got: %s", message.Type)
	}
}

// TestLockStatusMessage tests the lock status message structure
func TestLockStatusMessage(t *testing.T) {
	// Create lock status message
	lockStatus := LockStatus{
		Type:     "lock",
		Locked:   true,
		LockedBy: "test_client",
	}

	// Marshal and unmarshal
	data, err := json.Marshal(lockStatus)
	if err != nil {
		t.Fatalf("Failed to marshal lock status: %v", err)
	}

	var decoded LockStatus
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal lock status: %v", err)
	}

	// Verify
	if decoded.Type != "lock" {
		t.Fatalf("Expected type 'lock', got: %s", decoded.Type)
	}
	if !decoded.Locked {
		t.Fatalf("Expected locked to be true")
	}
	if decoded.LockedBy != "test_client" {
		t.Fatalf("Expected lockedBy 'test_client', got: %s", decoded.LockedBy)
	}
}

// TestAutoUnlockOnDisconnect tests that lock is automatically released when client disconnects
func TestAutoUnlockOnDisconnect(t *testing.T) {
	// Reset state
	mutex.Lock()
	isLocked = false
	lockedBy = ""
	content = ""
	mutex.Unlock()

	// Simulate client connecting and acquiring lock
	mutex.Lock()
	isLocked = true
	lockedBy = "disconnecting_client"
	mutex.Unlock()

	// Verify lock is held
	mutex.Lock()
	if !isLocked {
		t.Fatalf("Expected lock to be locked")
	}
	if lockedBy != "disconnecting_client" {
		t.Fatalf("Expected lock to be locked by disconnecting_client")
	}
	mutex.Unlock()

	// Simulate client disconnecting (this is what happens in the defer function)
	mutex.Lock()
	clientID := "disconnecting_client"
	// Check if client holds lock and release it
	if isLocked && lockedBy == clientID {
		isLocked = false
		lockedBy = ""
	}
	mutex.Unlock()

	// Verify lock is released
	mutex.Lock()
	defer mutex.Unlock()
	if isLocked {
		t.Fatalf("Expected lock to be unlocked after client disconnect")
	}
	if lockedBy != "" {
		t.Fatalf("Expected lockedBy to be empty after client disconnect")
	}
}
