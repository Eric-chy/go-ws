# go-ws
简单的go websocket示例，由于时间问题，只是简单的写了个demo，配置文件和目录没有进行划分，比 如路由可以单独的写一个文件，redis连接等可以放到配置文件里，比如推送的消息可以做成自定义的,**_另外服务重启或断开的时候需要自己清空一下redis里的连接，代码里面也没有写_**

如果只是单机版可以直接将连接存在map里而不需要存redis，如果多机部署将连接存redis的时候最好加上 机器的ip，推送的时候就能知道ws连接的是哪台机器，通过对应的机器去推送消息

user1.html和user2.html是前端websocket示例代码，表示两个不同的连接，服务端启动后，我们打开 user1.html就启动了用户1的websocket连接，打开user2.html就启动了用户2的连接，这个时候访问 http://localhost:8199/send-to-user?uid=1 **____**就会给用户1推送消息，访问http://localhost:8199/send-to-user?uid=2 就会给用户2推送消息，访问http://localhost:8199/send-to-users就会给用户1和2都推送消息

指定用户推送： http://localhost:8199/send-to-user?uid=1, uid=1表示推送给用户1，uid=2表 示推送给用户2，实际开发由于安全问题最好用token而不是uid

全体用户推送： http://localhost:8199/send-to-users
