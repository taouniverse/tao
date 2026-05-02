# Tao Core 测试工具

.PHONY: help test coverage clean

help:
	@echo "Tao Core 测试工具"
	@echo ""
	@echo "可用命令:"
	@echo "  make test        - 运行测试"
	@echo "  make coverage    - 生成测试覆盖率报告"
	@echo "  make clean       - 清理测试生成的文件"

test:
	go test -v ./...

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

clean:
	rm -f coverage.out coverage.html
