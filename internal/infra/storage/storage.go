package storage

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrInvalidPath  = errors.New("invalid path")
)

// FileInfo represents file metadata
type FileInfo struct {
	Path         string    `json:"path"`
	Name         string    `json:"name"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	IsDir        bool      `json:"is_dir"`
	MimeType     string    `json:"mime_type,omitempty"`
}

// Driver defines the storage driver interface
type Driver interface {
	// Put stores a file
	Put(ctx context.Context, path string, content []byte) error

	// PutStream stores a file from a reader
	PutStream(ctx context.Context, path string, reader io.Reader) error

	// Get retrieves a file's content
	Get(ctx context.Context, path string) ([]byte, error)

	// GetStream retrieves a file as a reader
	GetStream(ctx context.Context, path string) (io.ReadCloser, error)

	// Delete removes a file
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists
	Exists(ctx context.Context, path string) bool

	// Size returns the file size
	Size(ctx context.Context, path string) (int64, error)

	// LastModified returns the last modification time
	LastModified(ctx context.Context, path string) (time.Time, error)

	// Copy copies a file
	Copy(ctx context.Context, from, to string) error

	// Move moves a file
	Move(ctx context.Context, from, to string) error

	// URL returns the public URL for a file (if applicable)
	URL(path string) string

	// Files lists files in a directory
	Files(ctx context.Context, directory string) ([]FileInfo, error)

	// AllFiles lists all files recursively
	AllFiles(ctx context.Context, directory string) ([]FileInfo, error)

	// Directories lists directories
	Directories(ctx context.Context, directory string) ([]string, error)

	// MakeDirectory creates a directory
	MakeDirectory(ctx context.Context, path string) error

	// DeleteDirectory removes a directory
	DeleteDirectory(ctx context.Context, path string) error
}

// --- Storage Manager ---

// Manager manages storage disks
type Manager struct {
	mu       sync.RWMutex
	disks    map[string]Driver
	default_ string
}

var (
	manager *Manager
	once    sync.Once
)

// Global returns the global storage manager
func Global() *Manager {
	once.Do(func() {
		manager = New()
	})
	return manager
}

// New creates a new storage manager
func New() *Manager {
	return &Manager{
		disks:    make(map[string]Driver),
		default_: "local",
	}
}

// RegisterDisk registers a storage disk
func (m *Manager) RegisterDisk(name string, driver Driver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.disks[name] = driver
}

// SetDefault sets the default disk
func (m *Manager) SetDefault(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.default_ = name
}

// Disk returns a disk by name
func (m *Manager) Disk(name string) Driver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.disks[name]
}

// Default returns the default disk
func (m *Manager) Default() Driver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.disks[m.default_]
}

// --- Local Disk Driver ---

// LocalDisk implements local filesystem storage
type LocalDisk struct {
	root    string
	baseURL string
}

// LocalConfig holds local disk configuration
type LocalConfig struct {
	Root    string // Root directory for storage
	BaseURL string // Base URL for public files
}

// NewLocalDisk creates a new local disk driver
func NewLocalDisk(cfg LocalConfig) (*LocalDisk, error) {
	root := cfg.Root
	if root == "" {
		root = "./storage"
	}

	// Ensure root directory exists
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, err
	}

	return &LocalDisk{
		root:    root,
		baseURL: cfg.BaseURL,
	}, nil
}

// fullPath returns the full filesystem path
func (d *LocalDisk) fullPath(path string) string {
	// Prevent directory traversal
	clean := filepath.Clean(path)
	if strings.HasPrefix(clean, "..") {
		return ""
	}
	return filepath.Join(d.root, clean)
}

// Put stores a file
func (d *LocalDisk) Put(ctx context.Context, path string, content []byte) error {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return ErrInvalidPath
	}

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(fullPath, content, 0644)
}

// PutStream stores a file from a reader
func (d *LocalDisk) PutStream(ctx context.Context, path string, reader io.Reader) error {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return ErrInvalidPath
	}

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	return err
}

// Get retrieves a file's content
func (d *LocalDisk) Get(ctx context.Context, path string) ([]byte, error) {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return nil, ErrInvalidPath
	}

	content, err := os.ReadFile(fullPath)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	return content, err
}

// GetStream retrieves a file as a reader
func (d *LocalDisk) GetStream(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return nil, ErrInvalidPath
	}

	file, err := os.Open(fullPath)
	if os.IsNotExist(err) {
		return nil, ErrFileNotFound
	}
	return file, err
}

// Delete removes a file
func (d *LocalDisk) Delete(ctx context.Context, path string) error {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return ErrInvalidPath
	}

	err := os.Remove(fullPath)
	if os.IsNotExist(err) {
		return nil // Not an error if file doesn't exist
	}
	return err
}

// Exists checks if a file exists
func (d *LocalDisk) Exists(ctx context.Context, path string) bool {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return false
	}

	_, err := os.Stat(fullPath)
	return err == nil
}

// Size returns the file size
func (d *LocalDisk) Size(ctx context.Context, path string) (int64, error) {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return 0, ErrInvalidPath
	}

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return 0, ErrFileNotFound
	}
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// LastModified returns the last modification time
func (d *LocalDisk) LastModified(ctx context.Context, path string) (time.Time, error) {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return time.Time{}, ErrInvalidPath
	}

	info, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		return time.Time{}, ErrFileNotFound
	}
	if err != nil {
		return time.Time{}, err
	}
	return info.ModTime(), nil
}

// Copy copies a file
func (d *LocalDisk) Copy(ctx context.Context, from, to string) error {
	content, err := d.Get(ctx, from)
	if err != nil {
		return err
	}
	return d.Put(ctx, to, content)
}

// Move moves a file
func (d *LocalDisk) Move(ctx context.Context, from, to string) error {
	fromPath := d.fullPath(from)
	toPath := d.fullPath(to)
	if fromPath == "" || toPath == "" {
		return ErrInvalidPath
	}

	// Ensure destination directory exists
	dir := filepath.Dir(toPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.Rename(fromPath, toPath)
}

// URL returns the public URL for a file
func (d *LocalDisk) URL(path string) string {
	if d.baseURL == "" {
		return path
	}
	return strings.TrimSuffix(d.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")
}

// Files lists files in a directory
func (d *LocalDisk) Files(ctx context.Context, directory string) ([]FileInfo, error) {
	fullPath := d.fullPath(directory)
	if fullPath == "" {
		return nil, ErrInvalidPath
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Path:         filepath.Join(directory, entry.Name()),
			Name:         entry.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
			IsDir:        false,
		})
	}

	return files, nil
}

// AllFiles lists all files recursively
func (d *LocalDisk) AllFiles(ctx context.Context, directory string) ([]FileInfo, error) {
	fullPath := d.fullPath(directory)
	if fullPath == "" {
		return nil, ErrInvalidPath
	}

	var files []FileInfo

	err := filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(d.root, path)
		files = append(files, FileInfo{
			Path:         relPath,
			Name:         info.Name(),
			Size:         info.Size(),
			LastModified: info.ModTime(),
			IsDir:        false,
		})
		return nil
	})

	return files, err
}

// Directories lists directories
func (d *LocalDisk) Directories(ctx context.Context, directory string) ([]string, error) {
	fullPath := d.fullPath(directory)
	if fullPath == "" {
		return nil, ErrInvalidPath
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotFound
		}
		return nil, err
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// MakeDirectory creates a directory
func (d *LocalDisk) MakeDirectory(ctx context.Context, path string) error {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return ErrInvalidPath
	}
	return os.MkdirAll(fullPath, 0755)
}

// DeleteDirectory removes a directory
func (d *LocalDisk) DeleteDirectory(ctx context.Context, path string) error {
	fullPath := d.fullPath(path)
	if fullPath == "" {
		return ErrInvalidPath
	}
	return os.RemoveAll(fullPath)
}

// --- Convenience Functions ---

// RegisterDisk registers a disk with the global manager
func RegisterDisk(name string, driver Driver) {
	Global().RegisterDisk(name, driver)
}

// Disk returns a disk from the global manager
func Disk(name string) Driver {
	return Global().Disk(name)
}

// Default returns the default disk
func Default() Driver {
	return Global().Default()
}

// Put stores a file on the default disk
func Put(ctx context.Context, path string, content []byte) error {
	return Default().Put(ctx, path, content)
}

// Get retrieves a file from the default disk
func Get(ctx context.Context, path string) ([]byte, error) {
	return Default().Get(ctx, path)
}

// Delete removes a file from the default disk
func Delete(ctx context.Context, path string) error {
	return Default().Delete(ctx, path)
}

// Exists checks if a file exists on the default disk
func Exists(ctx context.Context, path string) bool {
	return Default().Exists(ctx, path)
}

// URL returns the URL for a file on the default disk
func URL(path string) string {
	return Default().URL(path)
}
