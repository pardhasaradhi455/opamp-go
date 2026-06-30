package controller

import (
	"net/http"

	"github.com/open-telemetry/opamp-go/opamp-server/data"
	"github.com/open-telemetry/opamp-go/opamp-server/login/models"
	"github.com/open-telemetry/opamp-go/opamp-server/login/response"
	"github.com/open-telemetry/opamp-go/opamp-server/login/service"
)

type LoginController struct{}

func NewLoginController(version string) *LoginController {
	return &LoginController{}
}

func (c *LoginController) GetRoles(w http.ResponseWriter, r *http.Request) {

	claimsAny := r.Context().Value(service.ClaimsKey)
	if claimsAny == nil {
		response.New(w, r).
			Status(http.StatusUnauthorized).
			Log("missing claims in context").
			JSON(models.ErrorResponse{Message: "missing claims in context"})
		return
	}

	claims, ok := claimsAny.(*models.IDPClaims)
	if !ok {
		response.New(w, r).
			Status(http.StatusInternalServerError).
			Log("invalid claims type").
			JSON(models.ErrorResponse{Message: "invalid claims type"})
		return
	}

	roles, permissions := service.ExtractRolesAndPermissions(claims)

	response.New(w, r).
		Status(http.StatusOK).
		Log("roles fetched").
		JSON(map[string]any{
			"roles":       roles,
			"permissions": permissions,
		})
}

func (c *LoginController) GetAgents(w http.ResponseWriter, r *http.Request) {

	var agentsInfo []models.AgentResponse

	agentsMap := data.AllAgents.GetAllAgentsReadonlyClone()

	for _ ,agent := range agentsMap {
		agentResponse, _ := service.ExtractAgentInfo(*agent)
		agentsInfo = append(agentsInfo, agentResponse)
	}

	response.New(w, r).
		Status(http.StatusOK).
		Log("agents fetched").
		JSON(agentsInfo)
}
