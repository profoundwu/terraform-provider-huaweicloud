package cfw

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccDataSourceDomainNameIdParseIpList_basic(t *testing.T) {
	dataSource := "data.huaweicloud_cfw_domain_name_id_parse_ip_list.test"
	rName := acceptance.RandomAccResourceName()
	dc := acceptance.InitDataSourceCheck(dataSource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceDomainNameIdParseIpList_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					dc.CheckResourceExists(),
					resource.TestCheckResourceAttrSet(dataSource, "data.#"),
				),
			},
		},
	})
}

func testDataSourceDomainNameIdParseIpList_basic(name string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_cfw_domain_name_id_parse_ip_list" "test" {
  domain_address_id     = data.huaweicloud_cfw_domain_name_groups.filter_by_id.records[0].domain_names[0].domain_address_id
  domain_set_id         = huaweicloud_cfw_domain_name_group.test.id
  fw_instance_id        = "%[2]s"
}
`, testDataSourceDomainNameIdParseIpList_base(name), acceptance.HW_CFW_INSTANCE_ID)
}

func testDataSourceDomainNameIdParseIpList_base(name string) string {
	return fmt.Sprintf(`
%[1]s

resource "huaweicloud_cfw_domain_name_group" "test" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  name           = "%[3]s"
  type           = 1
  description    = "network domain name group"
  
  domain_names {
    domain_name = "www.baidu.com"
    description = "baidu"
  }
}

data "huaweicloud_cfw_domain_name_groups" "filter_by_id" {
  fw_instance_id = "%[2]s"
  object_id      = data.huaweicloud_cfw_firewalls.test.records[0].protect_objects[0].object_id
  group_id       = huaweicloud_cfw_domain_name_group.test.id
  
  depends_on = [huaweicloud_cfw_domain_name_group.test]
}
`, testAccDatasourceFirewalls_basic(), acceptance.HW_CFW_INSTANCE_ID, name)
}
