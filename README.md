## How to run (with docker)

```
# redis 컨테이너가 실행되고 있지 않다면
$ docker run --name my-redis -p 6379:6379 -d redis

$ git clone https://github.com/meopin-top/convey-your-mind-WSS.git meopin-wss
$ cd meopin-wss
$ docker build -t meopin-wss .
$ docker run meopin-wss
```
