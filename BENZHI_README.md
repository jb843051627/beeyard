# Beeyard 打包说明

## 构建

```bash
./build_benzhi_docker.sh beeyard-1 linux/amd64
./build_benzhi_docker.sh beeyard-1 linux/arm64
```

镜像名：`benzhi/<name>:latest`

## 运行

```bash
docker run --rm -p 8080:8080 benzhi/beeyard:latest
```

## 验证

```bash
docker run --rm benzhi/beeyard:latest go version
docker run --rm benzhi/beeyard:latest go build ./...
docker run --rm benzhi/beeyard:latest go test ./...
```
