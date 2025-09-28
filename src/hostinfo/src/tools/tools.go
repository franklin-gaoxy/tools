package tools

/*
This is the global configuration or global variable
*/

type ServiceConfig struct {
	Port     int16 `yaml:"port"`
	Database struct {
		DataBaseType string `yaml:"databaseType"`
		ConnPath     string `yaml:"connPath"`
		//Type         string     `yaml:"type"`
		Path        string     `yaml:"path"`
		Host        string     `yaml:"host"`
		Port        int16      `yaml:"port"`
		AuthSource  string     `yaml:"authSource"`
		AuthType    string     `yaml:"authType"`
		Description UserConfig `yaml:"description"`
		BaseName    string     `yaml:"basename"`
	}
	Login struct {
		User UserConfig
	} `yaml:"login"`
}

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type CPU struct {
	Number    int      `json:"number"`
	ModelName string   `json:"modelName"`
	Cores     int32    `json:"cores"`
	Mhz       float64  `json:"mhz"`
	CacheSize int32    `json:"cacheSize"`
	Percent   float64  `json:"percent"`
	Flags     []string `json:"flags"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total_mb"`
	Available   uint64  `json:"available_mb"`
	Used        uint64  `json:"used_mb"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free_mb"`
	Cached      uint64  `json:"cached_mb"`
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Name        string  `json:"name"`
	Mountpoint  string  `json:"mountpoint"`
	Type        string  `json:"type"`
}

type NetworkInfo struct {
	Name    string   `json:"name"`
	Address []string `json:"address"`
}

type NodeInfo struct {
	Hostname        string `json:"hostname"`
	OS              string `json:"os"`
	Platform        string `json:"platform"`
	PlatformVersion string `json:"platform_version"`
	KernelVersion   string `json:"kernel_version"`
	Arch            string `json:"arch"`
}
