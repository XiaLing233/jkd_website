# 计科导大作业文件结构说明

## vue 前端

我需要完成的绝大部分内容存放在了 /src 下

其中，/src/router/ 存放了路由配置文件；/src/views/ 存放了 vue 的三段式文件；/src/assets/ 存放了图片资源。其余文件夹暂时没有投入使用。

## flask 后端

它是一个独立文件夹，在 /src/ics-flask/ 下。其中的 `app.py` 是编写的配置文件。

注：`app.py` 中数据库的用户名和密码已略去，而且其实服务器上使用的用户名和密码也不是这个。

## 环境配置

```bash
# 后端的配置
pip install flask flask-cors flask-Limiter requests
pip install configParser logging
pip install redis # 需要进行硬盘上的存储
pip install gunicorn # wsgi
```

这次的更新，主要是进行了限流政策，防止同一时间大量请求汇入。一会改成每分钟十次吧。走的弯路就是，flask 可能是单线程，用内存来监控请求数量就 OK 了，然而 gunicorn 可能是多线程，所以导致明明触发了阈值，但是并不会导致限速。以及请求是前端发送给 flask 的，所以获取到的 IP 只是 127.0.0.1。以及增加了日志。
