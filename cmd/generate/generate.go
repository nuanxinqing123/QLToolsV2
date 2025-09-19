package main

import (
	"flag"
	"fmt"

	"github.com/nuanxinqing123/QLToolsV2/internal/app/config"
	"github.com/nuanxinqing123/QLToolsV2/internal/app/initializer"
	"gorm.io/gen"
)

func main() {
	// 解析命令行参数，支持通过 -config 或 -c 指定配置文件路径
	// 示例：./app -config=/path/to/config.yaml 或 ./app -c ./configs/dev.yaml
	var configPath string
	flag.StringVar(&configPath, "config", "", "配置文件路径（默认：configs/config.yaml）")
	flag.StringVar(&configPath, "c", "", "配置文件路径（默认：configs/config.yaml）")
	flag.Parse()

	// 初始化配置文件
	config.VP = initializer.Viper(configPath)
	if config.VP != nil {
		fmt.Println("Viper初始化成功")
	}

	// 初始化日志
	config.Log = initializer.Zap()

	// 初始化数据库
	config.DB = initializer.Gorm()
	if config.DB == nil {
		config.Log.Fatal("数据库连接失败")
		return
	} else {
		// 初始化表
		initializer.RegisterTables(config.DB)
		config.Log.Info("数据库连接成功")
	}

	// 创建生成器实例
	g := gen.NewGenerator(gen.Config{
		// 输出路径
		OutPath: "./internal/repository",
		// 输出模式
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		// 表字段可为空值时，对应结构体字段使用指针类型
		FieldNullable: true,
		// 生成字段类型标签
		FieldWithTypeTag: true,
		// 生成字段DB标签
		FieldWithIndexTag: true,
	})

	// 使用数据库
	g.UseDB(config.DB)

	// 生成所有表的模型和查询代码
	// 也可以用 g.GenerateModel("users") 指定表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行代码生成
	g.Execute()
}
