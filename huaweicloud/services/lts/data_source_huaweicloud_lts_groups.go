// Generated by PMS #213
package lts

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

func DataSourceLtsGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLtsGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All log groups that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log group name.`,
						},
						"ttl_in_days": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The log expiration time(days).`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs to associate with the log group.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the log group belongs.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the log group, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

type GroupsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newGroupsDSWrapper(d *schema.ResourceData, meta interface{}) *GroupsDSWrapper {
	return &GroupsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceLtsGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newGroupsDSWrapper(d, meta)
	listLogGroupsRst, err := wrapper.ListLogGroups()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listLogGroupsToSchema(listLogGroupsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API LTS GET /v2/{project_id}/groups
func (w *GroupsDSWrapper) ListLogGroups() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "lts")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/groups"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *GroupsDSWrapper) listLogGroupsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("groups", schemas.SliceToList(body.Get("log_groups"),
			func(groups gjson.Result) any {
				return map[string]any{
					"id":          groups.Get("log_group_id").Value(),
					"name":        groups.Get("log_group_name").Value(),
					"ttl_in_days": groups.Get("ttl_in_days").Value(),
					// Using the `delete` method in the `ignoreSysEpsTag` method will change the original value,
					// so assign a value to `enterprise_project_id` parmater before assigning a value to `tags` parmater.
					"enterprise_project_id": groups.Get("tag._sys_enterprise_project_id").Value(),
					"tags":                  w.setLogGroupsTag(groups),
					"created_at":            w.setLogGroCreTime(groups),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*GroupsDSWrapper) setLogGroupsTag(data gjson.Result) map[string]interface{} {
	tags := data.Get("tag").Value()
	if tags != nil {
		return ignoreSysEpsTag(tags.(map[string]interface{}))
	}
	return nil
}

func (*GroupsDSWrapper) setLogGroCreTime(data gjson.Result) string {
	// The "creation_time" format is a UNIX timestamp.
	// Convert to the time corresponding to the local time zone of the computer.
	return utils.FormatTimeStampRFC3339(data.Get("creation_time").Int()/1000, false)
}
