# 青训营项目——抖音极简版

## 项目环境
- MySQL

## 项目启动

#### 克隆并进去项目
```shell
git clone git@gitee.com:loau/douyin.git
cd douyin
```

#### 修改配置文件
将 `config.bak.yaml` 复制为 `config.yaml`，并在 `config.yaml` 中根据需要修改配置。
```shell
cp config/config.yaml config.yaml
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