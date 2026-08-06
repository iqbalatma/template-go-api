package media

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Disk adalah abstraksi penyimpanan file. Implementasi baru (S3, GCS, dsb)
// cukup memenuhi interface ini tanpa mengubah repository.
type Disk interface {
	Name() string
	Put(relPath string, src io.Reader) error
	Open(relPath string) (io.ReadCloser, error)
	Delete(relPath string) error
	DeleteDirectory(relDir string) error
	URL(relPath string) string
}

// LocalDisk menyimpan file di filesystem lokal dan menyajikannya lewat baseURL.
type LocalDisk struct {
	name    string
	root    string
	baseURL string
}

func NewLocalDisk(name, root, baseURL string) *LocalDisk {
	return &LocalDisk{
		name:    name,
		root:    root,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (d *LocalDisk) Name() string {
	return d.name
}

func (d *LocalDisk) Root() string {
	return d.root
}

// resolve mencegah path traversal keluar dari root disk.
func (d *LocalDisk) resolve(relPath string) (string, error) {
	cleaned := filepath.Clean(filepath.Join(d.root, filepath.FromSlash(relPath)))

	root, err := filepath.Abs(d.root)
	if err != nil {
		return "", err
	}
	target, err := filepath.Abs(cleaned)
	if err != nil {
		return "", err
	}

	if target != root && !strings.HasPrefix(target, root+string(os.PathSeparator)) {
		return "", fmt.Errorf("media: path %q is outside disk root", relPath)
	}
	return target, nil
}

func (d *LocalDisk) Put(relPath string, src io.Reader) error {
	target, err := d.resolve(relPath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	dst, err := os.Create(target)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}
	return nil
}

func (d *LocalDisk) Open(relPath string) (io.ReadCloser, error) {
	target, err := d.resolve(relPath)
	if err != nil {
		return nil, err
	}
	return os.Open(target)
}

func (d *LocalDisk) Delete(relPath string) error {
	target, err := d.resolve(relPath)
	if err != nil {
		return err
	}
	if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (d *LocalDisk) DeleteDirectory(relDir string) error {
	target, err := d.resolve(relDir)
	if err != nil {
		return err
	}

	root, err := filepath.Abs(d.root)
	if err != nil {
		return err
	}
	// Jangan pernah menghapus root disk itu sendiri.
	if target == root {
		return fmt.Errorf("media: refusing to delete disk root")
	}

	if err := os.RemoveAll(target); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (d *LocalDisk) URL(relPath string) string {
	return d.baseURL + "/" + strings.TrimLeft(path.Clean(filepath.ToSlash(relPath)), "/")
}

// DiskManager menampung disk yang terdaftar beserta disk default.
type DiskManager struct {
	disks       map[string]Disk
	defaultDisk string
}

func NewDiskManager(defaultDisk string) *DiskManager {
	return &DiskManager{
		disks:       make(map[string]Disk),
		defaultDisk: defaultDisk,
	}
}

func (m *DiskManager) Register(disk Disk) {
	m.disks[disk.Name()] = disk
}

func (m *DiskManager) Disk(name string) (Disk, error) {
	if name == "" {
		name = m.defaultDisk
	}
	disk, ok := m.disks[name]
	if !ok {
		return nil, fmt.Errorf("media: disk %q is not registered", name)
	}
	return disk, nil
}

func (m *DiskManager) DefaultDiskName() string {
	return m.defaultDisk
}
