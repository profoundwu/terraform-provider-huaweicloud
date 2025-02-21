---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_domain_name_id_parse_ip_list"
description: |-
  Use this data source to get the DNS resolution result of a domain name.
---

# huaweicloud_cfw_domain_name_id_parse_ip_list

Use this data source to get the DNS resolution result of a domain name.

## Example Usage

```hcl
variable "fw_instance_id" {}
variable "domain_set_id" {}
variable "domain_address_id" {}

data "huaweicloud_cfw_domain_name_id_parse_ip_list" "test" {
  fw_instance_id    = var.fw_instance_id
  domain_set_id     = var.domain_set_id
  domain_address_id = var.domain_address_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `fw_instance_id` - (Required, String) Specifies the Firewall ID.

* `domain_set_id` - (Required, String) Specifies the domain name group ID.

* `domain_address_id` - (Required, String) Specifies the domain name ID.

* `address_type` - (Optional, Int) Specifies the address type.
  The valid value can be **0** (IPv4) or **1** (IPv6).

* `enterprise_project_id` - (Optional, String) Specifies the Enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The IP address resolution data of a domain name.
