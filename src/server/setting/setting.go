package setting

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Setting struct {
	Server ServerSetting
	Bt     BtServerSetting
	Log    LogSetting
	Redis  RedisSetting
	Mysql  MysqlSetting
}

type ServerSetting struct {
	Domain    string `yaml:"domain"`
	BoundIp   string `yaml:"boundIp"`
	Port      int
	RestPort  int    `yaml:"restPort"`
	WebPort   int    `yaml:"webPort"`
	CrtFile   string `yaml:"crtFile"`
	KeyFile   string `yaml:"keyFile"`
	MediaPath string `yaml:"mediaPath"`

	HlsPath    string
	PosterPath string
}

type BtServerSetting struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	SavePath string `yaml:"savePath"`
}

type LogSetting struct {
	FileName string `yaml:"fileName"`
	Level    string `yaml:"level"`
}

type RedisSetting struct {
	Ip       string `yaml:"ip"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type MysqlSetting struct {
	Ip           string
	Port         int
	User         string
	Password     string
	Dbname       string
	MaxOpenConns int `yaml:"maxOpenConns"`
	MaxIdleConns int `yaml:"maxIdleConns"`
}

var GS *Setting = new(Setting)

func Init() {
	file_name := "local_setting.yml"
	config_path := "."
	_, err := os.Stat(config_path + "/" + file_name)
	if err != nil {
		file_name = "setting.yml"
	}
	yamlFile, err := ioutil.ReadFile(config_path + "/" + file_name)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = yaml.Unmarshal(yamlFile, GS)
	if err != nil {
		fmt.Println(err.Error())
	}

	GS.Server.MediaPath = path.Clean(GS.Server.MediaPath)
	GS.Bt.SavePath = path.Clean(GS.Bt.SavePath)

	GS.Server.HlsPath = GS.Server.MediaPath + "/hls"
	GS.Server.PosterPath = GS.Server.MediaPath + "/poster"

	os.MkdirAll(GS.Server.HlsPath, 0755)
	os.MkdirAll(GS.Server.PosterPath, 0755)
	os.MkdirAll(GS.Bt.SavePath, 0755)
}

func GetMysqlConnectStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		GS.Mysql.User,
		GS.Mysql.Password,
		GS.Mysql.Ip,
		GS.Mysql.Port,
		GS.Mysql.Dbname)
}
