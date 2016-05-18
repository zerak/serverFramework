// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// core config base on beego config
package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TaXingTianJi/serverFramework/utils"
	"github.com/astaxie/beego/config"
)

// Config is the main struct for serverConfig
type Config struct {
	AppName             string //Application name
	RunMode             string //Running Mode: dev | prod
	RouterCaseSensitive bool
	ServerName          string
	RecoverPanic        bool
	MaxMemory           int64
	EnableErrorsShow    bool
	TCPAddr             string
	TCPPort             int
	MsgSize             int // client msg buffer size

	AdminConf AdminConfig // monitor config
	LogConf   LogConfig   // log config
	DBConf    DBConfig    // db config
}

// for debug print
func (c Config) String() {
	fmt.Printf("config[%v %v %v %v %v %v %v %v] admin[%v %v %v %v %v %v %v %v] log[%v %v %v]\n\n", c.AppName, c.RunMode, c.RouterCaseSensitive, c.ServerName, c.RecoverPanic, c.MaxMemory, c.EnableErrorsShow, c.TCPAddr,
		c.AdminConf.ServerTimeOut, c.AdminConf.ListenTCP4, c.AdminConf.EnableHTTP, c.AdminConf.HTTPAddr, c.AdminConf.HTTPPort, c.AdminConf.EnableAdmin, c.AdminConf.AdminAddr, c.AdminConf.AdminPort,
		c.LogConf.AccessLogs, c.LogConf.FileLineNum, c.LogConf.Outputs)
}

// AdminConfig holds for admin control http related config
type AdminConfig struct {
	ServerTimeOut int64
	ListenTCP4    bool
	EnableHTTP    bool
	HTTPAddr      string
	HTTPPort      int
	EnableAdmin   bool
	AdminAddr     string
	AdminPort     int
}

// LogConfig holds Log related config
type LogConfig struct {
	AccessLogs  bool              // admin check access log
	FileLineNum bool              // Set line num
	Outputs     map[string]string // Store Adaptor : config
}

type DBConfig struct {
	User string
	PW   string
	Addr string
	Port int
	DB   string
}

var (
	// SConfig is the default config for Application
	SConfig *Config

	// AppConfig is the instance of Config, store the config information from file
	appConfig *serverConfig

	// AppPath is the absolute path to the app
	appPath string

	// appConfigPath is the path to the config files
	appConfigPath string

	// appConfigProvider is the provider for the config, default is ini
	appConfigProvider = "ini"
)

func init() {
	appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(appPath)

	SConfig = &Config{
		AppName:             "server",
		RunMode:             DEV,
		RouterCaseSensitive: true,
		ServerName:          "serverFramework:" + VERSION,
		RecoverPanic:        true,
		MaxMemory:           1 << 26, //64MB
		EnableErrorsShow:    true,
		TCPAddr:             "127.0.0.1",
		TCPPort:             60060,
		MsgSize:             10000,
		AdminConf: AdminConfig{
			ServerTimeOut: 0,
			ListenTCP4:    false,
			EnableHTTP:    true,
			HTTPAddr:      "",
			HTTPPort:      8080,
			EnableAdmin:   true,
			AdminAddr:     "",
			AdminPort:     8088,
		},
		LogConf: LogConfig{
			AccessLogs:  false,
			FileLineNum: true,
			Outputs:     map[string]string{"console": ""},
		},
		DBConf: DBConfig{
			User: "user",
			PW:   "pw",
			Addr: "localhost",
			Port: 3306,
			DB:   "testDb",
		},
	}

	appConfigPath = filepath.Join(appPath, "conf", "app.conf")
	if !utils.FileExists(appConfigPath) {
		appConfig = &serverConfig{innerConfig: config.NewFakeConfig()}
		return
	}

	if err := parseConfig(appConfigPath); err != nil {
		panic(err)
	}
}

