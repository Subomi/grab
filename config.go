package grab

/*
	Url to download information
*/

type OverWriteableConfig struct {
	CacheSize int64
	CacheDir  string
	HomeDir   string
}

type SystemConfig struct {
	OverWriteableConfig
}

type UserConfig struct {
	OverWriteableConfig
}

