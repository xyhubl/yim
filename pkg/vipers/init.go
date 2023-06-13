package vipers

import (
	"errors"
	"flag"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
)

var (
	ErrConfTyp = errors.New("vipers: only a pointer to struct or map can be unmarshalled from config content")
)

type ViperConf struct {
	conf interface{}
	opts *options
}

func InitViperConf(conf interface{}, opts ...Option) {
	typ := reflect.TypeOf(conf)
	if typ.Kind() != reflect.Ptr {
		panic(ErrConfTyp)
	}
	v := &ViperConf{
		opts: new(options),
	}
	for _, o := range opts {
		o.apply(v.opts)
	}
	v.conf = conf
	v.initModule()
}

func (v *ViperConf) initModule() {
	path := flag.String("config", "", "please input config file like ./config/dev")
	flag.Parse()
	if path == nil || *path == "" {
		flag.Usage()
		os.Exit(1)
	}
	log.Println("*-----------------------let's go-------------------------------*")
	log.Printf("[INFO] config=%s\n", *path)
	log.Printf("[INFO] %s\n", "starting load resoursec.")
	if err := v.initViperConf(*path); err != nil {
		log.Printf("[ERROR] %s\n", err.Error())
		os.Exit(1)
	}
}

func (v *ViperConf) initViperConf(path string) error {
	var (
		err  error
		dir  string
		file *os.File
	)

	file, err = os.Open(path)
	if err != nil {
		return err
	}
	//获取项目的执行路径
	dir, err = os.Getwd()
	if err != nil {
		return err
	}

	vip := viper.New()
	vip.AddConfigPath(dir)
	vip.SetConfigName(file.Name())
	vip.SetConfigType(v.opts.configType)
	if err = vip.ReadInConfig(); err != nil {
		return err
	}
	if v.opts.openWatching {
		vip.WatchConfig()
		vip.OnConfigChange(func(in fsnotify.Event) {
			log.Println("[INFO]*------------------------------------------------------*", file.Name()+" has changed.")
			if err = vip.Unmarshal(v.conf); err != nil {
				return
			}
		})
	}
	if err = vip.Unmarshal(v.conf); err != nil {
		return err
	}
	return nil
}
