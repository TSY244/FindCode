# todo
- [x] 完成命令行ai 模式
- [x] 完成 -go_target 参数模式
- [x] 将ai 模式的结果通过前端展示
- [x] 给server 添加自动识别框架代码
- [x] 展示项目扫描结果
- [x] 完成命令行的替换，不需要提供框架类型​​​
- [x] 完成指纹提取​​​
- [x] 完成更深层次的子调用代码
- [x] 展示测试ai 检测结果​​​
- [x] 添加server 的自动扫描，完成命令行的替换，不需要提供框架类型​
- [ ] 将常规漏洞扫描结果存放进去数据库​​​
- [ ] 将ai 扫描结果存放进入数据库​​​
- [ ] 完成数据库的定时清理​​​
- [ ] 完成前端查看已经完成的扫描结果​​​
- [ ] 完成前端查看详细扫描结果​​​
- [x] 修复getFuncCode 使用path 的问题​​​


# 简介

该项目可以用于
- 在指定的函数中寻找特定的代码
- 只过滤出想获取的函数
- 查找指定的代码片段
- 具备大模型接入能力，可以帮助白盒理解函数的调用链
- 当接入大模型能力的时候，具备默认优化的提示词，用来判断是否越权风险

FindCodeServer 版本是command 版本和server 版本合二为一

FindCodeCommand 版本纯命令行版本



# 测试
## gin

``` yaml
task_name: gin扫描测试

mode:
  - go

go_mode_target_rule:
  rule: 'true'

func_rule:
  - func_name: # v2 版本trpc
      rule: 'beginWithLower()'
    param_name:
      rules:
        - size: 1
          rules:
            - 'equal("c")'
    param_type:
      rules:
        - size: 1
          rules:
            - 'equal("*gin.Context")'
    return_type:
      rules:
        - size: 0
          rules:
    recv_type:
      rule: 'beginStr("*")'
    recv_name:
      rule: 'true'


path:
  rule: 'true' # 可在此添加目录规则，如 'contain("service")'
  deepsize: 10

file:
  rule: 'true'  # 默认已排除 test 文件和非 go 文件，无需额外配置
```

上面就是所用到的函数规则，下面展示代码的扫描结果

```text

/product/code/routers/staff_key.go 扫描结果如下
开始行数:结尾行数:函数名字
40:61:apply
63:84:reApply
86:108:mySelf
110:132:myApprove
134:155:enable
157:174:delete
176:198:encrypt
200:222:decrypt


/product/routers/staff_key.go 扫描结果如下
开始行数:结尾行数:函数名字
40:61:apply
63:84:reApply
86:108:mySelf
110:132:myApprove
134:155:enable
157:174:delete
176:198:encrypt
200:222:decrypt


/product/routers/base_info.go 扫描结果如下
开始行数:结尾行数:函数名字
29:36:products


/product/code/routers/service.go 扫描结果如下
开始行数:结尾行数:函数名字
47:63:apply
65:83:mySelf
85:104:maybe
106:119:getById
121:133:delete
135:151:update
153:169:reApply
171:184:getSecretKey
186:198:newSecretKey
200:216:enable
218:235:myApprove
237:254:encrypt
256:273:decrypt
275:292:keyVersion


/product/code/routers/ioa.go 扫描结果如下
开始行数:结尾行数:函数名字
34:37:token
39:46:getStaff
48:50:error
52:54:ok


/product/routers/apply.go 扫描结果如下
开始行数:结尾行数:函数名字
34:50:approve
52:69:getApproveInfo


/product/routers/service.go 扫描结果如下
开始行数:结尾行数:函数名字
47:63:apply
65:83:mySelf
85:104:maybe
106:119:getById
121:133:delete
135:151:update
153:169:reApply
171:184:getSecretKey
186:198:newSecretKey
200:216:enable
218:235:myApprove
237:254:encrypt
256:273:decrypt
275:292:keyVersion


/product/routers/service_key.go 扫描结果如下
开始行数:结尾行数:函数名字
44:60:keyApply
62:79:keyList
81:98:myKeyList
100:116:reApplyKey
118:134:keyEnable
136:149:getKeyById
151:163:deleteKey
165:181:applyToken
183:199:tokenEnable
201:218:myToken
220:237:myTokenApprove


/product/routers/key.go 扫描结果如下
开始行数:结尾行数:函数名字
43:59:create
61:78:list
80:97:oldList
99:115:enable
117:133:planDelete
135:147:cancelDelete
149:161:delete
163:179:update
181:194:getKeyValue
196:212:updateKeyValue
214:231:myApprove


/product/code/routers/base_info.go 扫描结果如下
开始行数:结尾行数:函数名字
29:36:products


/product/code/routers/key.go 扫描结果如下
开始行数:结尾行数:函数名字
43:59:create
61:78:list
80:97:oldList
99:115:enable
117:133:planDelete
135:147:cancelDelete
149:161:delete
163:179:update
181:194:getKeyValue
196:212:updateKeyValue
214:231:myApprove


/product/code/routers/service_key.go 扫描结果如下
开始行数:结尾行数:函数名字
44:60:keyApply
62:79:keyList
81:98:myKeyList
100:116:reApplyKey
118:134:keyEnable
136:149:getKeyById
151:163:deleteKey
165:181:applyToken
183:199:tokenEnable
201:218:myToken
220:237:myTokenApprove


/product/routers/application.go 扫描结果如下
开始行数:结尾行数:函数名字
45:61:apply
63:79:openapiApply
81:99:mySelf
101:119:maybe
121:134:getById
136:148:delete
150:166:update
168:184:reApply
186:199:getSecretKey
201:213:newSecretKey
215:231:enable
233:250:myApprove


/product/code/routers/apply.go 扫描结果如下
开始行数:结尾行数:函数名字
34:50:approve
52:69:getApproveInfo


/product/routers/logs.go 扫描结果如下
开始行数:结尾行数:函数名字
33:50:list


/product/code/routers/application.go 扫描结果如下
开始行数:结尾行数:函数名字
45:61:apply
63:79:openapiApply
81:99:mySelf
101:119:maybe
121:134:getById
136:148:delete
150:166:update
168:184:reApply
186:199:getSecretKey
201:213:newSecretKey
215:231:enable
233:250:myApprove


/product/routers/ioa.go 扫描结果如下
开始行数:结尾行数:函数名字
34:37:token
39:46:getStaff
48:50:error
52:54:ok


/product/code/routers/logs.go 扫描结果如下
开始行数:结尾行数:函数名字
33:50:list

```


