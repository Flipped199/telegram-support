Telegram Support Bot

### 使用Telegram群组的话题功能实现与用户的双向联系

首先你需要有一个群组，建议设为**私有**，然后将你的bot加入群组并添加管理员，并且至少给予*话题管理*的权限。
到此、似乎已经没有更多需要配置的地方了。


### 说明

你可以在结束一个用户的会话后关闭话题，此时话题将显示一个小锁标志用以标明结束对话（会话结束后用户仍可以发送消息）。

假设你非要删除话题，这不会影响用户后续的联系，因为当用户联系bot时，若不存在对应的话题会重新创建一个话题。
但是我并不建议你删除话题，因为这样会导致用户回复消息引用失败（引用的源消息被删除）

我们实现了同步编辑消息的功能，或许可能存在一些bug，但经过粗略的测试，已经能够实现。如您在使用的过程中发现问题，欢迎issue。

在群组话题内使用以下命令:
```shell
/del # 对指定的消息回复/del命令会删除双方对话中关联的消息
```
~~`/ban # 封禁用户（用户无法通过机器人联系）`~~

~~`/unban # 解除封禁`~~
  
> 当回复消息时，引用了一条被删除的消息，bot将直接发送该消息到目标chat。

### Docker
```shell
docker run -d --name telegram-supprot-bot -v $PWD/config.toml:/app/conf/config.toml -v $PWD/data:/app/data flipped199/telegram-support-bot:latest
```