package setting

import (
	"fmt"
	"os"
	"path"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Setting struct {
	Server ServerSetting
	Bt     BtServerSetting
	Log    LogSetting
	Redis  RedisSetting
	Mysql  MysqlSetting
}

type ServerSetting struct {
	Domain        string `yaml:"domain"`
	BoundIp       string `yaml:"boundIp"`
	Port          int
	WebPort       int    `yaml:"webPort"`
	CrtFile       string `yaml:"crtFile"`
	KeyFile       string `yaml:"keyFile"`
	MediaPath     string `yaml:"mediaPath"`
	EnableCrawler bool   `yaml:"enableCrawler"`

	HlsPath    string
	PosterPath string
}

type BtServerSetting struct {
	BtClientAddress   string `yaml:"btClientAddress"`
	ListenInterfaces  string `yaml:"listenInterfaces"`
	SavePath          string `yaml:"savePath"`
	ProxyHost         string `yaml:"proxyHost"`
	ProxyPort         int32  `yaml:"proxyPort"`
	ProxyType         string `yaml:"proxyType"`
	UploadRateLimit   int32  `yaml:"uploadRateLimit"`
	DownloadRateLimit int32  `yaml:"downloadRateLimit"`
	HashingThreads    int32  `yaml:"hashingThreads"`
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

var setting atomic.Pointer[Setting]

type OnConfigChange func()

var onFunsMtx sync.Mutex
var onCfgChangeFuns map[string]OnConfigChange

func AddOnCfgChangeFun(name string, f OnConfigChange) {
	if f == nil {
		return
	}
	onFunsMtx.Lock()
	onCfgChangeFuns[name] = f
	onFunsMtx.Unlock()
}

func DelOnCfgChangeFun(name string) {
	onFunsMtx.Lock()
	delete(onCfgChangeFuns, name)
	onFunsMtx.Unlock()
}

func GS() *Setting {
	return setting.Load()
}

func Init(config_file_full_path string) {

	onFunsMtx.Lock()
	onCfgChangeFuns = make(map[string]OnConfigChange)
	onFunsMtx.Unlock()

	if _, err := os.Stat(config_file_full_path); err != nil {
		config_file_full_path = "./server.yml"
		if _, err := os.Stat(config_file_full_path); err != nil {
			config_file_full_path = "/etc/pnas/server.yml"
		}
	}

	yamlFile, err := os.ReadFile(config_file_full_path)
	if err != nil {
		fmt.Println(err.Error())
	}

	var s Setting
	err = yaml.Unmarshal(yamlFile, &s)
	if err != nil {
		fmt.Println(err.Error())
	}
	initPath(&s)
	setting.Store(&s)
	watcher, err := fsnotify.NewWatcher()
	if err == nil {
		watcher.Add(path.Dir(config_file_full_path))
		configFileName := path.Base(config_file_full_path)
		go func() {
			defer watcher.Close()
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if (event.Name == configFileName || event.Name == "./"+configFileName) &&
						(event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) {

						yamlFile, err := os.ReadFile(config_file_full_path)
						if err != nil {
							fmt.Println(err.Error())
						}

						var s Setting
						err = yaml.Unmarshal(yamlFile, &s)
						if err != nil {
							continue
						}
						initPath(&s)
						setting.Store(&s)
						funs := []OnConfigChange{}
						onFunsMtx.Lock()
						for _, f := range onCfgChangeFuns {
							funs = append(funs, f)
						}
						onFunsMtx.Unlock()
						for _, f := range funs {
							f()
						}
					}
				case _, ok := <-watcher.Errors:
					if !ok {
						return
					}
				}
			}
		}()
	}
}

func initPath(s *Setting) {
	s.Server.MediaPath = path.Clean(s.Server.MediaPath)
	s.Bt.SavePath = path.Clean(s.Bt.SavePath)
	s.Server.HlsPath = s.Server.MediaPath + "/hls"
	s.Server.PosterPath = s.Server.MediaPath + "/poster"

	os.MkdirAll(s.Server.HlsPath, 0755)
	os.MkdirAll(s.Server.PosterPath, 0755)
	os.MkdirAll(s.Bt.SavePath, 0755)
}

func GetMysqlConnectStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		GS().Mysql.User,
		GS().Mysql.Password,
		GS().Mysql.Ip,
		GS().Mysql.Port,
		GS().Mysql.Dbname)
}
