# kinship — Go 语言家族图谱与亲缘关系推算 HTTP 后端服务，支持 GEDCOM 解析、祖先/后代查询和亲属路径搜索

本 HTTP 后端服务根据 GEDCOM 家族图谱推算亲缘：给定人物标识与查询类型（祖先/后代/路径），返回关系路径；环路或缺人必须报错，不得静默丢代。

## 构建 / 运行 / 测试

```text
go build ./...                                    # 编译
go run . -addr :8080                              # 启动 HTTP 服务（/api/parse, /api/ancestors, /api/kin）
go test ./...                                     # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
