package rocketmq

const (
	DEFAULT_RETRY        = 3
	DEFAULT_MAX_CACHE    = 16
	DEFAULT_STORAGE_FILE = "./storagefile"
)

type Config struct {
	Address     []string
	SecretKey   string
	AccessKey   string
	NameSpace   string
	GroupName   string
	Retry       int
	Topic       string
	MaxCache    int
	StorageFile string
}

func newDefaultConfig() *Config {
	return &Config{
		Retry:       DEFAULT_RETRY,
		MaxCache:    DEFAULT_MAX_CACHE,
		StorageFile: DEFAULT_STORAGE_FILE,
	}
}
