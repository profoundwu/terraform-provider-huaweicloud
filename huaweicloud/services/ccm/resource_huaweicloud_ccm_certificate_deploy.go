// Generated by PMS #319
package ccm

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SCM POST /v3/scm/certificates/{certificate_id}/deploy
func ResourceCertificateDeploy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateDeployCreate,
		ReadContext:   resourceCertificateDeployRead,
		DeleteContext: resourceCertificateDeployDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The CCM SSL certificate ID to be deployed.`,
			},
			"service_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The target service to which the certificate is pushed.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: `The list of resources to be deployed.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The resource ID.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The resource type.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The domain name to be deployed.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `The enterprise project ID to which the resources to be deployed.`,
						},
					},
				},
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The name of the project where the deployed resources are located.`,
			},
		},
	}
}

func buildDeployResourcesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	resourcesRaw := d.Get("resources").([]interface{})
	if len(resourcesRaw) == 0 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, 0, len(resourcesRaw))
	for _, v := range resourcesRaw {
		resourceMap := v.(map[string]interface{})
		bodyParams = append(bodyParams, map[string]interface{}{
			"id":                    utils.ValueIgnoreEmpty(resourceMap["id"]),
			"type":                  utils.ValueIgnoreEmpty(resourceMap["type"]),
			"domain_name":           utils.ValueIgnoreEmpty(resourceMap["domain_name"]),
			"enterprise_project_id": utils.ValueIgnoreEmpty(resourceMap["enterprise_project_id"]),
		})
	}
	return bodyParams
}

func buildCertificateDeployBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"service_name": d.Get("service_name"),
		"project_name": utils.ValueIgnoreEmpty(d.Get("project_name")),
		"resources":    buildDeployResourcesBodyParams(d),
	}
	return bodyParams
}

// The response status code of the deployment API is always 200. Determine whether the deployment is successful by
// parsing the response body.
// Example of a deployment failure response: {"failure_list":[{"resource":"XXX","failure_info":"XXX"}]}
// Example of a successful deployment response: {"failure_list":[]}
func flattenDeployRespBody(deployRespBody interface{}) error {
	var (
		failureList = utils.PathSearch("failure_list", deployRespBody, make([]interface{}, 0)).([]interface{})
		mErr        *multierror.Error
	)

	for _, v := range failureList {
		resource := utils.PathSearch("resource", v, "").(string)
		failureInfo := utils.PathSearch("failure_info", v, "").(string)
		mErr = multierror.Append(mErr, fmt.Errorf("resource: %s; error message: %s", resource, failureInfo))
	}
	return mErr.ErrorOrNil()
}

func resourceCertificateDeployCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/scm/certificates/{certificate_id}/deploy"
		product       = "scm"
		certificateID = d.Get("certificate_id").(string)
		serviceName   = d.Get("service_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SCM client: %s", err)
	}

	deployPath := client.Endpoint + httpUrl
	deployPath = strings.ReplaceAll(deployPath, "{certificate_id}", certificateID)
	deployOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCertificateDeployBodyParams(d)),
	}

	deployResp, err := client.Request("POST", deployPath, &deployOpt)
	if err != nil {
		return diag.Errorf("error deploying CCM SSL certificate (%s) to service (%s): %s",
			certificateID, serviceName, err)
	}

	deployRespBody, err := utils.FlattenResponse(deployResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if deployErr := flattenDeployRespBody(deployRespBody); deployErr != nil {
		return diag.Errorf("some errors occurred while deploying CCM SSL certificate (%s) to service (%s): %s",
			certificateID, serviceName, deployErr)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCertificateDeployRead(ctx, d, meta)
}

func resourceCertificateDeployRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCertificateDeployDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used only for deploy certificate. Deleting this resource
will not clear the certificate binding relationship, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
