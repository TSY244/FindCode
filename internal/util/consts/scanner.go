package consts

const (
	ArgsSize   = 2
	BeginLevel = 0
)

// function names
const (
	ContainFunc    = "contain"
	BeginStrFunc   = "beginStr"
	EndStrFunc     = "endStr"
	RegFunc        = "reg"
	BeginWithLower = "beginWithLower"
	BeginWithUpper = "beginWithUpper"
	EqualFunc      = "equal"
)

// mode
const (
	StrMode string = "str"
	GoMode  string = "go"
	AiMode  string = "ai"
)

// 用于传参数
const (
	IsUseCtxKey     = "isUseCtx"
	AiConfigKey     = "aiConfig"
	IsReturnBoolKey = "isReturnBool"
)

// 用于指定file path
const (
	GinRule       = "rule/find_gin_api.yaml"
	GoSwaggerRule = "rule/find_go_swagger_api.yaml"
	TrpcRule      = "rule/find_trpc_api.yaml"
)

const (
	FirstLevel = 1
	MaxLevel   = 4
)

const (
	MeetTheRules   = true  // 满足规则，可能是存在越权漏洞
	NoMeetTheRules = false // 不满足规则，可能不存在越权漏洞
)

const (
	ProductDetail    = "我是一个新手开发者，我简单对这个项目进行一个简介："
	PermissionDetail = "我是一个新手开发者，我觉得该项目的权限模型是: %s \n该权限可能有这些特点: %s"
)

const (
	PublicPermission              = "public"
	UserGroupPermission           = "userGroups"
	RBACPermission                = "RBAC"
	RBACWithConstraintsPermission = "RBACWithConstraints"
	ABACOrPBACPermission          = "ABAC/PBAC"
	DACPermission                 = "DAC"
	MACPermission                 = "MACPermission"
)

