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
服务心跳检测注册卸载，sdk支持服务自动注册

正在进行：（todo
http连接池修复
api模块的controller改造适配数据结构。
地址注册后需要重启网关才生效，可能没放到内存
service自动注册
服务注册需要密码
jwt网关层校验（登录之后，发送事件，统一存到redis，网关会先和redis校验，再决定是否转发）（与此同时，就可以做账号下线功能）
黑名单：ip访问如果被认定为恶意（比如说一段时间内一直被权限拦截且访问频率过高），会加入黑名单（mysql）。
sdk发布打tag(还差http的，grpc的已经完成)
sdk的使用要更加便捷化。修复api重复注册就会报错的问题
修复：grpc反射描述符信息缓存（目前暴力解决了，但舍弃了缓存的优势）
fix: 创建service如果是已有，不会校验protocol

未来计划：
api自动注册的diff算法
2.解耦注册中心
3.服务注册sdk
一些todo的优化。一些遍历可能可以性能优化。
网关CICD
白名单，比如内部服务注册后自动白名单
