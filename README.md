# 配合 go-get-v2-chrome-plugin 对中国主流网站视频进行下载

读取Chrome插件[go-get-v2-chrome-plugin](https://github.com/schwarzeni/go-get-v2-chrome-plugin)生成对配置文件，对视频进行下载

支持网站：爱奇艺，腾讯，bilibili，优酷

**看看代码就行，爱奇艺，腾讯，bilibili，优酷的视频接口未来可能会改变，但是本人因为要准备准备考研了所以不会去维护了的**

相关总结个人博客文章见此，仅供参考

- [科普向：下载PC端b站视频的思路](http://blog.schwarzeni.com/2018/05/14/%E7%A7%91%E6%99%AE%E5%90%91%EF%BC%9A%E4%B8%8B%E8%BD%BDPC%E7%AB%AFb%E7%AB%99%E8%A7%86%E9%A2%91%E7%9A%84%E6%80%9D%E8%B7%AF/)

- [科普向：下载pc端爱奇艺视频的思路](http://blog.schwarzeni.com/2018/05/29/%E7%A7%91%E6%99%AE%E5%90%91%EF%BC%9A%E4%B8%8B%E8%BD%BDpc%E7%AB%AF%E7%88%B1%E5%A5%87%E8%89%BA%E8%A7%86%E9%A2%91%E7%9A%84%E6%80%9D%E8%B7%AF/)

- [网站视频抓取程序总结](http://blog.schwarzeni.com/2018/05/29/%E7%BD%91%E7%AB%99%E8%A7%86%E9%A2%91%E6%8A%93%E5%8F%96%E7%A8%8B%E5%BA%8F%E6%80%BB%E7%BB%93/)

---

## 环境需求

使用前先使用go对src/github.com/schwarzeni/go-get-v2进行编译，将bin目录加入到环境变量中

需要[ffmpeg](https://www.ffmpeg.org/download.html)安装并加入到环境变量中

需要有bash环境，windows尤其请注意，可以使用git安装时自带的bash模拟终端，获取bash路径后bin中的concat_file.sh第一行进行修改

---

## 执行

执行go-get-v2就可以了，以下为参数列表

- -h 帮助

- -w 并行下载数量，默认为20

- -p **必须填写** 配置文件路径，此文件由chrome插件[go-get-v2-chrome-plugin]()自动生成

- -y 不确认保存路径，不填写的话在程序执行的初始需要确认以下文件保存路径
