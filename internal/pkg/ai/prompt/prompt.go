package prompt

import (
	"ScanIDOR/pkg/logger"
	"os"
)

const (
	JsonSystem     = "你是一个专业的API安全分析引擎，输出要求：最终只返回指定JSON格式，reason字段需详细说明判定依据，不输出思考过程。json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。"
	CheckApiPrompt = `

下面是需要检测代码 api 函数代码：

:--------

%s


:--------

下面是需要检测代码 api 调用的函数代码：
%s




最后再确定一遍，不需要在输出的时候展示思考过程，只需要返回json。
result 是bool 返回true/false
当存在漏洞的时候返回false 
不存在漏洞返回true

reason 是string 返回给出result 的原因，使用中文
`

	ReturnBoolPrompt = `
再次强调，返回一个json，json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。
`
)

const (
	base = "你是一个专业的API安全分析引擎，执行规则：仅分析越权漏洞；不提供修复方案；不推测业务上下文；多角色平台需严格权限分级。输出要求：最终只返回指定JSON格式，reason字段需详细说明判定依据，不输出思考过程。json 拥有两个字段，result 是一个bool 类型。reason 是为什么给出这个result。某些API接口是公开的，任何用户都可以访问，无需进行越权检查。请确保你能够识别这些接口，并将其排除在越权检查之外。我提供的子调用代码可能是同名函数，并非真正的zi"
)

var (
	CheckApiSystem  = base + codePrompt + ruleConstraints
	ruleConstraints = ""
	codePrompt      = ""
)

const (
	DefaultRuleConstraintsPath = "prompt/rulePrompt.txt"
	DefaultCodePrompt          = "prompt/codePrompt.txt"
)

