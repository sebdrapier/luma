package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"elano.fr/src/backend/models"
	"github.com/google/uuid"
)

type ProjectStore interface {
	Get() *models.Project
	Save(*models.Project) error
	Delete() error
	Backup() error
	GetPath() string
}

type YAMLStore struct {
	mu           sync.RWMutex
	project      *models.Project
	filepath     string
	backupDir    string
	maxBackups   int
	lastModified time.Time
}

func NewYAMLStore(path string) (*YAMLStore, error) {
	project, err := LoadProjectFromFile(path)
	if err != nil {
		return nil, err
	}

	store := &YAMLStore{
		project:      project,
		filepath:     path,
		backupDir:    filepath.Join(filepath.Dir(path), "backups"),
		maxBackups:   10,
		lastModified: time.Now(),
	}

	if err := os.MkdirAll(store.backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return store, nil
}

func NewYAMLStoreWithDefault(path string) (*YAMLStore, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	project := &models.Project{
		ID:           uuid.New().String(),
		Name:         "Default Project",
		USBInterface: "/dev/ttyUSB0",
		Fixtures:     []models.Fixture{},
		Presets:      []models.Preset{},
		Shows:        []models.Show{},
	}

	if err := SaveProjectToFile(project, path); err != nil {
		return nil, fmt.Errorf("failed to save default project: %w", err)
	}

	return NewYAMLStore(path)
}

func (s *YAMLStore) Get() *models.Project {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.project == nil {
		return nil
	}

	projectCopy := *s.project

	projectCopy.Fixtures = make([]models.Fixture, len(s.project.Fixtures))
	copy(projectCopy.Fixtures, s.project.Fixtures)

	projectCopy.Presets = make([]models.Preset, len(s.project.Presets))
	copy(projectCopy.Presets, s.project.Presets)

	projectCopy.Shows = make([]models.Show, len(s.project.Shows))
	copy(projectCopy.Shows, s.project.Shows)

	return &projectCopy
}

func (s *YAMLStore) Save(p *models.Project) error {
	if p == nil {
		return fmt.Errorf("cannot save nil project")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.project != nil && time.Since(s.lastModified) > 5*time.Minute {
		if err := s.createBackup(); err != nil {
			fmt.Printf("Warning: failed to create backup: %v\n", err)
		}
	}

	if err := s.validateProject(p); err != nil {
		return fmt.Errorf("project validation failed: %w", err)
	}

	tempFile := s.filepath + ".tmp"
	if err := SaveProjectToFile(p, tempFile); err != nil {
		return fmt.Errorf("failed to save to temporary file: %w", err)
	}

	if err := os.Rename(tempFile, s.filepath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to save project: %w", err)
	}

	s.project = p
	s.lastModified = time.Now()
	return nil
}

func (s *YAMLStore) Delete() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.project != nil {
		backupPath := filepath.Join(s.backupDir, fmt.Sprintf("deleted_%s.yaml", time.Now().Format("20060102_150405")))
		SaveProjectToFile(s.project, backupPath)
	}

	if err := os.Remove(s.filepath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete project file: %w", err)
	}

	s.project = nil
	return nil
}

func (s *YAMLStore) Backup() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.project == nil {
		return fmt.Errorf("no project to backup")
	}

	return s.createBackup()
}

func (s *YAMLStore) GetPath() string {
	return s.filepath
}

func (s *YAMLStore) createBackup() error {
	if s.project == nil {
		return nil
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(s.backupDir, fmt.Sprintf("backup_%s.yaml", timestamp))

	if err := SaveProjectToFile(s.project, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	go s.cleanOldBackups()

	return nil
}

func (s *YAMLStore) cleanOldBackups() {
	files, err := os.ReadDir(s.backupDir)
	if err != nil {
		return
	}

	var backups []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
			backups = append(backups, file)
		}
	}

	if len(backups) <= s.maxBackups {
		return
	}

	for i := 0; i < len(backups)-s.maxBackups; i++ {
		os.Remove(filepath.Join(s.backupDir, backups[i].Name()))
	}
}

func (s *YAMLStore) validateProject(p *models.Project) error {
	if p.ID == "" {
		return fmt.Errorf("project ID is required")
	}
	if p.Name == "" {
		return fmt.Errorf("project name is required")
	}
	if p.USBInterface == "" {
		return fmt.Errorf("USB interface is required")
	}

	if p.Fixtures == nil {
		p.Fixtures = []models.Fixture{}
	}
	if p.Presets == nil {
		p.Presets = []models.Preset{}
	}
	if p.Shows == nil {
		p.Shows = []models.Show{}
	}

	fixtureIDs := make(map[string]bool)
	for _, f := range p.Fixtures {
		if f.ID == "" {
			return fmt.Errorf("fixture ID cannot be empty")
		}
		if fixtureIDs[f.ID] {
			return fmt.Errorf("duplicate fixture ID: %s", f.ID)
		}
		fixtureIDs[f.ID] = true
	}

	presetIDs := make(map[string]bool)
	for _, pr := range p.Presets {
		if pr.ID == "" {
			return fmt.Errorf("preset ID cannot be empty")
		}
		if presetIDs[pr.ID] {
			return fmt.Errorf("duplicate preset ID: %s", pr.ID)
		}
		presetIDs[pr.ID] = true
	}

	showIDs := make(map[string]bool)
	for _, s := range p.Shows {
		if s.ID == "" {
			return fmt.Errorf("show ID cannot be empty")
		}
		if showIDs[s.ID] {
			return fmt.Errorf("duplicate show ID: %s", s.ID)
		}
		showIDs[s.ID] = true
	}

	return nil
}
