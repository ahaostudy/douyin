# 字节青训营 —— 抖音极简版

<img src="public/image/img.png" style="max-height: 400px;">

## 目录

- [项目环境](#项目环境)
    - [安装FFmpeg](#安装FFmpeg)
    - [配置MySQL](#配置MySQL)
    - [Redis与RabbitMQ](#Redis与RabbitMQ)
- [项目启动](#项目启动)
    - [clone并进入项目](#clone并进入项目)
    - [编辑配置文件](#编辑配置文件)
    - [同步项目依赖](#同步项目依赖)
    - [启动项目](#启动项目)

## 项目环境

- Golang 1.20
- MySQL
- Redis
- RabbitMQ
- FFmpeg

### 安装FFmpeg

下载并解压 [ffmpeg-master-latest-win64-gpl.zip](https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-win64-gpl.zip)
，将解压后的 `bin` 目录添加到环境变量即可。

### 配置MySQL

本项目运行需要将`sql_mode`中的`ONLY_FULL_GROUP_BY`删掉。

先登录MySQL，查询`sql_mode`

```mysql
SELECT @@global.sql_mode;
```

```text
+---------------------------------------------------------------+
| @@global.sql_mode                                             |
+---------------------------------------------------------------+
| ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION |
+---------------------------------------------------------------+
```

删除结果中的 ONLY_FULL_GROUP_BY

```mysql
SET GLOBAL sql_mode = 'STRICT_TRANS_TABLES,NO_ENGINE_SUBSTITUTION';
```

### Redis与RabbitMQ

自行安装

## 项目启动

### clone并进入项目

```shell
git clone git@gitee.com:loau/douyin.git
cd douyin
```

### 编辑配置文件

将 `config.bak.yaml` 复制为 `config.yaml`，并在 `config.yaml` 中根据需要修改配置。

```shell
cp config/config.bak.yaml config/config.yaml
vim config/config.yaml
```

### 同步项目依赖

```shell
go mod init main
go mod tidy
```

### 启动项目

```shell
go run main.go
```