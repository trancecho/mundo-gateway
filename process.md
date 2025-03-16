现在已经完成：
1.gateway的初始化。service和api的注册。
3.注册serviec和api的路由
2.跑通pingdemo http转发方式(
4.连接池(使用client)
2.跑通ping demo grpc转发方式(
服务注册sdk
5.上线网关（zyb doing）
grpc ping demo 跑通
实现：删除服务顺带删除API(zyb)
实现：不能同时插入多个相同路由API(zyb)
修复路径匹配问题，从遍历数据库改成匹配第一个/之后内容
sdk注册服务之后需要调用网关的刷新


正在进行：
http连接池修复
api模块的controller改造适配数据结构。

未来计划：
api自动注册的diff算法
2.解耦注册中心
3.服务注册sdk
一些todo的优化。一些遍历可能可以性能优化。
网关CICD
