---
subcategory: "Policy - DNS"
layout: "nsxt"
page_title: "NSXT: nsxt_policy_dns_forwarder_zone"
description: A resource to configure DNS Forwarder Zone.
---

# nsxt_policy_dns_forwarder_zone

This resource provides a method for the management of DNS Forwarder Zone.

This resource is applicable to NSX Global Manager, NSX Policy Manager and VMC.

## Example Usage

```hcl
resource "nsxt_policy_dns_forwarder_zone" "test" {
  display_name     = "test"
  description      = "Terraform provisioned Zone"
  dns_domain_names = ["test.domain.org"]
  upstream_servers = ["33.14.0.2"]
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) Display name of the resource.
* `description` - (Optional) Description of the resource.
* `tag` - (Optional) A list of scope + tag pairs to associate with this resource.
* `nsx_id` - (Optional) The NSX ID of this resource. If set, this ID will be used to create the resource.
* `upstream_servers` - (Required) List of server IP addresses for this Forwarder Zone.
* `dns_domain_names` - (Optional) For conditional (FQDN) zones, a list of maximum 5 domains. For Default Forwarder Zone, this attribute should not be specified.
* `source_ip` - (Optional) The source IP address used by the DNS Forwarder zone.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - ID of the resource.
* `revision` - Indicates current revision number of the object as seen by NSX-T API server. This attribute can be useful for debugging.
* `path` - The NSX path of the policy resource.

## Importing

An existing object can be [imported][docs-import] into this resource, via the following command:

[docs-import]: https://www.terraform.io/cli/import

```
terraform import nsxt_policy_dns_forwarder_zone.test ID
```

The above command imports PolicyDnsForwarderZone named `test` with the NSX Dns Forwarder Zone ID `ID`.
