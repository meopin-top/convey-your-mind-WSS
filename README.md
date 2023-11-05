## How to run (with docker)

```
# redis 컨테이너가 실행되고 있지 않다면
$ docker run --name my-redis -p 6379:6379 -d redis

$ git clone https://github.com/meopin-top/convey-your-mind-WSS.git meopin-wss
$ cd meopin-wss
$ docker build -t meopin-wss .
$ docker run meopin-wss
```

## How to use

1. `http://localhost:3000/ws/{project_id}`으로 웹소켓 연결
2. 연결이 성립되면 현재 프로젝트의 데이터를 불러옴.
```json
{
    "status": "active",
    "project_id": "abc",
    "contents": [
        {
            "user_id": "byungwook",
            "content_id": "123",
            "content_type": "text",
            "x": 100,
            "y": 100,
            "width": 200,
            "height": 300,
            "text": "hello world",
            "image_url": ""
        }
    ]
}
```
3. 이후 데이터를 전송하게되면 같은 프로젝트에 접속중인 모든 유저에게 데이터 브로드캐스팅 됨
```json
{
    "project_id": "abc",
    "user_id": "qwer",
    "content": {
        "user_id": "qwer",
        "content_id": "1234",
        "content_type": "text",
        "x": 100,
        "y": 100,
        "width": 200,
        "height": 300,
        "text": "hello zz",
        "image_url": ""
    }
}
```
