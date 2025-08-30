package task_04

type ConfigDB struct {
	Database struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`
}
