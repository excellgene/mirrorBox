package smb

import (
	"context"
	"io"

	"excellgene.com/symbaSync/internal/infra/fs"
)

// Client provides an abstraction for SMB operations.
// Responsibility: Define interface for remote filesystem access.
// Implementations should handle SMB protocol details, authentication, etc.
//
// This interface keeps sync logic decoupled from SMB implementation.
// Tests can use mock implementations.
type Client interface {
	// Connect establishes connection to SMB share.
	// Must be called before other operations.
	Connect(ctx context.Context) error

	// Disconnect closes the SMB connection.
	Disconnect() error

	// Walk traverses the remote directory tree.
	// Returns a Walker that can enumerate remote files.
	Walk(sharePath string) (fs.Walker, error)

	// Upload copies a file to the remote share.
	// dst is the remote path, src provides the file data.
	Upload(ctx context.Context, dst string, src io.Reader, size int64) error

	// Delete removes a file from the remote share.
	Delete(ctx context.Context, path string) error

	// MkdirAll creates a directory and all necessary parents.
	MkdirAll(ctx context.Context, path string) error
}

// Config holds SMB connection parameters.
type Config struct {
	Host     string // SMB server hostname or IP
	Port     int    // SMB server port (usually 445)
	Share    string // Share name
	Username string // Authentication username
	Password string // Authentication password
	Domain   string // Windows domain (optional)
}

// MockClient is a placeholder implementation for testing and development.
// Replace with real SMB client (e.g., hirochachacha/go-smb2) in production.
type MockClient struct {
	connected bool
}

// NewMockClient creates a mock SMB client for testing.
func NewMockClient(cfg Config) *MockClient {
	return &MockClient{}
}

func (m *MockClient) Connect(ctx context.Context) error {
	m.connected = true
	return nil
}

func (m *MockClient) Disconnect() error {
	m.connected = false
	return nil
}

func (m *MockClient) Walk(sharePath string) (fs.Walker, error) {
	// In real implementation, return an SMB walker
	// For now, return nil (caller should handle gracefully)
	return nil, nil
}

func (m *MockClient) Upload(ctx context.Context, dst string, src io.Reader, size int64) error {
	// Placeholder: in real impl, copy src to SMB destination
	return nil
}

func (m *MockClient) Delete(ctx context.Context, path string) error {
	// Placeholder: in real impl, delete from SMB share
	return nil
}

func (m *MockClient) MkdirAll(ctx context.Context, path string) error {
	// Placeholder: in real impl, create directory on SMB share
	return nil
}
