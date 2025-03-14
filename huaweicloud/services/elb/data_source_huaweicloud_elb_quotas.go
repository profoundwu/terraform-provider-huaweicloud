// Generated by PMS #578
package elb

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceElbQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceElbQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the project ID.`,
			},
			"loadbalancer": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the load balancer quota.`,
			},
			"listener": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the listener quota.`,
			},
			"l7policy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the forwarding policy quota.`,
			},
			"pool": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the backend server group quota.`,
			},
			"member": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the backend server quota.`,
			},
			"security_policy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the custom security policy quota.`,
			},
			"ipgroup": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the IP address group quota.`,
			},
			"ipgroup_max_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of IP addresses that can be added to an IP address group.`,
			},
			"healthmonitor": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the health check quota.`,
			},
			"certificate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the certificate quota.`,
			},
			"ipgroup_bindings": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of listeners that can be associated with an IP address group.`,
			},
			"listeners_per_loadbalancer": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of listeners that can be associated with a load balancer.`,
			},
			"listeners_per_pool": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of listeners that can be associated with a backend server group.`,
			},
			"l7policies_per_listener": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of forwarding policies that can be configured for a listener.`,
			},
			"ipgroups_per_listener": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of IP address groups that can be associated with a listener.`,
			},
			"members_per_pool": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of backend servers in a backend server group.`,
			},
			"pools_per_l7policy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of backend server groups that can be used by a forwarding policy.`,
			},
			"condition_per_policy": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum number of forwarding rules per forwarding policy.`,
			},
		},
	}
}

type QuotasDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newQuotasDSWrapper(d *schema.ResourceData, meta interface{}) *QuotasDSWrapper {
	return &QuotasDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceElbQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newQuotasDSWrapper(d, meta)
	showQuotaRst, err := wrapper.ShowQuota()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showQuotaToSchema(showQuotaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API ELB GET /v3/{project_id}/elb/quotas
func (w *QuotasDSWrapper) ShowQuota() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "elb")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/elb/quotas"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *QuotasDSWrapper) showQuotaToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("project_id", body.Get("quota.project_id").Value()),
		d.Set("loadbalancer", body.Get("quota.loadbalancer").Value()),
		d.Set("listener", body.Get("quota.listener").Value()),
		d.Set("l7policy", body.Get("quota.l7policy").Value()),
		d.Set("pool", body.Get("quota.pool").Value()),
		d.Set("member", body.Get("quota.member").Value()),
		d.Set("security_policy", body.Get("quota.security_policy").Value()),
		d.Set("ipgroup", body.Get("quota.ipgroup").Value()),
		d.Set("ipgroup_max_length", body.Get("quota.ipgroup_max_length").Value()),
		d.Set("healthmonitor", body.Get("quota.healthmonitor").Value()),
		d.Set("certificate", body.Get("quota.certificate").Value()),
		d.Set("ipgroup_bindings", body.Get("quota.ipgroup_bindings").Value()),
		d.Set("listeners_per_loadbalancer", body.Get("quota.listeners_per_loadbalancer").Value()),
		d.Set("listeners_per_pool", body.Get("quota.listeners_per_pool").Value()),
		d.Set("l7policies_per_listener", body.Get("quota.l7policies_per_listener").Value()),
		d.Set("ipgroups_per_listener", body.Get("quota.ipgroups_per_listener").Value()),
		d.Set("members_per_pool", body.Get("quota.members_per_pool").Value()),
		d.Set("pools_per_l7policy", body.Get("quota.pools_per_l7policy").Value()),
		d.Set("condition_per_policy", body.Get("quota.condition_per_policy").Value()),
	)
	return mErr.ErrorOrNil()
}