// now only support ini, next will support json.
func parseConfig(appConfigPath string) (err error) {
	appConfig, err = newAppConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return err
	}
	// set the run mode first
	if envRunMode := os.Getenv("SERVER_RUNMODE"); envRunMode != "" {
		SConfig.RunMode = envRunMode
	} else if runMode := appConfig.String("RunMode"); runMode != "" {
		SConfig.RunMode = runMode
	}

	SConfig.AppName = appConfig.DefaultString("AppName", SConfig.AppName)
	SConfig.RecoverPanic = appConfig.DefaultBool("RecoverPanic", SConfig.RecoverPanic)
	SConfig.RouterCaseSensitive = appConfig.DefaultBool("RouterCaseSensitive", SConfig.RouterCaseSensitive)
	SConfig.ServerName = appConfig.DefaultString("ServerName", SConfig.ServerName)
	SConfig.MaxMemory = appConfig.DefaultInt64("MaxMemory", SConfig.MaxMemory)
	SConfig.EnableErrorsShow = appConfig.DefaultBool("EnableErrorsShow", SConfig.EnableErrorsShow)
	SConfig.TCPAddr = appConfig.DefaultString("TCPAddr", SConfig.TCPAddr)
	SConfig.TCPPort = appConfig.DefaultInt("TCPPort", SConfig.TCPPort)
	SConfig.MsgSize = appConfig.DefaultInt("MsgSize", SConfig.MsgSize)

	SConfig.AdminConf.HTTPAddr = appConfig.String("HTTPAddr")
	SConfig.AdminConf.HTTPPort = appConfig.DefaultInt("HTTPPort", SConfig.AdminConf.HTTPPort)
	SConfig.AdminConf.ListenTCP4 = appConfig.DefaultBool("ListenTCP4", SConfig.AdminConf.ListenTCP4)
	SConfig.AdminConf.EnableHTTP = appConfig.DefaultBool("EnableHTTP", SConfig.AdminConf.EnableHTTP)
	SConfig.AdminConf.EnableAdmin = appConfig.DefaultBool("EnableAdmin", SConfig.AdminConf.EnableAdmin)
	SConfig.AdminConf.AdminAddr = appConfig.DefaultString("AdminAddr", SConfig.AdminConf.AdminAddr)
	SConfig.AdminConf.AdminPort = appConfig.DefaultInt("AdminPort", SConfig.AdminConf.AdminPort)
	SConfig.AdminConf.ServerTimeOut = appConfig.DefaultInt64("ServerTimeOut", SConfig.AdminConf.ServerTimeOut)

	SConfig.LogConf.AccessLogs = appConfig.DefaultBool("LogAccessLogs", SConfig.LogConf.AccessLogs)
	SConfig.LogConf.FileLineNum = appConfig.DefaultBool("LogFileLineNum", SConfig.LogConf.FileLineNum)

	SConfig.DBConf.User = appConfig.DefaultString("DBUser", SConfig.DBConf.User)
	SConfig.DBConf.PW = appConfig.DefaultString("DBPW", SConfig.DBConf.PW)
	SConfig.DBConf.Addr = appConfig.DefaultString("DBAddr", SConfig.DBConf.Addr)
	SConfig.DBConf.Port = appConfig.DefaultInt("DBPort", SConfig.DBConf.Port)
	SConfig.DBConf.DB = appConfig.DefaultString("DBName", SConfig.DBConf.DB)

	if lo := appConfig.String("LogOutputs"); lo != "" {
		los := strings.Split(lo, ";")
		for _, v := range los {
			if logType2Config := strings.SplitN(v, ",", 2); len(logType2Config) == 2 {
				SConfig.LogConf.Outputs[logType2Config[0]] = logType2Config[1]
			} else {
				continue
			}
		}
	}

	//init log
	ServerLogger.Reset()
	for adaptor, config := range SConfig.LogConf.Outputs {
		err = ServerLogger.SetLogger(adaptor, config)
		if err != nil {
			fmt.Printf("%s with the config `%s` got err:%s\n", adaptor, config, err)
		}
	}
	if SConfig.RunMode == DEV {
		SetLevel(LevelInformational)
	} else if SConfig.RunMode == PROD {
		SetLevel(LevelWarning)
	}
	SetLogFuncCall(SConfig.LogConf.FileLineNum)

	//fmt.Print(SConfig)
	ServerLogger.Info("%v", SConfig)
	return nil
}

// LoadAppConfig allow developer to apply a config file
func LoadAppConfig(adapterName, configPath string) error {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return err
	}

	if !utils.FileExists(absConfigPath) {
		return fmt.Errorf("the target config file: %s don't exist", configPath)
	}

	if absConfigPath == appConfigPath {
		return nil
	}

	appConfigPath = absConfigPath
	appConfigProvider = adapterName

	return parseConfig(appConfigPath)
}

type serverConfig struct {
	innerConfig config.Configer
}

func newAppConfig(appConfigProvider, appConfigPath string) (*serverConfig, error) {
	ac, err := config.NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &serverConfig{ac}, nil
}

func (b *serverConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(SConfig.RunMode+"::"+key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *serverConfig) String(key string) string {
	if v := b.innerConfig.String(SConfig.RunMode + "::" + key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *serverConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(SConfig.RunMode + "::" + key); v[0] != "" {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *serverConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(SConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *serverConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(SConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *serverConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(SConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *serverConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(SConfig.RunMode + "::" + key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *serverConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serverConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *serverConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *serverConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}
