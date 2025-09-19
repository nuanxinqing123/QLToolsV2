package autoload

type App struct {
	Mode    string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Address string `mapstructure:"address" json:"address" yaml:"address"`
	Port    int    `mapstructure:"port" json:"port" yaml:"port"`
}
