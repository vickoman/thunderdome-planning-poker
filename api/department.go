package api

import (
	"net/http"
	"strings"

	"github.com/StevenWeathers/thunderdome-planning-poker/model"
	"github.com/gorilla/mux"
)

type departmentResponse struct {
	Organization     *model.Organization `json:"organization"`
	Department       *model.Department   `json:"department"`
	OrganizationRole string              `json:"organizationRole"`
	DepartmentRole   string              `json:"departmentRole"`
}

type departmentTeamResponse struct {
	Organization     *model.Organization `json:"organization"`
	Department       *model.Department   `json:"department"`
	Team             *model.Team         `json:"team"`
	OrganizationRole string              `json:"organizationRole"`
	DepartmentRole   string              `json:"departmentRole"`
	TeamRole         string              `json:"teamRole"`
}

// handleGetOrganizationDepartments gets a list of departments associated to the organization
// @Summary Get Departments
// @Description get list of organizations departments
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID to get departments for"
// @Success 200 object standardJsonResponse{data=[]model.Department}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments [get]
func (a *api) handleGetOrganizationDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Departments := a.db.OrganizationDepartmentList(OrgID, Limit, Offset)

		Success(w, r, http.StatusOK, Departments, nil)
	}
}

// handleGetDepartmentByUser gets a department with user role
// @Summary Get Department
// @Description Gets an organization department with users role
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID to get"
// @Success 200 object standardJsonResponse{data=departmentResponse}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId} [get]
func (a *api) handleGetDepartmentByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		DepartmentRole := r.Context().Value(contextKeyDepartmentRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := a.db.DepartmentGet(DepartmentID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentResponse{
			Organization:     Organization,
			Department:       Department,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
		}

		Success(w, r, http.StatusOK, result, nil)
	}
}

// handleCreateDepartment handles creating an organization department
// @Summary Create Department
// @Description Create an organization department
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID to create department for"
// @Param name body string true "the department name"
// @Success 200 object standardJsonResponse{data=model.Department}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments [post]
func (a *api) handleCreateDepartment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)

		OrgName := keyVal["name"].(string)
		OrgID := vars["orgId"]
		NewDepartment, err := a.db.DepartmentCreate(OrgID, OrgName)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, NewDepartment, nil)
	}
}

// handleGetDepartmentTeams gets a list of teams associated to the department
// @Summary Get Department Teams
// @Description Gets a list of organization department teams
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID to get teams for"
// @Success 200 object standardJsonResponse{data=[]model.Team}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams [get]
func (a *api) handleGetDepartmentTeams() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Teams := a.db.DepartmentTeamList(DepartmentID, Limit, Offset)

		Success(w, r, http.StatusOK, Teams, nil)
	}
}

// handleGetDepartmentUsers gets a list of users associated to the department
// @Summary Get Department Users
// @Description get a list of organization department users
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Success 200 object standardJsonResponse{data=[]model.DepartmentUser}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users [get]
func (a *api) handleGetDepartmentUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		Limit, Offset := getLimitOffsetFromRequest(r, w)

		Users := a.db.DepartmentUserList(DepartmentID, Limit, Offset)

		Success(w, r, http.StatusOK, Users, nil)
	}
}

// handleCreateDepartmentTeam handles creating an department team
// @Summary Create Department Team
// @Description Create a department team
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param name body string true "the team name"
// @Success 200 object standardJsonResponse{data=model.Team}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams [post]
func (a *api) handleCreateDepartmentTeam() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		keyVal := getJSONRequestBody(r, w)

		TeamName := keyVal["name"].(string)
		DepartmentID := vars["departmentId"]
		NewTeam, err := a.db.DepartmentTeamCreate(DepartmentID, TeamName)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, NewTeam, nil)
	}
}

// handleDepartmentAddUser handles adding user to an organization department
// @Summary Add Department User
// @Description Add a department User
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param email body string true "the users email"
// @Param role body string true "the users department role" Enums(MEMBER, ADMIN)
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users [post]
func (a *api) handleDepartmentAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		keyVal := getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		DepartmentId := vars["departmentId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, err := a.db.DepartmentAddUser(DepartmentId, User.Id, Role)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentRemoveUser handles removing user from a department (and department teams)
// @Summary Remove Department User
// @Description Remove a department User
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param userId path string true "the user ID"
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/users/{userId} [delete]
func (a *api) handleDepartmentRemoveUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		vars := mux.Vars(r)
		DepartmentID := vars["departmentId"]
		UserID := vars["userId"]

		err := a.db.DepartmentRemoveUser(DepartmentID, UserID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamAddUser handles adding user to a team so long as they are in the department
// @Summary Add Department Team User
// @Description Add a User to Department Team
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param teamId path string true "the team ID"
// @Param email body string true "the users email"
// @Param role body string true "the users team role" Enums(MEMBER, ADMIN)
// @Success 200 object standardJsonResponse{}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams/{teamId}/users [post]
func (a *api) handleDepartmentTeamAddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		keyVal := getJSONRequestBody(r, w)

		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]
		UserEmail := strings.ToLower(keyVal["email"].(string))
		Role := keyVal["role"].(string)

		User, UserErr := a.db.GetUserByEmail(UserEmail)
		if UserErr != nil {
			Failure(w, r, http.StatusInternalServerError, Errorf(ENOTFOUND, "USER_NOT_FOUND"))
			return
		}

		_, DepartmentRole, roleErr := a.db.DepartmentUserRole(User.Id, OrgID, DepartmentID)
		if DepartmentRole == "" || roleErr != nil {
			Failure(w, r, http.StatusInternalServerError, Errorf(EUNAUTHORIZED, "DEPARTMENT_USER_REQUIRED"))
			return
		}

		_, err := a.db.TeamAddUser(TeamID, User.Id, Role)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Success(w, r, http.StatusOK, nil, nil)
	}
}

// handleDepartmentTeamByUser gets a team with users roles
// @Summary Get Department Team
// @Description Get a department team with users role
// @Tags organization
// @Produce  json
// @Param orgId path string true "the organization ID"
// @Param departmentId path string true "the department ID"
// @Param teamId path string true "the team ID"
// @Success 200 object standardJsonResponse{data=departmentTeamResponse}
// @Failure 500 object standardJsonResponse{}
// @Security ApiKeyAuth
// @Router /organizations/{orgId}/departments/{departmentId}/teams/{teamId} [get]
func (a *api) handleDepartmentTeamByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.config.OrganizationsEnabled {
			Failure(w, r, http.StatusBadRequest, Errorf(EINVALID, "ORGANIZATIONS_DISABLED"))
			return
		}
		OrgRole := r.Context().Value(contextKeyOrgRole).(string)
		DepartmentRole := r.Context().Value(contextKeyDepartmentRole).(string)
		TeamRole := r.Context().Value(contextKeyTeamRole).(string)
		vars := mux.Vars(r)
		OrgID := vars["orgId"]
		DepartmentID := vars["departmentId"]
		TeamID := vars["teamId"]

		Organization, err := a.db.OrganizationGet(OrgID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Department, err := a.db.DepartmentGet(DepartmentID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		Team, err := a.db.TeamGet(TeamID)
		if err != nil {
			Failure(w, r, http.StatusInternalServerError, err)
			return
		}

		result := &departmentTeamResponse{
			Organization:     Organization,
			Department:       Department,
			Team:             Team,
			OrganizationRole: OrgRole,
			DepartmentRole:   DepartmentRole,
			TeamRole:         TeamRole,
		}

		Success(w, r, http.StatusOK, result, nil)
	}
}
