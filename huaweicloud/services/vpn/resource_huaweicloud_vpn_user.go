package vpn

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var userNonUpdatableParams = []string{"vpn_server_id", "name"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}
// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}/reset-password
func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		UpdateContext: resourceUserUpdate,
		ReadContext:   resourceUserRead,
		DeleteContext: resourceUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(userNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The VPN server ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user name.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user password.`,
			},
			"description": {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      `The description of the user.`,
				DiffSuppressFunc: utils.SuppressCaseDiffs,
			},
			"user_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user group ID.`,
			},
			"user_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user group name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createUserHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users"
		createUserProduct = "vpn"
	)
	createUserClient, err := conf.NewServiceClient(createUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createUserPath := createUserClient.Endpoint + createUserHttpUrl
	createUserPath = strings.ReplaceAll(createUserPath, "{project_id}", createUserClient.ProjectID)
	createUserPath = strings.ReplaceAll(createUserPath, "{vpn_server_id}", d.Get("vpn_server_id").(string))

	createUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createUserOpt.JSONBody = utils.RemoveNil(buildCreateUserBodyParams(d))
	createUserResp, err := createUserClient.Request("POST", createUserPath, &createUserOpt)
	if err != nil {
		return diag.Errorf("error creating VPN user: %s", err)
	}

	createUserRespBody, err := utils.FlattenResponse(createUserResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("user.id", createUserRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VPN user: ID is not found in API response")
	}
	d.SetId(id)
	// The creation interface is asynchronous. If creation fails, the generated information will be removed.
	// Wait for a while to check if the creation is successful.
	time.Sleep(25 * time.Second)

	return resourceUserRead(ctx, d, meta)
}

func buildCreateUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user": map[string]interface{}{
			"name":          d.Get("name"),
			"password":      d.Get("password"),
			"description":   d.Get("description"),
			"user_group_id": utils.ValueIgnoreEmpty(d.Get("user_group_id")),
		},
	}
	return bodyParams
}

func resourceUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	getUserProduct := "vpn"
	getUserClient, err := conf.NewServiceClient(getUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	getUserRespBody, err := GetUser(getUserClient, serverId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN user")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("user.name", getUserRespBody, nil)),
		d.Set("description", utils.PathSearch("user.description", getUserRespBody, nil)),
		d.Set("user_group_id", utils.PathSearch("user.user_group_id", getUserRespBody, nil)),
		d.Set("user_group_name", utils.PathSearch("user.user_group_name", getUserRespBody, nil)),
		d.Set("created_at", utils.PathSearch("user.created_at", getUserRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("user.updated_at", getUserRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetUser(client *golangsdk.ServiceClient, serverId, id string) (interface{}, error) {
	getUserHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}"
	getUserPath := client.Endpoint + getUserHttpUrl
	getUserPath = strings.ReplaceAll(getUserPath, "{project_id}", client.ProjectID)
	getUserPath = strings.ReplaceAll(getUserPath, "{vpn_server_id}", serverId)
	getUserPath = strings.ReplaceAll(getUserPath, "{user_id}", id)

	getUserOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getUserResp, err := client.Request("GET", getUserPath, &getUserOpt)

	if err != nil {
		return nil, err
	}

	getUserRespBody, err := utils.FlattenResponse(getUserResp)
	if err != nil {
		return nil, err
	}

	return getUserRespBody, nil
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	updateUserClient, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	updateUserhasChanges := []string{
		"description",
		"user_group_id",
	}

	if d.HasChanges(updateUserhasChanges...) {
		updateUserHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}"

		updateUserPath := updateUserClient.Endpoint + updateUserHttpUrl
		updateUserPath = strings.ReplaceAll(updateUserPath, "{project_id}", updateUserClient.ProjectID)
		updateUserPath = strings.ReplaceAll(updateUserPath, "{vpn_server_id}", d.Get("vpn_server_id").(string))
		updateUserPath = strings.ReplaceAll(updateUserPath, "{user_id}", d.Id())

		updateUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		updateUserOpt.JSONBody = utils.RemoveNil(buildUpdateUserBodyParams(d))
		_, err = updateUserClient.Request("PUT", updateUserPath, &updateUserOpt)
		if err != nil {
			return diag.Errorf("error updating VPN user: %s", err)
		}
	}

	if d.HasChange("password") {
		resetPasswordUserHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}/reset-password"

		resetPasswordUserPath := updateUserClient.Endpoint + resetPasswordUserHttpUrl
		resetPasswordUserPath = strings.ReplaceAll(resetPasswordUserPath, "{project_id}", updateUserClient.ProjectID)
		resetPasswordUserPath = strings.ReplaceAll(resetPasswordUserPath, "{vpn_server_id}", d.Get("vpn_server_id").(string))
		resetPasswordUserPath = strings.ReplaceAll(resetPasswordUserPath, "{user_id}", d.Id())

		resetPasswordUserOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resetPasswordUserOpt.JSONBody = map[string]interface{}{
			"new_password": d.Get("password"),
		}
		_, err = updateUserClient.Request("POST", resetPasswordUserPath, &resetPasswordUserOpt)
		if err != nil {
			return diag.Errorf("error resetting VPN user password: %s", err)
		}
	}

	return resourceUserRead(ctx, d, meta)
}

func buildUpdateUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user": map[string]interface{}{
			"description":   utils.ValueIgnoreEmpty(d.Get("description")),
			"user_group_id": utils.ValueIgnoreEmpty(d.Get("user_group_id")),
		},
	}
	return bodyParams
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		deleteUserHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/users/{user_id}"
		deleteUserProduct = "vpn"
	)
	deleteUserClient, err := conf.NewServiceClient(deleteUserProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	deleteUserPath := deleteUserClient.Endpoint + deleteUserHttpUrl
	deleteUserPath = strings.ReplaceAll(deleteUserPath, "{project_id}", deleteUserClient.ProjectID)
	deleteUserPath = strings.ReplaceAll(deleteUserPath, "{vpn_server_id}", d.Get("vpn_server_id").(string))
	deleteUserPath = strings.ReplaceAll(deleteUserPath, "{user_id}", d.Id())

	deleteConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteUserClient.Request("DELETE", deleteUserPath, &deleteConnectionOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN user")
	}

	return nil
}
