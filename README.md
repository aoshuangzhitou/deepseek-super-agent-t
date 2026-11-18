# deepseek-super-agent-t

用于试用 DeepSeek Harness Superagent 的测试仓库。

## hello web 服务

`cmd/hello-web` 是一个只输出 `hello` 的最小 web 服务：访问根路径返回 `hello`，
其他路径返回 404，非 GET/HEAD 方法返回 405。

根路径按 `Accept` 头做内容协商：浏览器（`Accept` 含 `text/html`）会拿到一张
自带样式的单页卡片，`curl` / 脚本等发送 `*/*` 的客户端仍然只收到纯文本
`hello`，`Content-Type` 分别是 `text/html; charset=utf-8` 和
`text/plain; charset=utf-8`。

### 本地启动

```sh
./scripts/run-hello-web.sh            # 默认监听 127.0.0.1:8123
./scripts/run-hello-web.sh 127.0.0.1:9000
```

脚本会编译到 `.run/hello-web`、在后台启动服务，并把 PID 写入
`.run/hello-web.pid`、日志写入 `.run/hello-web.log`。

访问地址：<http://127.0.0.1:8123/>

停止服务：

```sh
kill "$(cat .run/hello-web.pid)"
```

也可以直接运行：

```sh
go run ./cmd/hello-web -addr 127.0.0.1:8123
```

### 测试

```sh
go test ./...
```