## go swagger (很少使用的框架)
拿到这个项目，只需要，抽离出api 的规律，就能找到所有的api 函数。
```yaml
task_name: go_swagger_扫描测试

mode:
  - go

go_mode_target_rule:
  rule: 'true'

func_rule:
  - func_name: # v2 版本trpc
      rule: 'beginWithUpper()'
    param_name:
      rules:
        - size: 0
          rules:
    param_type:
      rules:
        - size: 0
          rules:
    return_type:
      rules:
        - size: 1
          rules:
            - 'beginStr("operations")'
    recv_type:
      rule: 'equal("*ApiManager")'
    recv_name:
      rule: 'equal("apiManager")'


path:
  rule: 'true' # 可在此添加目录规则，如 'contain("router")'
  deepsize: 10

file:
  rule: 'true'  # 默认已排除 test 文件和非 go 文件，无需额外配置
```

扫描结果:
```text
不存在filter


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetClientsGet.go 扫描结果如下
开始行数:结尾行数:函数名字
15:41:GetClientsGetHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetLoginEventPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
15:70:GetLoginEventGetApi


/product/code/go/augeu/backEnd/internal/pkg/web/api/LoginPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
17:98:LoginPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/UploadUserInfoPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
16:56:UploadUserInfoPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetFileReportPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
16:94:GetFileReportPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetVersionApi.go 扫描结果如下
开始行数:结尾行数:函数名字
10:16:GetVersionApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/UploadLoginEventPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
14:69:UploadLoginEventApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/RegisterPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
15:67:RegisterPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/UploadRdpEventPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
16:69:UploadRdpEventPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetClientIdPostApi.go 扫描结果如下
开始行数:结尾行数:函数名字
18:104:GetClientIdPostApiHandlerFunc


/product/code/go/augeu/backEnd/internal/pkg/web/api/GetRulesGetApi.go 扫描结果如下
开始行数:结尾行数:函数名字
12:30:GetRulesGetHandlerFunc

```
这个是自己的项目中的api
https://github.com/TSY244/augeu

