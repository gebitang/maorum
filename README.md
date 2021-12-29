# MaoRum: 定投记录小助手

*Mixin & Rum 开发练手作品*

格式化信息，生成番茄记录。自动生成当日番茄记录图，可将对应的图发送到Rum的指定群组中。可获取指定群组的最新信息

### 快速入门

- 创建mixin机器人，[创建说明](https://prsdigg.com/articles/4bb154a0-30b7-478d-88ea-8d5d68a1bafd)
- 将私钥文件命名为`config.json`报错在`config`目录下
- 同目录下保持`rum.yml`文件，内容如下
- 下载编译好的文件直接启动，或下载[本仓库]()后只需`go run main.go`

```yaml
rum:
  host: 127.0.0.1 #
  port: 1799 # 节点与网络-> 节点参数 -> 端口
  cert.file: path\to\rum-app\resources\quorum_bin\certs\server.crt # rum安装目录下的文件 rum-app\resources\quorum_bin\certs\server.crt
  post.group.id: valid-post-rum-group-id # 群组详情 -> ID  分享tomato的群组
  read.group.id: valid-read-rum-group-id # 群组详情 -> ID  读取信息的群组
```

config.json文件内容类似——

```json
{
 "pin": "11111",
 "client_id": "xxx-xx-xx-xx-xxx",
 "session_id": "xxx-1e57-4bxxx65-xxx-xxxx",
 "pin_token": "xxxxx-xxxx",
 "private_key": "xxxxxxxx"
}
```

![Solstice Sun and Milky Way](https://s3-img.meituan.net/v1/mss_3d027b52ec5a4d589e68050845611e68/ff/n0/0m/zh/a8_156586.jpg) 

### 交互逻辑

- 回复`GAO`，答复所有可选项
- 回复`gao`，答复选项类型，选择对应的类型后，自动生成一条番茄记录
- 回复`mao`，生成当日的番茄图
- 回复`rum`，将当日的番茄图和总结信息发送到指定的Rum群组(需要配置rum信息)
- 回复`mur`(`rum`的倒序)，读取指定Rum群组的最新一条信息

### 下一步计划

- 番茄项目格式信息可配置(目前采用的是"锤子手机文青“配色)
- 支持指定开始时间(目前的逻辑是`记录的时刻 - 时长`为默认的开始时间)
- 部署机器人，直接添加机器人即可使用(依赖rum的linux环境编译和群组添加功能，应该已经支持，我还没进行实践)

### 开发手记

[BotDemo](https://github.com/gebitang/botdemo)是最初的学习，使用golang模仿官方视频的Node实现，学习mixin机器人api。

后续还给golang的sdk提了[两个pr](https://github.com/MixinNetwork/bot-api-go-client/commits?author=gebitang)。发送群组按钮的功能还没有发布，所有依赖需要下载官方的master版本之后，使用本地代码替代。

`replace github.com/MixinNetwork/bot-api-go-client => D:\gogogo\bot-api-go-client`

常用的开发操作都有涉及：字符串处理、文件处理、画图、数据库操作、网络交互、配置信息

- 色彩转换使用[golang playground的colors包](https://github.com/go-playground/colors)
- 画图操作根据plot的[wiki实现](https://github.com/gonum/plot/wiki/Creating-Custom-Plotters:-A-tutorial-on-creating-custom-Plotters)
- 数据库使用gorm包，采用Sqllite数据库
- rum的交互是本地生成了[quorum](https://github.com/rumsystem/quorum)的swagger文件

>启动swagger服务依赖swaggo生成的docs包。scripts文件夹下提供了对应的生成脚本，结果我自己研究了半天，只生成了个空的swagger接口文档。
>最终在wsl环境下生成，复用即可。

陆续开发了一段时间，应该一开始就使用git记录，查看本地文件的修改时间，可以确定最早开工时间是`2021-12-20 11:25`

### 最初的设想 

[MAO台番茄酱厂][purpose]

Mixin Autonomous Organization

需要开源一个这样的Mixin机器人：交互归集每天的

- 一次性自定义收集项目，例如：1读书、2工作、3玩、4其他
- 每次交互： 1-25 表示刚刚完成了25分钟的读书番茄，自动记录这一次番茄时钟。机器人根据收到的时间戳指定这个番茄实际发生的时间段
- 回复“酱”(当然也允许自定义关键字)，生成一天的番茄时钟文字+图表
- 回复“rum”，可以将当天的总结内容+图片自动同步到Rum上的群组中
- 本地有rum环境，或者云上的rum环境，知道本地的证书文件、网络端口、

- 还可以格式化记录所有的信息，后续机器人本地还可以支持浏览器访问，可以时间段进行聚合展示


还可以建一个群，群规就一条：只能发每天的番茄总结图，不能任何其他内容。
所有的内容都在本地完成。Rum目前可以做到完全匿名效果，所以也不涉及到隐私
一群人默默前行的感觉应该很酷的吧？

[purpose]: https://s3-img.meituan.net/v1/mss_3d027b52ec5a4d589e68050845611e68/ff/n0/0m/zh/az_156567.jpg@596w_1l.jpg