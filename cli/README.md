# ticket-cli

石河子大学 AI 迎新助手工单系统命令行工具。让本地 Agent / 运维 / 业务人员能高效访问工单系统，支持列表查询、详情、提交、答复、归档、删除、批量操作、导出等所有 v3 后端能力。

## 安装

```bash
# 方式 1: go install
cd cli
go install .

# 方式 2: 下载预编译
# (暂无 release，请自行编译)
cd cli && go build -o ticket.exe .      # Windows
cd cli && go build -o ticket .          # Linux/macOS
```

## 配置

首次运行自动创建 `~/.ticket-cli.yaml`，或手动复制 `examples/config.yaml`：

```yaml
api_base: http://120.48.98.164:8080/api/v1
api_key: sk-ticket-prod-shzu-2026       # 用于 submit（智能体模式）
username: admin                          # 用于 login
password: admin123                       # 用于 login
output_format: table                     # table | json | pretty
```

命令行参数优先级最高：
```bash
ticket --api-base http://localhost:8000/api/v1 --api-key sk-xxx list
```

## 命令清单

### 登录
```bash
ticket login                              # 用配置文件中的用户名密码
ticket login -u admin -p admin123         # 命令行覆盖
```

### 列表查询
```bash
ticket list                               # 默认第一页 20 条
ticket list --size 50                     # 50 条/页
ticket list --status pending              # 按状态筛选
ticket list --category 宿舍                # 按分类筛选
ticket list --archived true               # 只看已归档
ticket list --q "宿舍" --size 10           # 关键字搜索
ticket list --start 2026-07-20 --end 2026-07-23   # 时间范围
ticket list --all                         # 翻所有页
ticket list --json                        # JSON 输出
```

### 详情
```bash
ticket get T20260723-ABCDEF
ticket get T20260723-ABCDEF --json
```

### 提交工单（智能体模式，需 API Key）
```bash
ticket submit --question "宿舍分配标准"
ticket submit --question "..." --user-id sess_001 --rag "检索片段" --source hiagent_chat
```

### 答复 / 关闭 / 重开
```bash
ticket answer T20260723-ABCDEF --text "2026 年 8 月 29 日报到"
ticket answer T20260723-ABCDEF --file answer.md   # 从文件读
ticket close T20260723-ABCDEF --reason "已处理"
ticket reopen T20260723-ABCDEF
```

### 归档 / 取消归档
```bash
ticket archive T20260723-ABCDEF
ticket unarchive T20260723-ABCDEF
```

### 删除
```bash
ticket delete T20260723-ABCDEF            # 软删
ticket delete T20260723-ABCDEF --hard     # 硬删（不可恢复）
ticket batch-delete T1 T2 T3              # 批量删除
```

### 导出
```bash
ticket export --format csv --output tickets.csv
ticket export --format csv --status closed --output closed.csv
ticket export --format json --output all.json
ticket export --format json --archived true --output archived.json
```

## 典型工作流

### 场景 1：批量关闭 7 天前的已答复工单
```bash
# 1. 列出候选工单
ticket list --status answered --start 2026-07-15 --end 2026-07-22 --json > pending.json

# 2. 提取 ID 列表
ticket list --status answered --start 2026-07-15 --end 2026-07-22 --all --json \
  | jq -r '.[].ticket_id' > ids.txt

# 3. 批量关闭（需要先实现批量 close，当前 CLI 只有 batch-delete）
# 也可以在管理后台多选批量操作
```

### 场景 2：按关键字搜索并导出 CSV
```bash
ticket list --q "奖学金" --all --json | jq '.' > search.json
# 或者
ticket export --format csv --output scholarships.csv
# 然后在管理后台 UI 搜索"奖学金"后点导出
```

### 场景 3：Agent 自动化——批量处理同一类问题
```bash
# 1. 列出所有"宿舍"相关工单
ticket list --q "宿舍" --all --json | jq -r '.[].ticket_id' | while read id; do
  ticket answer "$id" --text "石河子大学宿舍标准请见迎新网"
done
```

## 鉴权

- **API Key**（`api_key`）—— 智能体调用 `submit` 用
- **JWT**（`token`，登录后自动获取）—— 管理员所有其他操作

CLI 自动选择：
- 有 token → 用 JWT
- 无 token + 有 api_key → 用 API Key（仅适用于双轨鉴权的端点）

## 跨平台

| 平台 | 二进制 |
|------|-------|
| Windows | `ticket.exe` |
| Linux | `ticket` |
| macOS | `ticket` |

编译：
```bash
GOOS=windows GOARCH=amd64 go build -o ticket.exe .    # Windows
GOOS=linux GOARCH=amd64 go build -o ticket .          # Linux
GOOS=darwin GOARCH=arm64 go build -o ticket-mac .     # macOS Apple Silicon
```

## 退出码

- `0` — 成功
- `1` — 失败（HTTP 错误 / 网络错误 / 参数错误）

## 常见问题

**Q: 报 `unauthorized`？**
A: 先 `ticket login`，或者在配置里填 `api_key`。

**Q: token 过期了？**
A: 重新 `ticket login` 即可。

**Q: 怎么列出所有工单？**
A: `ticket list --all`。
