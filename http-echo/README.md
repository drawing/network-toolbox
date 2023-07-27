
# http-echo

Program startup parameters:

| startup parameter | description       |
| :---------------- | :---------------- |
| host              | listening address |
| port              | listening port    |

# request parameters

## content

The description of the different values of the content parameter is as follows:

| value  | description                                                                                     | example path                |
| :----- | :---------------------------------------------------------------------------------------------- | :-------------------------- |
| random | The response body returns random characters, and the size is determined by the parameter length | /?content=random&length=100 |
| echo   | Return the value of the parameter data as it is                                                 | /?content=echo&data=abc     |

## speed

return packet speed, Controls the number of bytes that can be returned per second.

example: /?content=random&length=1000&speed=100

The server will return 10 times, which takes 10 seconds.

## Docker images

```shell
docker pull fancymore/http-echo:latest
docker run -it fancymore/http-echo:latest /http-echo --port 80
```
