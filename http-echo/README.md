# http-echo

## 程序启动参数

| 启动参数 | 描述 |
| :------- | :--- |
| host     | 监听地址 |
| port     | 监听端口 |

## 请求参数

### content

content 参数不同值的说明如下：

| 值    | 描述                                                               | 示例路径                   |
| :----- | :----------------------------------------------------------------- | :------------------------- |
| random | 响应体返回随机字符，大小由参数 length 决定                          | /?content=random&length=100 |
| echo   | 原样返回参数 data 的值                                             | /?content=echo&data=abc     |

### speed

返回包速度，控制每秒可以返回的字节数。

示例：/?content=random&length=1000&speed=100

服务器将返回 10 次，耗时 10 秒。

## Docker 镜像

```shell
docker pull fancymore/http-echo:latest
docker run -p 80:80 fancymore/http-echo /http-echo --port 80
```