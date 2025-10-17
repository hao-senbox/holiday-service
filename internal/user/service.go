package user

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"worktime-service/pkg/constants"
	"worktime-service/pkg/consul"

	"github.com/hashicorp/consul/api"
)

type UserService interface {
	GetUserInfor(ctx context.Context, userID string) (*UserInfor, error)
	GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error)
	GetCurrentUser(ctx context.Context) (*CurrentUser, error)
	GetTeacherInforByOrg(ctx context.Context, teacherID, orgID string) (*UserInfor, error)
}

type userService struct {
	client *callAPI
}

type callAPI struct {
	client       consul.ServiceDiscovery
	clientServer *api.CatalogService
}

var (
	mainService = "go-main-service"
)

func NewUserService(client *api.Client) UserService {
	mainServiceAPI := NewServiceAPI(client, mainService)
	return &userService{
		client: mainServiceAPI,
	}
}

func NewServiceAPI(client *api.Client, serviceName string) *callAPI {
	sd, err := consul.NewServiceDiscovery(client, serviceName)
	if err != nil {
		fmt.Printf("Error creating service discovery: %v\n", err)
		return nil
	}

	var service *api.CatalogService

	for i := 0; i < 10; i++ {
		service, err = sd.DiscoverService()
		if err == nil && service != nil {
			break
		}
		fmt.Printf("Waiting for service %s... retry %d/10\n", serviceName, i+1)
		time.Sleep(3 * time.Second)
	}

	if service == nil {
		fmt.Printf("Service %s not found after retries, continuing anyway...\n", serviceName)
	}

	if os.Getenv("LOCAL_TEST") == "true" {
		fmt.Println("Running in LOCAL_TEST mode — overriding service address to localhost")
		service.ServiceAddress = "localhost"
	}

	return &callAPI{
		client:       sd,
		clientServer: service,
	}
}

func (u *userService) GetCurrentUser(ctx context.Context) (*CurrentUser, error) {

	token, ok := ctx.Value(constants.Token).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getCurrentUser(token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", token)
	}

	return data, nil
}

func (u *userService) GetUserInfor(ctx context.Context, userID string) (*UserInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getUserInfor(userID, token)
	if err != nil {
		log.Printf("[userService] getUserInfor call error: %v", err)
		return nil, nil
	}

	if data == nil {
		log.Printf("[userService] getUserInfor: empty data for userID=%s", userID)
		return nil, nil
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil
}

func (u *userService) GetStudentInfor(ctx context.Context, studentID string) (*UserInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStudentInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:         fmt.Sprintf("%v", innerData["id"]),
		UserName:       fmt.Sprintf("%v", innerData["name"]),
		OrganizationID: fmt.Sprintf("%v", innerData["organization_id"]),
		Avartar:        avatar,
	}, nil

}

func (u *userService) GetTeacherInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getTeacherInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil

}

func (u *userService) GetListTeacherInfor(ctx context.Context, userID string) ([]*UserInfor, error) {

	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getListTeacherInfor(userID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", userID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	fmt.Printf("innerData: %v\n", innerData)
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	return nil, nil
}

func (u *userService) GetStaffInfor(ctx context.Context, studentID string) (*UserInfor, error) {
	token, ok := ctx.Value(constants.TokenKey).(string)

	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getStaffInfor(studentID, token)
	if err != nil {
		return nil, err
	}

	if data == nil {
		return nil, fmt.Errorf("no user data found for userID: %s", studentID)
	}

	innerData, ok := data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'data' field")
	}

	var avatar Avatar
	if rawAvatar, exists := innerData["avatar"].(map[string]interface{}); exists {
		avatar = Avatar{
			ImageID:  uint64(castToInt64(rawAvatar["image_id"])),
			ImageKey: fmt.Sprintf("%v", rawAvatar["image_key"]),
			ImageUrl: fmt.Sprintf("%v", rawAvatar["image_url"]),
			Index:    int(castToInt64(rawAvatar["index"])),
			IsMain:   castToBool(rawAvatar["is_main"]),
		}
	}

	return &UserInfor{
		UserID:   fmt.Sprintf("%v", innerData["id"]),
		UserName: fmt.Sprintf("%v", innerData["name"]),
		Avartar:  avatar,
	}, nil

}

