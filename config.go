package common

import "github.com/spf13/viper"

type (
	ServiceConfig interface {
		ConfigFromFileName(config string) ServiceConfig
	}
	ServiceConfigImpl struct {
		ServiceName   string
		LogFile       string
		JaegerAddr    string
		ConsulAddr    string
		ConsulPort    int
		HttpAddr      string
		DbUri         string
		Port          int
		RedisHost     string
		RedisPort     string
		RedisPassword string
		RedisDB       int
		MqUri         string
	}
)

func (serviceConfig ServiceConfigImpl) ConfigFromFileName(config string) ServiceConfig {
	viper.SetConfigFile(config)
	if err := viper.ReadInConfig(); err != nil {
		return serviceConfig
	}
	serviceConfig = ServiceConfigImpl{
		ServiceName:   viper.GetString("service_name"),
		LogFile:       viper.GetString("log_file"),
		JaegerAddr:    viper.GetString("jaeger_addr"),
		ConsulAddr:    viper.GetString("consul_addr"),
		ConsulPort:    viper.GetInt("consul_port"),
		HttpAddr:      viper.GetString("address"),
		DbUri:         viper.GetString("db_uri"),
		Port:          viper.GetInt("port"),
		RedisHost:     viper.GetString("redis_host"),
		RedisPort:     viper.GetString("redis_port"),
		RedisPassword: viper.GetString("redis_password"),
		RedisDB:       viper.GetInt("redis_db"),
		MqUri:         viper.GetString("mq_uri"),
	}

	return serviceConfig
}

//
//func (serviceConfig ServiceConfigImpl) GetDB() (db *gorm.DB, err error) {
//	db, err = gorm.Open("mysql", serviceConfig.DbUri)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	db.SingularTable(true)
//	db.LogMode(true)
//	db.DB().SetMaxIdleConns(10)
//	db.DB().SetMaxOpenConns(100)
//	db.DB().SetConnMaxLifetime(time.Hour)
//	return
//}
//
//func (conf ServiceConfigImpl) Migrate(db *gorm.DB, models []BaseModel) (err error) {
//	for _, value := range models {
//		db.AutoMigrate(value)
//	}
//	return nil
//}
