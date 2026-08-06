package media

// Registry dan DiskManager disimpan sebagai state package mengikuti pola
// config.DB / config.AppConfig pada project ini, supaya resource dan model
// bisa membangun URL tanpa harus mengalirkan dependency ke mana-mana.
var (
	disks       *DiskManager
	collections = NewRegistry()
)

type Config struct {
	// DefaultDisk adalah nama disk yang dipakai bila collection tidak
	// menentukan disk secara eksplisit.
	DefaultDisk string
	// LocalRoot adalah folder tempat file disimpan, mis. "storage/media".
	LocalRoot string
	// LocalBaseURL adalah prefix URL publik untuk LocalRoot, mis. "/storage/media".
	LocalBaseURL string
}

const DiskLocal = "local"

// Boot menyiapkan disk yang tersedia. Dipanggil sekali saat aplikasi start.
func Boot(cfg Config) {
	manager := NewDiskManager(cfg.DefaultDisk)
	manager.Register(NewLocalDisk(DiskLocal, cfg.LocalRoot, cfg.LocalBaseURL))
	disks = manager
}

// Disks mengembalikan disk manager aktif.
func Disks() *DiskManager {
	return disks
}

// Register mendaftarkan collection milik sebuah model type.
func Register(modelType string, cols ...Collection) {
	collections.Register(modelType, cols...)
}

// URLOf mengembalikan URL publik file asli.
func URLOf(m *Media) string {
	disk, err := disks.Disk(m.Disk)
	if err != nil {
		return ""
	}
	return disk.URL(m.Path())
}

// ConversionURLOf mengembalikan URL publik sebuah conversion. String kosong
// berarti conversion tersebut tidak ada untuk media ini.
func ConversionURLOf(m *Media, conversion string) string {
	relPath := m.ConversionPath(conversion)
	if relPath == "" {
		return ""
	}

	disk, err := disks.Disk(m.Disk)
	if err != nil {
		return ""
	}
	return disk.URL(relPath)
}