func init() {
	ruleConstraintsFile, err := os.ReadFile(DefaultRuleConstraintsPath)
	if err != nil {
		ruleConstraints = "1. 检查垂直越权，水平越权\n1.1 水平越权（Horizontal Privilege Escalation）\n定义：攻击者访问或操作同级别用户的资源或数据。\n示例：用户A和用户B权限相同，但通过修改URL参数（如/user?id=B改为/user?id=A），用户B可以查看或修改用户A的个人信息。\n常见场景：\n订单ID、用户ID等参数未校验归属权。\n直接使用可预测的参数（如递增的数字ID）。\n\n1.2 垂直越权（Vertical Privilege Escalation）\n定义：攻击者获取更高权限的功能或数据（如普通用户提升为管理员）。\n示例：普通用户通过伪造Cookie或API请求访问管理员后台（如/admin/delete_user）。\n常见场景：\n功能接口未校验角色权限。\n前端隐藏控件或API路由暴露高权限操作。\n\n2. 执行约束\n2.1 不分析修复方案\n2.2 不推测业务上下文\n2.3 仅基于代码可见逻辑判断\n2.4 返回格式\n最后只返回一个JSON格式：{ \"result\": \"\", \"reason\": \"\" }\n\n3. 资源访问控制\n    公开的资源不需要检查越权。\n\n\n4. 详细的返回reason\n    返回的reason使用中文。\n\n5. 输出要求\n    5.1 不需要在输出的时候展示思考过程，只需要返回JSON。\n\n6. 漏洞结果\n\n7. 存在漏洞result=false，不存在漏洞result=true。\n\n\n8. 其他越权漏洞原因\n8.1 访客数据未做主态判断。\n8.2 非好友拉取战绩时，返回结果包含用户Uin，存在泄露用户关系链风险。\n8.3 对用户传入的订单ID和用户ID进行权限校验，避免越权评论。\n8.4 参数校验不严谨，没有强校验优惠券ID和批次ID的对应关系。\n8.5 学习计划鲤鱼平台接口鉴权方式配置错误，未登录不能透传。\n8.6 直接原因是vcuid没有做登录态校验。\n8.7 优惠券冻结时未传入资源校验id及资源类型。\n8.8 接口未判断请求的用户是否有访问请求的资源的权限。\n\n9. 越权的修复方法\n9.1 授权时的权限校验添加了对各角色可授予的权限的明确校验。\n9.2 接口增加权限判断，只有管理员能拿到信息。\n9.3 直接修复：禁止往“内置数据源”类型进行注册。\n9.4 鉴权应用id与资源的应用id强校验。\n9.5 针对前端传到后端的作品ID，严格比对作品的用户ID和当前登录态中的用户ID是否一致。"
		logger.Error("ruleConstraints 加载失败，加载了默认的ruleConstraints")
	} else {
		ruleConstraints = string(ruleConstraintsFile)
		logger.Info("ruleConstraints 加载成功")
	}

	codePromptFile, err := os.ReadFile(DefaultCodePrompt)
	if err != nil {
		codePrompt = "越权可能的代码有（并非完整go 函数代码，只有部分片段）：\n1.1\n越权id 类漏洞示例：\nfunc (v *vPageCouponServiceImpl) TakeCoupon(ctx context.Context,\n     req *pb.TakeCouponReq, rsp *pb.TakeCouponRsp) (err error) {\n     ret, err := v.areaLogic.TakeCoupon(ctx, req.GetCouponBatchId()) // 出现ID越权的地方\n     if err != nil {\n          log.ErrorContextf(ctx, \\\"batchID: %s TakeCoupon err: %v\\\", req.GetCouponBatchId(), err)\n          return err\n     }\n\n原因： 可能被遍历id，导致隐藏卷被领取\n1.2\n缺少鉴权示例：\nif m, err = meta.GetMeta(ctx, app, attr.FileID); err != nil {\n          return err\n     }\n     attr.Version = m.RevisionVersion + 1\n     // 更新修订\n     if err := commitRevertRevision(ctx, app, attr); err != nil {\n这个代码中没有\n     if err := checkPermssion(ctx, app, m, attr); err != nil {\n          return err\n     }\n\n1.3 \n缺少对资料进行所属判断示例\nif resp := middleware2.CheckJwt(params.HTTPRequest, apiManager.s); resp != nil {\n               return resp\n          }\n          clientIds, err := apiManager.s.GetAllClientId()\n          clients, err := HostInfo.GetAgentsByClientIds(apiManager.s.DBM.DB, clientIds)\n          clientInfo := convert2.ArrayCopy(clients, convert2.DbHostinfo2moduleHostinfo)\n          return operations.NewGetGetClientsOK().WithPayload(&models.GetClientsResponse{\n               Clients: clientInfo,\n          })\n没有对client id 判断是否有权限读\n\n\n1.3 \n缺属性校验漏洞示例：\nif param.Req.GetModifySeq() == 0 {\n          return nil\n     }\n\n     // 打卡和作业不允许修改, 接龙和普通收集表都可以修改.\n     if param.Meta.GetFormType() == formpb.FormType_FORM_TYPE_SIGN ||\n          param.Meta.GetFormType() == formpb.FormType_FORM_TYPE_PAPER {\n缺少对文档属性的检查，没有判断该文档是否能改\n完整代码：\nif param.Req.GetModifySeq() == 0 {\n          return nil\n     }\n     // 检查可修改设置是否打开.\n     if !param.Meta.GetEnableModifySubmit() {\n          return errs.New(int(code.Code_ERR_PERMISSION_DENIED), \\\"modification not enabled\\\")\n     }\n     // 打卡和作业不允许修改, 接龙和普通收集表都可以修改.\n     if param.Meta.GetFormType() == formpb.FormType_FORM_TYPE_SIGN ||\n          param.Meta.GetFormType() == formpb.FormType_FORM_TYPE_PAPER {\n\n1.4\n越权刷评论漏洞：\nfunc (c *consumerImpl) checkPublishReq(ctx context.Context, req *BatchPublishCommentReq) error {\n     // 1. 校验图片URL是否合法\n     validURLs := c.cfg.Get().ValidPicURLs\n     for _, item := range req.Comments {\n          for _, pic := range item.PictureList {\n               if !isValidPicURL(pic.GetUrl(), validURLs) {\n                    // 任意一个图片不合法，返回错误\n                    return errors.New(\\\"picture url invalid\\\")\n               }\n          }\n     }\n}\n但是没有检查订单信息是否一致\n准确的代码如下：\nfunc (c *consumerImpl) checkPublishReq(ctx context.Context, req *BatchPublishCommentReq) error {\n     // 1. 校验图片URL是否合法\n     validURLs := c.cfg.Get().ValidPicURLs\n     for _, item := range req.Comments {\n          for _, pic := range item.PictureList {\n               if !isValidPicURL(pic.GetUrl(), validURLs) {\n                    // 任意一个图片不合法，返回错误\n                    return errors.New(\\\"picture url invalid\\\")\n               }\n          }\n     }\n}\n     // 2. 校验订单信息是否一致，避免越权刷评论等问题\n     orderInfo, err := c.orderMgr.GetExpressOrder(ctx, req.OrderID)\n     if err != nil {\n          log.ErrorContextf(ctx, \\\"get express order info fail, orderID:%s, err:%+v\\\", req.OrderID, err)\n          return err\n     }\n     if err = checkOrderInfo(orderInfo, req); err != nil {\n          return err\n     }\n     return nil\n\n1.5 \n缺少权限校验\n// PostCommentFeed 帖子回复\n// nolint\nfunc (c *CommentImpl) PostCommentFeed(ctx context.Context, req *pb.PostCommentRequest,\n     rsp *pb.PostCommentResponse) error {\n     // 1. 发表接口开关检查\n     if config.RainbowConfig.CommentSwitch == constant.OFF {\n          return errs.New(errcode.RetErrPubOFF, \\\"comment switch off\\\")\n     }\n     // 2. 入参校验\n     if !isValidInput(req) {\n          return errs.Newf(errcode.RetErrInvalidParam, \\\"input invalid, %s\\\", strs.Format(req))\n     }\n     // 3. 获取用户登录信息\n     userInfo, err := logicutil.GetUserInfo(ctx)\n     if err != nil {\n          log.ErrorContextf(ctx, \\\"[Comment] get UseInfo failed, err:<%v>\\\", err)\n          return errs.Newf(errcode.RetErrLogin, \\\"get user info failed, err: %+v\\\", err)\n     }\n     c.vuid = cast.ToString(userInfo.Vuid) // 目前只是对齐弹幕回复，该字段在贴子场景未用上\n     // 4. 数据转换\n     msg.go, ctxMap := transFeedCommentData(ctx, req, userInfo)\n     if msg.go == nil {\n          log.ErrorContextf(ctx, \\\"[Comment] trans terminal pb to MsgDomainInfo failed,req:<%s>\\\", strs.Format(req))\n          return errs.Newf(errcode.RetErrPBDecode, \\\"trans backend protocol failed, %s\\\", strs.Format(req))\n     }\n     // 5. 检查内容的图片url\n     if err := logicutil.CheckContentPictureURL(msg.go); err != nil {\n          log.ErrorContextf(ctx, \\\"[Publish] CheckContentPictureURL error, msg.go:%+v, err:%v\\\", msg.go, err)\n          return errs.Newf(errcode.RetErrPicURL, \\\"CheckContentPictureURL err:%v\\\", err)\n     }\n\n正确的代码如下：\n应该添加\n     // 3.1 权限校验\n     if err := publishFeedPermission(ctx, userInfo); err != nil {\n          return err\n     }\n\n// publishFeedPermission 是否有发表的权限\nfunc publishFeedPermission(ctx context.Context, user *unionmodel.UserCpInfo) error {\n     // 用户发表类型为创作号，检查创作号的准确性\n     vcuid, creatorStatus := logicutil.GetVcuid(ctx)\n     if vcuid != \\\"\\\" && creatorStatus != \\\"\\\" && vcuid != user.Vcuid {\n          logger.ErrorF(\\\"[publishFeedPermission] vcuid data inconsistency, ctx vcuid(%s) user vcuid(%s)\\\",\n               vcuid, user.Vcuid)\n          return errs.Newf(errcode.RetErrInvalidParam, \\\"vcuid data inconsistency, ctx vcuid(%s) user vcuid(%s)\\\",\n               vcuid, user.Vcuid)\n     }\n\n     return nil\n}\n1.6 \nvuid助力作为参数，没有加密。黑产可通过用遍历vuid的方式来助力进而获得奖励。\ncase ShareAndAssist:\n               errCode = ShareAssistByGuest(ctx, req, assistVuseridCode, assistVuserid, vuserid, actEndTime, reportData, actID, isPrepublish, modID)\n          case ShareStatus:\n               id, err := clients.DeCode(ctx, assistVuseridCode)\n          \n               reqParams := ReqParams{\n                    ActID:         req.ActId,\n                    ModID:         req.ModId,\n                    IsPrepublish:  strtoInt32(req.Context.CgiReqData.Query[\\\"is_prepublish\\\"]),\n                    Vuserid:       strtoInt64(req.Context.LoginInfo.Vuserid),\n                    MasterVuserid: id,\n               }\n               shareStatus, err := GetShareStatus(ctx, &reqParams)\n               if err != nil {\n                    lo\n                    g.Errorf(\\\"GetShareStatus err, err = %v\\\", err)\n正确代码\ncase ShareAndAssist:\n               errCode = ShareAssistByGuest(ctx, req, assistVuseridCode, assistVuserid, vuserid, actEndTime, reportData, actID, isPrepublish, modID)\n          case ShareStatus:\n               id, err := clients.DeCode(ctx, assistVuseridCode)\n               if actID >= svrConfig.GraysActID {\n                    if err != nil || cast.ToInt64(id) <= 10000 {\n                         id = \\\"\\\"\n                    }\n               } else {\n                    if err != nil || cast.ToInt64(id) <= 10000 {\n                         id = cast.ToString(assistVuserid)\n                    }\n               }\n               masterVuid := cast.ToInt64(id)\n               reqParams := ReqParams{\n                    ActID:         req.ActId,\n                    ModID:         req.ModId,\n                    IsPrepublish:  strtoInt32(req.Context.CgiReqData.Query[\\\"is_prepublish\\\"]),\n                    Vuserid:       strtoInt64(req.Context.LoginInfo.Vuserid),\n                    MasterVuserid: masterVuid,\n               }\n               shareStatus, err := GetShareStatus(ctx, &reqParams)\n               if err != nil {\n                    log.Errorf(\\\"GetShareStatus err, err = %v\\\", err)\n\n\n"
		logger.Error("codePromptFile 加载失败，加载了默认的 codePromptFile")
	} else {
		codePrompt = string(codePromptFile)
		logger.Info("codePromptFile 记载成功")
	}
}
