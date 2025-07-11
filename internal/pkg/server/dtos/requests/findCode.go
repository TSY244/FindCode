package requests

type APIScanRequest struct {
	Token  string `form:"token"`
	GitURL string `form:"gitUrl" binding:"required"`
	//Type                string   `form:"type" binding:"required"`
	IsUseAi             bool     `form:"isUseAi"`
	IsUseAiPrompt       bool     `form:"isUseAiPrompt"`
	AiPrompt            string   `form:"aiPrompt"`
	IsReturnBool        bool     `form:"isReturnBool"`
	Model               string   `form:"aiModel"`
	PermissionModel     string   `form:"permissionModel"`
	ProjectType         string   `form:"projectType"`
	AuthenticationCodes []string `form:"authenticationCodes"`
}
