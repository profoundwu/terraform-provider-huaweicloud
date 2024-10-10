package vpn

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/vpn"
)

func getUserFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	region := acceptance.HW_REGION_NAME

	getConnectionProduct := "vpn"
	client, err := conf.NewServiceClient(getConnectionProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating Connection Client: %s", err)
	}

	return vpn.GetUser(client, state.Primary.Attributes["vpn_server_id"], state.Primary.ID)
}

func TestAccUser_basic(t *testing.T) {
	var obj interface{}

	name := acceptance.RandomAccResourceName()
	rName := "huaweicloud_vpn_user.test"

	rc := acceptance.InitResourceCheck(
		rName,
		&obj,
		getUserFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckVPNP2cServer(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testConnection_basic(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name),
					resource.TestCheckResourceAttr(rName, "vpn_type", "STATIC"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "86400"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-128"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "3600"),
					resource.TestCheckResourceAttrPair(rName, "gateway_id",
						"huaweicloud_vpn_gateway.test", "id"),
					resource.TestCheckResourceAttrPair(rName, "gateway_ip",
						"huaweicloud_vpn_gateway.test", "master_eip.0.id"),
					resource.TestCheckResourceAttrPair(rName, "customer_gateway_id",
						"huaweicloud_vpn_customer_gateway.test", "id"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar"),
				),
			},
			{
				Config: testConnection_update(name, ipAddress),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(rName, "name", name+"-update"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ikepolicy.0.lifetime_seconds", "172800"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.authentication_algorithm", "sha2-512"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.encryption_algorithm", "aes-256"),
					resource.TestCheckResourceAttr(rName, "ipsecpolicy.0.lifetime_seconds", "7200"),
					resource.TestCheckResourceAttr(rName, "tags.key", "val"),
					resource.TestCheckResourceAttr(rName, "tags.foo", "bar-update"),
				),
			},
			{
				ResourceName:      rName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"psk",
				},
			},
		},
	})
}