func (u *userService) GetTeacherInforByOrg(ctx context.Context, teacherID, orgID string) (*UserInfor, error) {
	if u.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	token, ok := ctx.Value(constants.TokenKey).(string)
	if !ok {
		return nil, fmt.Errorf("token not found in context")
	}

	data, err := u.client.getTeacherInforByOrg(teacherID, orgID, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get teacher info: %w", err)
	}

	return parseUserInforSafely(data)
}

func (c *callAPI) getTeacherInforByOrg(teacherID, orgID string, token string) (map[string]interface{}, error) {

	if c == nil || c.client == nil || c.clientServer == nil {
		return nil, fmt.Errorf("client is not properly initialized")
	}

	endpoint := fmt.Sprintf("/v1/gateway/teachers/organization/%s/user/%s", orgID, teacherID)
	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		return nil, fmt.Errorf("API call failed: %w", err)
	}

	if res == "" {
		return nil, nil
	}

	var userData interface{}
	if err := json.Unmarshal([]byte(res), &userData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if userData == nil {
		return nil, nil
	}

	if myMap, ok := userData.(map[string]interface{}); ok {
		return myMap, nil
	}

	return nil, fmt.Errorf("unexpected response format")
}

func parseUserInforSafely(data map[string]interface{}) (*UserInfor, error) {
	if data == nil {
		return nil, nil
	}

	innerData, ok := safeGetMapString(data["data"])
	if !ok || innerData == nil {
		return nil, nil
	}

	avatar := parseAvatarSafely(innerData)

	return &UserInfor{
		UserID:   safeGetString(innerData["id"]),
		UserName: safeGetString(innerData["name"]),
		Avartar:  avatar,
	}, nil
}

func safeGetMapString(v interface{}) (map[string]interface{}, bool) {
	if v == nil {
		return nil, false
	}
	m, ok := v.(map[string]interface{})
	return m, ok
}

func safeGetString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func parseAvatarSafely(data map[string]interface{}) Avatar {
	var avatar Avatar

	if rawAvatar, exists := data["avatar"]; exists && rawAvatar != nil {
		if avatarMap, ok := safeGetMapString(rawAvatar); ok {
			avatar = Avatar{
				ImageID:  uint64(castToInt64(avatarMap["image_id"])),
				ImageKey: safeGetString(avatarMap["image_key"]),
				ImageUrl: safeGetString(avatarMap["image_url"]),
				Index:    int(castToInt64(avatarMap["index"])),
				IsMain:   castToBool(avatarMap["is_main"]),
			}
		}
	}

	return avatar
}

func (c *callAPI) getUserInfor(userID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/users/%s", userID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil

}

func (c *callAPI) getStudentInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/students/%s", studentID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil
}

func (c *callAPI) getTeacherInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/teachers/%s", studentID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil
}

func (c *callAPI) getStaffInfor(studentID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/staffs/%s", studentID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil
}

func (c *callAPI) getListTeacherInfor(userID string, token string) (map[string]interface{}, error) {

	endpoint := fmt.Sprintf("/v1/gateway/teachers/get-by-user/%s", userID)

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var userData interface{}

	err = json.Unmarshal([]byte(res), &userData)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	myMap := userData.(map[string]interface{})

	return myMap, nil
}

func (c *callAPI) getCurrentUser(token string) (*CurrentUser, error) {

	endpoint := "/v1/user/current-user/"

	header := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token,
	}

	res, err := c.client.CallAPI(c.clientServer, endpoint, http.MethodGet, nil, header)
	if err != nil {
		fmt.Printf("Error calling API: %v\n", err)
		return nil, err
	}

	var data APIGateWayResponse[CurrentUser]
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		fmt.Printf("Error unmarshalling response: %v\n", err)
		return nil, err
	}

	return &data.Data, nil
}

func castToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case float64:
		return int64(val)
	case int:
		return int64(val)
	default:
		return 0
	}
}

func castToBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "1"
	default:
		return false
	}
}
