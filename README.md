### Pigeond

Pigeond服务启动后会监听本地/var/run/pigeond.socket，
pigeon客户端连接到socket后通过向socket中发送指令，等待
pigeond启动server时定义的callback函数处理完指令后将结
果通过socket连接返回给pigeon。

* socket package
    里面定义LaunchServer来启动unix socket server，
    启动时需要传入msg callback函数用于受到消息后进行
    后续的处理。
    server中用于启动socket server，handle中定义
    socket connection的处理函数， error中封装了服务
    的error类型。

* log package
    配置日志。

* tasks package
    封装TaskProxy和tasks，TaskProxy用于根据msg内容
    执行不同的task，并返回处理结果。

***

##### 其他

* 利用socat连接unix socket进行调试
    ```
    sudo socat UNIX-CONNECT:/var/run/pigeond.socket -
    ```