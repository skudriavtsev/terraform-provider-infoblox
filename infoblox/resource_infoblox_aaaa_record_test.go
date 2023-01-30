package infoblox

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
)

func testAccCheckAAAARecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_aaaa_record" {
			continue
		}
		connector := meta.(ibclient.IBConnector)
		objMgr := ibclient.NewObjectManager(connector, "terraform_test", "test")
		rec, _ := objMgr.GetAAAARecordByRef(rs.Primary.ID)
		if rec != nil {
			return fmt.Errorf("record not found")
		}

	}
	return nil
}

func testAccAAAARecordCompare(t *testing.T, resPath string, expectedRec *ibclient.RecordAAAA) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, found := s.RootModule().Resources[resPath]
		if !found {
			return fmt.Errorf("not found: %s", resPath)
		}
		if res.Primary.ID == "" {
			return fmt.Errorf("ID is not set")
		}
		meta := testAccProvider.Meta()
		connector := meta.(ibclient.IBConnector)
		objMgr := ibclient.NewObjectManager(connector, "terraform_test", "test")

		rec, _ := objMgr.GetAAAARecordByRef(res.Primary.ID)
		if rec == nil {
			return fmt.Errorf("record not found")
		}

		if !reflect.DeepEqual(rec.Name, expectedRec.Name) {
			return fmt.Errorf(
				"'fqdn' does not match: got '%s', expected '%s'",
				safePtrValue(rec.Name), safePtrValue(expectedRec.Name))
		}
		if !reflect.DeepEqual(rec.Ipv6Addr, expectedRec.Ipv6Addr) {
			return fmt.Errorf(
				"'ipv6address' does not match: got '%s', expected '%s'",
				safePtrValue(rec.Ipv6Addr), safePtrValue(expectedRec.Ipv6Addr))
		}
		if rec.View != expectedRec.View {
			return fmt.Errorf(
				"'dns_view' does not match: got '%s', expected '%s'",
				rec.View, expectedRec.View)
		}
		if !reflect.DeepEqual(rec.UseTtl, expectedRec.UseTtl) {
			return fmt.Errorf(
				"the value of 'use_ttl' field does not match: got '%s', expected '%s'",
				safePtrValue(rec.UseTtl), safePtrValue(expectedRec.UseTtl))
		} else if rec.UseTtl != nil && *rec.UseTtl {
			if !reflect.DeepEqual(rec.Ttl, expectedRec.Ttl) {
				return fmt.Errorf(
					"the value of 'ttl' field does not match: got '%s', expected '%s'",
					safePtrValue(rec.Ttl), safePtrValue(expectedRec.Ttl))
			}
		}
		if !reflect.DeepEqual(rec.Comment, expectedRec.Comment) {
			return fmt.Errorf(
				"'comment' does not match: got '%s', expected '%s'",
				safePtrValue(rec.Comment), safePtrValue(expectedRec.Comment))
		}
		return validateEAs(rec.Ea, expectedRec.Ea)
	}
}

func TestAccResourceAAAARecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAAAARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "infoblox_aaaa_record" "foo"{
						fqdn = "name1.test.com"
						ipv6_addr = "2000::1"
						dns_view = "default"
						comment = "test comment 1"
						ext_attrs = jsonencode({
							"Tenant ID"="terraform_test_tenant"
							"Location"="Test loc"
							"Site"="Test site"
							"TestEA1"=["text1","text2"]
						})
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::1"),
						Name:     utils.StringPtr("name1.test.com"),
						View:     "default",
						Ttl:      utils.Uint32Ptr(0),
						UseTtl:   utils.BoolPtr(false),
						Comment:  utils.StringPtr("test comment 1"),
						Ea: ibclient.EA{
							"Tenant ID": "terraform_test_tenant",
							"Location":  "Test loc",
							"Site":      "Test site",
							"TestEA1":   []string{"text1", "text2"},
						},
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_aaaa_record" "foo2"{
						fqdn = "name3.test.com"
						ipv6_addr = "2000::3"
						ttl = 155
						dns_view = "default"
						comment = "test comment 2"
						ext_attrs = jsonencode({
							"Tenant ID"="terraform_test_tenant"
							"Location"="Test loc"
						})
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::3"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "default",
						Ttl:      utils.Uint32Ptr(155),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("test comment 2"),
						Ea: ibclient.EA{
							"Tenant ID": "terraform_test_tenant",
							"Location":  "Test loc",
						},
					}),
				),
			},
		},
	})
}
