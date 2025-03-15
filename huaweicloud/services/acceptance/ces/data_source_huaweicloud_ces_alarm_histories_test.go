package ces

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceCesAlarmHistories_basic(t *testing.T) {
	dataSource := "data.huaweicloud_ces_alarm_histories.test"
	dc := acceptance.InitDataSourceCheck(dataSource)
	rName := acceptance.RandomAccResourceName()
	baseConfig := testDataSourceCesAlarmHistories_base(rName)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceCesAlarmHistories_basic(baseConfig),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.record_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.alarm_id"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.name"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.status"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.level"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.begin_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.end_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.last_alarm_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.first_alarm_time"),
					resource.TestCheckResourceAttrSet(dataSource, "alarm_histories.0.data_points.#"),
					resource.TestCheckOutput("is_default_filter_useful", "true"),
					resource.TestCheckOutput("is_alarm_id_filter_useful", "true"),
					resource.TestCheckOutput("is_name_filter_useful", "true"),
					resource.TestCheckOutput("is_level_filter_useful", "true"),
				),
			},
		},
	})
}

func testDataSourceCesAlarmHistories_basic(config string) string {
	return fmt.Sprintf(`
%[1]s

locals {
  alarm_id = huaweicloud_ces_alarmrule.test.id
  name     = huaewicloud_ces_alarmrule.test.alarm_name
  level    = huaewicloud_ces_alarmrule.test.condition.alarm_level
}

data "huaweicloud_ces_alarm_histories" "test" {
  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_default_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_histories.test.alarm_histories) >= 1 
}

data "huaweicloud_ces_alarm_histories" "filter_by_alarm_id" {
  alarm_id = huaweicloud_ces_alarmrule.test.id

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_alarm_id_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_histories.filter_by_alarm_id.alarm_histories) >= 1 && alltrue([
    for alarm in data.huaweicloud_ces_alarm_histories.filter_by_alarm_id.alarm_histories[*] : alarm.alarm_id == local.alarm_id
  ])
}

data "huaweicloud_ces_alarm_histories" "filter_by_name" {
  name = huaweicloud_ces_alarmrule.test.alarm_name

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_name_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_histories.filter_by_name.alarm_histories) >= 1 && alltrue([
    for alarm in data.huaweicloud_ces_alarm_histories.filter_by_name.alarm_histories[*] : alarm.alarm_name == local.name
  ])
}

data "huaweicloud_ces_alarm_histories" "filter_by_level" {
  name = huaweicloud_ces_alarmrule.test.alarm_name

  depends_on = [huaweicloud_ces_alarmrule.test]
}

output "is_level_filter_useful" {
  value = length(data.huaweicloud_ces_alarm_histories.filter_by_level.alarm_histories) >= 1 && alltrue([
    for alarm in data.huaweicloud_ces_alarm_histories.filter_by_level.alarm_histories[*] : alarm.condition.alarm_level == local.
  ])
}
`, config)
}

func testDataSourceCesAlarmHistories_base(name string) string {
	return fmt.Sprintf(`
resource "huaweicloud_smn_topic" "test" {
  name = "smn-%[1]s"
}

resource "huaweicloud_obs_bucket" "bucket" {
  bucket        = "%[1]s"
  acl           = "public-read"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "tracker" {
  bucket_name        = huaweicloud_obs_bucket.bucket.bucket
  file_prefix        = "cts"
  lts_enabled        = false
  compress_type      = "gzip"
  is_sort_by_service = false

  tags = {
    foo = "bar"
  }
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name           = "rule-%[1]s"
  alarm_action_enabled = true
  alarm_type           = "ALL_INSTANCE"

  metric {
    namespace = "SYS.OBS"
  }

  resources {
    dimensions {
      name  = "bucket_name"
	  value = huaweicloud_obs_bucket.bucket.bucket
    }
  }

  condition  {
    alarm_level         = 2
    period              = 1
    filter              = "average"
    comparison_operator = ">="
    count               = 1
	value               = 10
	unit                = "Bytes"
	suppress_duration   = 300
    metric_name         = "capacity_total"
  }

  alarm_actions {
    type              = "notification"
    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }

  depends_on = [huaweicloud_cts_tracker.test]
}
`, name)
}
