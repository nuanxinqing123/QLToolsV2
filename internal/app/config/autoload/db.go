package autoload

type DB struct {
	Type         string `mapstructure:"type" json:"type" yaml:"type"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Port         int    `mapstructure:"port" json:"port" yaml:"port"`
	Name         string `mapstructure:"name" json:"name" yaml:"name"`
	UserName     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`
	LogLevel     string `mapstructure:"log-level" json:"log-level" yaml:"log-level"`
}
