package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config应用配置
type Config struct {
	Hertz    HertzConfig    `mapstructure:"hertz"`
	Log      LogConfig      `mapstructure:"log"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Kitex    KitexConfig    `mapstructure:"kitex"`
}

// Hertz配置
type HertzConfig struct {
	Port    int           `mapstructure:"port"`
	Mode    string        `mapstructure:"mode"`
	Timeout time.Duration `mapstructure:"timeout"`
}

// Log配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// Database配置
type DatabaseConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql"`
}

// MySQL配置
type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	User            string `mapstructure:"user"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// Kafka配置
type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Version string   `mapstructure:"version"`
}

// JWT配置
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

// Kitex配置
type KitexConfig struct {
	Port          int `mapstructure:"port"`
	ClientTimeout int `mapstructure:"client_timeout"`
	ServerTimeout int `mapstructure:"server_timeout"`
}

// LoadConfig
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 添加配置搜索路径
	viper.AddConfigPath("./conf")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../conf")
	viper.AddConfigPath("../../conf")
	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("未找到配置文件，使用默认值")
		} else {
			log.Printf("读取配置文件错误: %v", err)
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func setDefaults() {
	// Hertz默认值
	viper.SetDefault("hertz.port", 8081)
	viper.SetDefault("hertz.mode", "debug")
	viper.SetDefault("hertz.timeout", "10s")

	// Log默认值
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "text")
	viper.SetDefault("log.output", "stdout")

	// Database默认值
	viper.SetDefault("database.mysql.host", "localhost")
	viper.SetDefault("database.mysql.port", 3306)
	viper.SetDefault("database.mysql.user", "root")
	viper.SetDefault("database.mysql.password", "")
	viper.SetDefault("database.mysql.dbname", "ecommerce")
	viper.SetDefault("database.mysql.max_open_conns", 100)
	viper.SetDefault("database.mysql.max_idle_conns", 10)
	viper.SetDefault("database.mysql.conn_max_lifetime", 3600)

	// Redis默认值
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)

	// Kafka默认值
	viper.SetDefault("kafka.brokers", []string{"localhost:9092"})
	viper.SetDefault("kafka.version", "2.8.0")

	// JWT默认值
	viper.SetDefault("jwt.secret", "change-this-secret-in-production")
	viper.SetDefault("jwt.expire_hours", 24)

	// Kitex默认值
	viper.SetDefault("kitex.port", 50051)
	viper.SetDefault("kitex.client_timeout", 3000)
	viper.SetDefault("kitex.server_timeout", 5000)
}
