# 配合 go-get-v2-chrome-plugin 对中国主流网站视频进行下载

读取Chrome插件[go-get-v2-chrome-plugin](https://github.com/schwarzeni/go-get-v2-chrome-plugin)生成对配置文件，对视频进行下载

支持网站：爱奇艺，腾讯，bilibili，优酷

**看看代码就行，爱奇艺，腾讯，bilibili，优酷的视频接口未来可能会改变，但是本人因为要准备准备考研了所以不会去维护了的**

---

## 环境需求

使用前先使用go对src/github.com/schwarzeni/go-get-v2进行编译，将bin目录加入到环境变量中

需要[ffmpeg](https://www.ffmpeg.org/download.html)安装并加入到环境变量中

需要有bash环境，windows尤其请注意，可以使用git安装时自带的模拟终端

---

## 执行

执行go-get-v2就可以了，以下为参数列表

- -h 帮助

- -w 并行下载数量，默认为20

- -p **必须填写** 配置文件路径，此文件由chrome插件[go-get-v2-chrome-plugin]()自动生成

- -y 不确认保存路径，不填写的话在程序执行的初始需要确认以下文件保存路径
