## ZeroMicroServices
本仓库为基于go-zero框架的微服务项目
参考仓库:https://github.com/Mikaelemmmm/go-zero-looklook

## 项目结构

项目目录结构如下：

- app：所有业务代码包含api、rpc以及mq（消息队列、延迟队列、定时任务）
- common：通用组件 error、middleware、interceptor、tool、ctxdata等
- data：该项目包含该目录依赖所有中间件(mysql、es、redis、grafana等)产生的数据，此目录下的所有内容应该在git忽略文件中，不需要提交。
- deploy：

    - filebeat: docker部署filebeat配置
    - go-stash：go-stash配置
    - nginx: nginx网关配置
    - prometheus ： prometheus配置
    - script：
        - gencode：生成api、rpc，以及创建kafka语句，复制粘贴使用
        - mysql：生成model的sh工具
    - goctl: 该项目goctl的template，goctl生成自定义代码模版，template用法可参考go-zero文档，复制到家目录下.goctl即可
- doc : 该项目系列文档
- modd.conf :  modd热加载配置文件


## 架构图

### 项目架构图

![gozerolooklook](./doc/images/gozerolooklook.png)

### 业务架构图
![gozerolooklook](./doc/images/go-zero-looklook-service.png)

## 开发建议

- app下放所有业务服务代码

- common放所有服务的公共基础库

- data项目依赖中间件产生的数据，实际开发中应该在git中忽略此目录以及此目录下产生的数据

- 生成api、rpc代码参考 deploy\script\gencode\gen.sh

- api项目中的.api文件我们做了拆分，统一放到每个api的desc文件夹下

- 生成model、错误处理时候使用了template重新定义，该项目用到的自定义的goctl的模版在项目deploy/goctl下
```powershell 

goctl template update -c model --home {TEMPLATE_PATH}

```