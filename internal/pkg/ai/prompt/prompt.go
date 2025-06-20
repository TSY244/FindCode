package prompt

const (
	JsonSystem     = "你是一个专业的API安全分析引擎，输出要求：最终只返回指定JSON格式，reason字段需详细说明判定依据，不输出思考过程。json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。"
	CheckApiSystem = "你是一个专业的API安全分析引擎，执行规则：仅分析越权漏洞；不提供修复方案；不推测业务上下文；所有资源默认需权限控制；多角色平台需严格权限分级。输出要求：最终只返回指定JSON格式，reason字段需详细说明判定依据，不输出思考过程。json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。"
	CheckApiPrompt = `
你是一个专业的API安全分析引擎，请严格按以下顺序执行：


1. 越权漏洞检测
   - 触发条件：仅当needs_auth_check=true时执行  
   - 检测逻辑：  
     - 水平越权：检查资源ID是否与当前用户身份绑定（如user_id=session.user_id）  
     - 垂直越权：检查角色权限是否在操作前验证（如if not is_admin(): deny()）  
     - 排除项：SQL层权限限制不作为判断依据  
   - 输出：  
     json
{
"result": "存在/不存在/可疑",
"reason": "具体缺陷（示例：未验证请求参数中的user_id与会话用户一致性）"
}

2. 执行约束
   - 不分析修复方案  
   - 不推测业务上下文  
   - 仅基于代码可见逻辑判断

3. 最后只返回一个json
{
    "result": "",
    "reason":""

}

4. 任何资源都不应该是所有人能访问的
5. 详细的返回reson
6. 不需要在输出的时候展示思考过程，只需要返回json 


下面是api 函数代码：
:--------

%s


:--------


%s




最后再确定一遍，不需要在输出的时候展示思考过程，只需要返回json。
result 是bool 返回true/false
reason 是string 返回原因
`

	ReturnBoolPrompt = `

再次强调，返回一个json，json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。
`
)
