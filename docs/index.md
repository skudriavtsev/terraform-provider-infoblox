# Infoblox IPAM Driver for Terraform

## Resources

There are resources for the following objects, supported by the plugin:

-   Network view
-   Network container
-   Network
-   A-record
-   AAAA-record
-   PTR-record
-   CNAME-record
-   Host record

Network container and network resources have two versions: IPv4 and IPv6. In
addition, there are two operations which are implemented as resources:
IP address allocation and IP address association with a network host
(ex. VM in a cloud environment); they have three versions: IPv4
and IPv6 separately, and IPv4/IPv6 combined.

The recommendation is to avoid using separate IPv4 and IPv6 versions of
IP allocation and IP association resources.
This recommendation does not relate to network container and network-related resources.

To work with DNS records a user must ensure that appropriate DNS zones
exist on the NIOS side, because currently the plugin does not support
creating a DNS zone.

Every resource has common attributes: 'comment' and 'ext_attrs'.
'comment' is text which describes the resource. 'ext_attrs' is a set of
NIOS Extensible Attributes attached to the resource, read more on this
attribute in a separate clause.

For DNS-related resources there is 'ttl' attribute as well, it specifies
TTL value (in seconds) for appropriate record. There is no default
value, zone's TTL is used if omitted. TTL value of 0 (zero) means
caching should be disabled for this record.

All the resources have 'comment' and 'ext_attrs' attributes,
additionally DNS-related records have 'ttl' attribute. They are all
optional. In this document, a resource's description implies that there
may be no explicit note in the appropriate clauses.

## Data sources

There are data sources for the following objects:

- A-record
- CNAME-record
- IPv4 Network

## Importing existing resources

Terraform has the capability to import existing infrastructure. This allows users to take resources they have created by some other means and bring it under Terraform management.

As of now, Infoblox provider plugin lacks this: you have to write full resource's definition yourself.

In general, the process of importing an existing resource looks like this:

- write a new Terraform file (ex. a-record-imported.tf) with the content:
  ```
  resource "infoblox_a_record" "a_rec_1_imported" {
    fqdn = "rec-a-1.imported.test.com"
    ip_addr = "192.168.1.2"
    ttl = 10
    comment = "A-record to be imported"
    ext_attrs = jsonencode({
      "Location" = "New office"
    })
  }
  ```
- get a reference for the resource you want to import (ex. by using `curl` tool)
- issue a command of the form `terraform import RESOURCE_TYPE.RESOURCE_NAME RESOURCE_REFERENCE`.
  Example: `terraform import infoblox_a_record.a_rec_1_imported record:a/ZG5zLmJpbmRfYSQuX2RlZmF1bHQub3JnLmV4YW1wbGUsc3RhdGljMSwxLjIuMy40:rec-a-1.imported.test.com/default`

Please note that if some of resources' properties (supported by the Infoblox provider plugin) are not defined or empty for the object on NIOS side, then appropriate resource's property must be empty or not defined. Otherwise you will get a difference in the resource's actual state and resource's description you specified, and thus you will get a resource's update performed on the next `terraform apply` command invocation, which will actually set the value of the property to the one which you defined (ex. empty value).
