# 青训营项目——抖音极简版

## 项目环境

- MySQL
- Redis
- RabbitMQ
- [FFmpeg](https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip)

## 项目启动

#### 安装环境

MySQL、Redis、RabbitMQ：略

FFmpeg：下载并解压 [ffmpeg-master-latest-win64-gpl.zip](https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip)
，将解压后的 `bin` 目录添加到环境变量即可。

#### 克隆并进入项目

```shell
git clone git@gitee.com:loau/douyin.git
cd douyin
```

#### 修改配置文件

将 `config.bak.yaml` 复制为 `config.yaml`，并在 `config.yaml` 中根据需要修改配置。

```shell
cp config/config.bak.yaml config.yaml
```

#### 同步项目依赖

```shell
go mod init main
go mod tidy
```

#### 启动项目

```shell
go run main.go
```