const (
	PublicPermissionDetail              = "无权限模型 / 公开访问\n\n描述： 所有用户（包括未登录的访客）对所有资源拥有完全相同的访问权限（通常是只读）。\n适用场景： 纯静态网站、公开信息展示页、某些只读API端点。\n特点： 实现最简单，无需用户系统或权限管理。"
	UserGroupPermissionDetail           = "用户组模型 (User Groups)\n\n描述： 用户被分配到不同的组（如 普通用户组、VIP用户组）。权限直接关联到组。用户加入哪个组就拥有该组的所有权限。\n适用场景： 权限划分相对简单、用户群体分类明确且权限需求同质化较高的系统（如某些内部工具、简单的内容管理系统）。\n特点： 比无权限模型有区分，管理比单独管理用户权限方便，但灵活性较低。注意： 有时会和角色模型混淆，但核心区别是权限直接绑定组，而非通过角色抽象。"
	RBACPermissionDetail                = "简单角色模型 (Simple RBAC - Role-Based Access Control Level 0)\n\n描述： 系统预定义几个固定角色（如 Admin, User, Guest）。每个角色拥有固定的权限集合。用户被直接分配一个或多个角色，从而获得该角色的权限。\n适用场景： 中小型系统，权限需求相对固定，角色数量不多（如博客后台：管理员、编辑、订阅者；小型SaaS应用：管理员、付费用户、免费用户）。\n特点： 最常见的入门级权限模型，易于理解和实现，管理相对简单（用户-角色关联）。你例子中的“只有admin和user的项目”就属于这种。\n变体： 只有 Admin 和 User 是最典型的简单角色模型。\n标准角色模型 (Hierarchical RBAC - RBAC Level 1 & 2)\n\n描述： 在简单角色模型基础上增加了：\n角色继承： 角色可以形成层级（如 管理员 > 部门经理 > 普通员工），高级角色自动继承低级角色的权限。\n权限-角色分离： 权限 (Permission) 被定义为独立的实体（如 create_post, delete_user）。角色 (Role) 是权限的集合。用户 (User) 被分配角色。三者关系清晰分离。\n适用场景： 中大型系统，组织结构有层级，权限需求较复杂且需要复用（如企业内部管理系统、ERP、CRM、复杂的后台管理系统）。你例子中的“多角色多权限的后台项目”通常指这种或更复杂的RBAC。\n特点： 灵活性、可扩展性、可管理性大大增强。权限变更只需修改角色定义，无需逐个修改用户。角色层级简化了权限分配。\n"
	RBACWithConstraintsPermissionDetail = "带约束的角色模型 (Constrained RBAC - RBAC Level 3)\n\n描述： 在标准RBAC基础上增加了约束，以实施更严格的安全策略，例如：\n职责分离： 互斥角色（如一个用户不能同时拥有 采购员 和 审批员 角色）。\n基数约束： 限制一个角色可分配的用户数，或一个用户可拥有的角色数。\n先决条件： 用户必须拥有角色A才能被分配角色B。\n适用场景： 对安全性和合规性要求极高的系统（如金融系统、医疗系统、政府系统）。\n特点： 提供更强的安全控制，但配置和管理更复杂。"
	ABACOrPBACPermissionDetail          = "基于属性的访问控制 (Attribute-Based Access Control - ABAC)\n\n描述： 权限决策基于属性的动态评估。属性可以来自：\n用户属性： 部门、职位、安全等级、国籍等。\n资源属性： 文件所有者、创建时间、敏感级别、标签等。\n环境属性： 当前时间、访问位置、设备类型、网络状态等。\n操作属性： 尝试执行的动作（读、写、删除）。\n策略： 管理员定义策略规则（Policy），规则由条件（基于属性）和结果（允许/拒绝）组成。例如：IF 用户.部门 == 资源.所属部门 AND 当前时间在 9:00-18:00 THEN ALLOW 编辑。\n适用场景： 需要非常细粒度、动态、上下文感知权限控制的复杂系统（如云平台 - AWS IAM Policies 的核心思想、大型企业复杂的数据访问控制、需要根据多种条件动态授权的场景）。\n特点： 提供极高的灵活性和表达能力，能实现非常精细的控制。但策略定义、管理和计算开销可能较大，理解和调试更复杂。\n基于策略的访问控制 (Policy-Based Access Control - PBAC)\n\n描述： PBAC 通常被视为 ABAC 的一个超集或紧密相关的概念。它更强调使用集中管理的、声明式的策略语言来定义访问规则。这些策略规则同样会利用用户、资源、环境等属性进行决策。核心在于策略的集中化、标准化和外部化管理。\n适用场景： 与 ABAC 高度重叠，尤其在大规模分布式系统、云原生应用、需要统一策略管理的场景中更强调 PBAC。\n特点： 强调策略与业务逻辑分离，策略可复用，易于统一审计和管理。实现通常依赖策略引擎（如 OPA - Open Policy Agent）。"
	DACPermissionDetail                 = "自主访问控制 (Discretionary Access Control - DAC)\n\n描述： 资源（如文件、文档）的所有者有权决定谁可以访问该资源以及拥有何种访问权限（读、写、执行）。通常通过访问控制列表实现。\n适用场景： 操作系统文件权限（Unix/Linux权限位、Windows ACL）、个人文档协作工具（如网盘设置文件共享）。\n特点： 控制权分散在资源所有者手中，灵活性高，但难以进行集中管理和实施统一的安全策略，安全性相对较低。"
	MACPermissionDetail                 = "强制访问控制 (Mandatory Access Control - MAC)\n\n描述： 由系统管理员或安全策略强制定义的安全标签（如密级：公开、秘密、机密；范畴：部门A、项目X）分配给所有主体（用户）和客体（资源）。访问决策严格基于主体标签和客体标签的比较规则（如“不向上读，不向下写”）。\n适用场景： 对国家安全、军事、情报等要求极高保密性的系统（如 SELinux, TrustedBSD）。\n特点： 提供最强的强制安全保证，但配置极其复杂，灵活性极低，用户无法更改权限。"
)
