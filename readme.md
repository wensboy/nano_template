## Nano Template

一个实验性的项目, 用于个人在 ai 时代下做一些关于个人开发效率提升的实验. 如果你对这个项目感兴趣, 直接 clone 即可.

```md
### Git Commit Message Rules

> format

<type><scope>: <subject>
<body>
<footer>

**type**
| 类型         | 说明                          |
| ---------- | --------------------------- |
| `feat`     | 新功能（feature）                |
| `fix`      | 修复 bug                      |
| `docs`     | 文档（documentation）           |
| `style`    | 格式（不影响代码运行的变动）              |
| `refactor` | 重构（即不是新增功能，也不是修改 bug 的代码变动） |
| `perf`     | 性能优化                        |
| `test`     | 增加测试                        |
| `chore`    | 构建过程或辅助工具的变动                |
| `ci`       | 持续集成相关                      |
| `build`    | 构建系统或外部依赖变动                 |
| `revert`   | 回滚 commit                   |
**scope**
可选，用于说明 commit 影响的范围. 例如: 
auth
user
api
ui
deps
**body**
每行不超过 72 个字符
说明改动原因和与上一版本的差异
**footer**
用于不兼容变动（BREAKING CHANGE:）
或关闭 Issue（Closes #123, #456）

> examples

---
feat(auth): add oauth2 login support

Implement Google and GitHub OAuth2 authentication.
Add user profile synchronization on first login.

Closes #42
---
fix(api): resolve user list pagination bug

Offset calculation was incorrect when page size exceeded
total records. Fixed by adding boundary check.

Fixes #88
---
refactor(user): extract validation logic to service layer

Move input validation from controller to dedicated
UserValidationService for better testability.

BREAKING CHANGE: validateUser() signature changed
---
```