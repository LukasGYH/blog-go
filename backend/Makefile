.PHONY: run build tidy test reset-db

# 启动后端服务
run:
	go run .

# 编译到 bin/blog-api
build:
	go build -o bin/blog-api .

# 拉取 / 整理依赖
tidy:
	go mod tidy

# 运行单元测试
test:
	go test ./...

# 删除 SQLite 数据库文件，下次启动自动重建并灌入种子数据
reset-db:
	rm -f blog.db data/blog.db
