<!--suppress HtmlDeprecatedAttribute -->
<p align="center">
  <a href="https://github.com/whyour/qinglong">
    <img width="150" src="https://z3.ax1x.com/2021/11/18/I7MpAe.png" alt="Img">
  </a>
</p>

<h1 align="center">青龙Tools V2</h1>

## 🍭功能介绍
- 开发版本：2.17.2 _[理论不支持v2.11.3以前的版本]_
- 多面板&多变量设置
- 轮询分布上传数据
- 变量支持限量(限额)上传
- [开发中] 支持上传插件 & 定时插件
- [开发中] 支持插件在线编辑 & 保存
- 其余内容正在规划中...

## 🍳技术栈
- 语言: Golang
- 框架: Gin框架
- 数据库: SQLite (使用非CGO的依赖库 _[ 相比于CGO版本性能要低一丢丢, 主要是开发方便 ]_ )

## 🐳Docker运行
```shell
# 运行镜像
docker run --restart=always -itd --name QLToolsV2 -v $PWD/config:/QLToolsV2/config -p 1500:1500 nuanxinqing123/qltools_v2:latest

# docker compose 运行
docker-compose up -d
```


## Dockerfile构建运行
```shell
# 编译镜像
docker build -t qltools_v2:latest .

# 运行镜像
docker run -d -p 1500:1500 qltools_v2:latest
```

## 🎯开发计划

- 开发文档: [ApiFox ->](https://apifox.com/apidoc/shared-85e6a037-263e-4a8f-8e96-3c56a924c242)

【联系方式】

TG：[@Nuanxinqing](https://t.me/Nuanxinqing)

【交流群】

青龙Tools频道：[@qltool](https://t.me/qltool)

青龙Tools吹水群：[@qltools](https://t.me/qltools)

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=nuanxinqing123/QLToolsV2&type=Date)](https://star-history.com/#nuanxinqing123/QLToolsV2&Date)