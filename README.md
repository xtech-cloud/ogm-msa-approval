# StartKit

See [omo-msa-startkit](https://github.com/xtech-cloud/omo-msa-startkit)

# Protoc

See [omo-msp-account](https://github.com/xtech-cloud/omo-msp-account)

# Docker

See [omo-docker-account](https://github.com/xtech-cloud/omo-docker-account)


# 存储方式

## 七牛云

假设需要上传a.jpg到test存储桶中
先在客户端解析出a.jpg的MD5值为7e45f060cf1fe26b01cc27b277b34113,文件大小为342325(单位byte) 
客户端调用Object.Prepare("test", "7e45f060cf1fe26b01cc27b277b34113.jpg", 342325)准备一个对象，获取到engine(存储引擎),address(存储地址),accessToken(存储令牌)
客户端使用七牛API或SDK完成上传,上传时的文件名必须为文件的md5值，例如7e45f060cf1fe26b01cc27b277b34113.jpg
客户端调用Object.Sync("test", "7e45f060cf1fe26b01cc27b277b34113.jpg", "a-1.jpg")完成

再次上传a-副本.jpg到test存储桶中，a-副本.jpg和a.jpg是同一个文件
客户端解析出a-副本.jpg的MD5值为7e45f060cf1fe26b01cc27b277b34113,文件大小为342325(单位byte) 
客户端调用Object.Prepare("test", "7e45f060cf1fe26b01cc27b277b34113.jpg", 342325)准备一个对象，因为存储引擎中已经存在同样的文件，所以返回一个200的错误
客户端调用Object.Flush("test", "7e45f060cf1fe26b01cc27b277b34113.jpg", "a-2.jpg")完成
此时存储桶test中存在两个对象，a-1.jpg和a-2.jpg,都指向存储引擎中的7e45f060cf1fe26b01cc27b277b34113.jpg

## MinIO


# 消息订阅

- 地址
  omo.msa.account.notification

- 消息
  | Action | Head | Body|
  |:--|:--|:--|
  |/signup||uuid|
  |/signin|accessToken|uuid|
  |/signout|accessToken||
  |/reset/password|accessToken||
  |/profile/update|accessToken||
  |/profile/query|accessToken||
