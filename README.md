## PostSystem

PostSystem 是一个使用 **Go + Gin + GORM** 实现的轻量级新闻发布系统，支持用户注册、登录、发布/编辑/删除新闻以及浏览新闻列表。前端采用原生 HTML + jQuery，并通过统一的 `main.css` 提供现代化 UI。

### 功能特性
- 用户注册 / 登录 / 退出，密码在前端做 MD5 哈希
- 基于 JWT 的免登录访问控制，中间件自动解析 Cookie
- 新闻的发布、列表展示、详情查看、更新与删除
- 个人密码修改、权限校验（仅作者可编辑/删除自己的帖子）

### 本地运行
1. **安装依赖**
   ```bash
   go mod tidy
   ```
2. **配置数据库 / JWT**
   - 数据库配置：`conf/db.yaml`
   - JWT 配置：`conf/jwt.yaml`
3. **启动服务**
   ```bash
   go run main.go
   ```
   服务器默认监听 `http://localhost:8080`
4. **访问页面**
   - 新闻列表：首页：`http://localhost:8080/news`
   - 登录页：`http://localhost:8080/login`
   - 注册页：`http://localhost:8080/register`

### 项目结构
```
.
├── conf/               # 数据库、JWT 配置
├── database/           # GORM 数据访问层
├── handler/            # Gin Handler（接口层）
├── util/               # 日志、配置、JWT 等工具
├── views/
│   ├── html/           # 页面模板
│   ├── css/main.css    # 前端样式
│   └── session_login.html
└── main.go             # 应用入口
```

欢迎根据业务需求扩展功能，例如评论、标签、富文本等。若在使用过程中遇到问题，可通过 README 中的信息快速对照定位。