## trpc 腾讯开源的rpc 框架
规则我是从https://github.com/trpc-group/trpc-cmdline/tree/main/install/protobuf/asset_go
中提取出来。类似的这种使用工具生成的代码，都可以使用该项目准确定位到api
```yaml
task_name: trpc扫描测试

mode:
  - go

go_mode_target_rule:
  rule: 'true'

func_rule:
  - func_name: # v2 版本trpc
      rule: 'beginWithUpper()'
    param_name:
      rules:
        - size: 2
          rules:
            - 'equal("ctx")'
            - 'equal("req")'
    param_type:
      rules:
        - size: 2
          rules:
            - 'equal("context.Context")'
            - 'beginStr("*pb.")'
    return_type:
      rules:
        - size: 2
          rules:
            - 'beginStr("*pb.")'
            - 'equal("error")'
    recv_type:
      rule: 'endStr("Impl")'
    recv_name:
      rule: "true"
  - func_name: # v1 版本trpc
      rule: 'beginWithUpper()'
    param_name:
      rules:
        - size: 3
          rules:
            - 'equal("ctx")'
            - 'equal("req")'
            - 'equal("rsp")'
    param_type:
      rules:
        - size: 3
          rules:
            - 'equal("context.Context")'
            - 'beginStr("*pb.")'
            - 'beginStr("*pb.")'
    return_type:
      rules:
        - size: 1
          rules:
            - 'equal("error")'
    recv_type:
      rule: 'endStr("Impl")'
    recv_name:
      rule: "true"
  - func_name: # $method.ClientStreaming $method.ServerStreaming 和$method.ClientStreaming
      rule: 'beginWithUpper()'
    param_name:
      rules:
        - size: 1
          rules:
            - 'equal("stream")'
    param_type:
      rules:
        - size: 1
          rules:
            - 'reg("_[^ ]+Server")'
    return_type:
      rules:
        - size: 1
          rules:
            - 'equal("error")'
    recv_type:
      rule: 'endStr("Impl")'
    recv_name:
      rule: "true"
  - func_name: # $method.ServerStreaming
      rule: 'beginWithUpper()'
    param_name:
      rules:
        - size: 2
          rules:
            - 'equal("req")'
            - 'equal("stream")'
    param_type:
      rules:
        - size: 2
          rules:
            - 'beginStr("*pb.")'
            - 'reg("_[^ ]+Server")'
    return_type:
      rules:
        - size: 1
          rules:
            - 'equal("error")'
    recv_type:
      rule: 'endStr("Impl")'
    recv_name:
      rule: "true"


path:
  rule: 'true' # 可在此添加目录规则，如 'contain("router")'
  deepsize: 10

file:
  rule: 'true'  # 默认已排除 test 文件和非 go 文件，无需额外配置
```

扫描结果
```text

/product/code/go/.../internal/server/xxxx.go 扫描结果如下
开始行数:结尾行数:函数名字
30:37:GetActList
40:47:GetActInfo
50:57:GetProdActList
60:67:GetProdConfActList
70:77:GetActUnconfProd
80:87:GetActByLock
90:97:GetCheckActCnt
100:107:GetProdActSku
110:117:AddAct
120:127:OperateAct
130:137:CopyAct
140:147:ExcludeProduct
150:157:AddActSku
160:167:ActTimerJob
170:177:GetMultiProdAct
180:187:GetActSkuInfo
190:207:SendBoxActNotifyMsg
210:220:GetSkuStandardPercent
223:233:CheckActSkuUpdate
267:275:AdjustSkuBatchQuantity
```


# 其他尝试
## 基于静态的鉴权扫描
   本质是找到鉴权框架之后，添加规则
   `!contain("CheckAuth")`
   如果扫描出的代码中不存在这个CheckAuth 那么就很有可能存在越权。

## 调用链输出-未实现

## 恶意三方库引入


# 使用
## 命令行参数
### 静态代码扫描逻辑
1. 最简单的，通过制定项目是什么框架之后进行扫描

   >  ./FindCode -l /product/path/

   自动识别项目框架

2. 使用手动指定框架类型

   > ./FindCode -l /product/path/ -r rule_path

    这个地方使用自己的扫描规则
    默认的规则有find_gin_api.yaml，find_go_swagger_api.yaml，find_trpc_api.yaml

3. 添加扫描目标逻辑

   > ./FindCode -l /product/path/  -go_target 'contain("asdfasdf")'

   可以在-go_target 添加更加具体的逻辑，表示的是搜寻

### 结合ai 扫描
1. 结合api代码+一层下游代码 询问大模型
   > ./FindCode -ai true -l /product/path/ -r/-f data -prompt_file /Prompt.txt -ai_cycle num -o result.txt
   
   表示通过提供的Prompt 对每一个函数进行轮训 -prompt_file 表示提示词的文件的地址。

   需要添加两个%s
   第一个%s 表示api对应的code 源码
   第二个%s 表示api 下游，一层被调用的函数code 源码

   -ai_cycle 表示每一个函数询问的次数

   -o 输出文件到result.txt


2. 使用大模型询问，并得到结果，限定大模型返回bool
   通过大模型返回的bool 类型判读api+下游代码是否满足要求
   > ./FindCode -ai -l /product/path/ -r/-f data -prompt_file /Prompt.txt -ai_cycle num -o result.txt -return_true true



## docker-compose 
```
AiSk=xxxxx docker-compose up -d
```
需要添加一个环境变量

