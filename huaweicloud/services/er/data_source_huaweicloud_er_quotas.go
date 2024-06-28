// Generated by PMS #233
package er

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceErQuotas() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceErQuotasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The quota type to be queried.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The instance ID.`,
			},
			"route_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The route table ID.`,
			},
			"quotas": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All quotas that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of quota used.`,
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unit of usage.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The quota type.`,
						},
						"limit": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of available quotas, ` + "`" + `-1` + "`" + ` means unlimited.`,
						},
					},
				},
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

func dataSourceErQuotasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newQuotasDSWrapper(d, meta)
	showQuotasRst, err := wrapper.ShowQuotas()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showQuotasToSchema(showQuotasRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API ER GET /v3/{project_id}/enterprise-router/quotas
func (w *QuotasDSWrapper) ShowQuotas() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "er")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/enterprise-router/quotas"
	params := map[string]any{
		"type":         w.Get("type"),
		"erId":         w.Get("instance_id"),
		"routeTableId": w.Get("route_table_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *QuotasDSWrapper) showQuotasToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("quotas", schemas.SliceToList(body.Get("quotas"),
			func(quotas gjson.Result) any {
				return map[string]any{
					"used":  quotas.Get("used").Value(),
					"unit":  quotas.Get("unit").Value(),
					"type":  quotas.Get("quota_key").Value(),
					"limit": quotas.Get("quota_limit").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}