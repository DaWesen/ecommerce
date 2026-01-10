package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Log         LogConfig         `mapstructure:"log"`
	Services    ServicesConfig    `mapstructure:"services"`
	JWT         JWTConfig         `mapstructure:"jwt"`
	RateLimit   RateLimitConfig   `mapstructure:"rate_limit"`
	Middleware  MiddlewareConfig  `mapstructure:"middleware"`
	HealthCheck HealthCheckConfig `mapstructure:"health_check"`
}

type ServerConfig struct {
	Port          int           `mapstructure:"port"`
	Env           string        `mapstructure:"env"`
	Name          string        `mapstructure:"name"`
	Timeout       time.Duration `mapstructure:"timeout"`
	EnableSwagger bool          `mapstructure:"enable_swagger"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
	Path   string `mapstructure:"path"`
}

type ServicesConfig struct {
	UserService    ServiceConfig `mapstructure:"user_service"`
	ProductService ServiceConfig `mapstructure:"product_service"`
	OrderService   ServiceConfig `mapstructure:"order_service"`
}

type ServiceConfig struct {
	Name    string `mapstructure:"name"`
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
}

type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
	Issuer      string `mapstructure:"issuer"`
}

type RateLimitConfig struct {
	Enable    bool          `mapstructure:"enable"`
	Requests  int           `mapstructure:"requests"`
	Window    time.Duration `mapstructure:"window"`
	Burst     int           `mapstructure:"burst"`
	IPLimit   bool          `mapstructure:"ip_limit"`
	UserLimit bool          `mapstructure:"user_limit"`
}

type MiddlewareConfig struct {
	EnableAuth    bool `mapstructure:"enable_auth"`
	EnableLog     bool `mapstructure:"enable_log"`
	EnableTrace   bool `mapstructure:"enable_trace"`
	EnableMetrics bool `mapstructure:"enable_metrics"`
}

type HealthCheckConfig struct {
	Enable    bool          `mapstructure:"enable"`
	Interval  time.Duration `mapstructure:"interval"`
	Timeout   time.Duration `mapstructure:"timeout"`
	Threshold int           `mapstructure:"threshold"`
	Path      string        `mapstructure:"path"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 搜索路径
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("../config")
	viper.AddConfigPath("../../config")

	// 设置环境变量前缀
	viper.SetEnvPrefix("GATEWAY")
	viper.AutomaticEnv()

	// 设置默认值
	setDefaults()

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Printf("未找到配置文件，将使用默认值和环境变量")
		} else {
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
	// 服务器配置
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.env", "development")
	viper.SetDefault("server.name", "ecommerce-gateway")
	viper.SetDefault("server.timeout", "30s")
	viper.SetDefault("server.enable_swagger", true)

	// 日志配置
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.path", "./logs")

	// 服务配置
	viper.SetDefault("services.user_service.host", "localhost")
	viper.SetDefault("services.user_service.port", 50052)
	viper.SetDefault("services.user_service.name", "user.service")
	viper.SetDefault("services.user_service.timeout", 5000)

	viper.SetDefault("services.product_service.host", "localhost")
	viper.SetDefault("services.product_service.port", 50051)
	viper.SetDefault("services.product_service.name", "product.service")
	viper.SetDefault("services.product_service.timeout", 5000)

	viper.SetDefault("services.order_service.host", "localhost")
	viper.SetDefault("services.order_service.port", 50053)
	viper.SetDefault("services.order_service.name", "order.service")
	viper.SetDefault("services.order_service.timeout", 10000)

	// JWT 配置
	viper.SetDefault("jwt.secret", "your-secret-key-change-in-production")
	viper.SetDefault("jwt.expire_hours", 24)
	viper.SetDefault("jwt.issuer", "ecommerce-gateway")

	// 限流配置
	viper.SetDefault("rate_limit.enable", true)
	viper.SetDefault("rate_limit.requests", 100)
	viper.SetDefault("rate_limit.window", "1m")
	viper.SetDefault("rate_limit.burst", 20)
	viper.SetDefault("rate_limit.ip_limit", true)
	viper.SetDefault("rate_limit.user_limit", false)

	// 中间件配置
	viper.SetDefault("middleware.enable_auth", true)
	viper.SetDefault("middleware.enable_log", true)
	viper.SetDefault("middleware.enable_trace", true)
	viper.SetDefault("middleware.enable_metrics", true)

	// 健康检查配置
	viper.SetDefault("health_check.enable", true)
	viper.SetDefault("health_check.interval", "30s")
	viper.SetDefault("health_check.timeout", "5s")
	viper.SetDefault("health_check.threshold", 3)
	viper.SetDefault("health_check.path", "/health")
}
