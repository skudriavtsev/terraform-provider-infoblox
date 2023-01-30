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

func testAccCheckARecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_a_record" {
			continue
		}
		connector := meta.(ibclient.IBConnector)
		objMgr := ibclient.NewObjectManager(connector, "terraform_test", "test")
		rec, _ := objMgr.GetARecordByRef(rs.Primary.ID)
		if rec != nil {
			return fmt.Errorf("record not found")
		}

	}
	return nil
}

func testAccARecordCompare(t *testing.T, resPath string, expectedRec *ibclient.RecordA) resource.TestCheckFunc {
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

		rec, _ := objMgr.GetARecordByRef(res.Primary.ID)
		if rec == nil {
			return fmt.Errorf("record not found")
		}

		if !reflect.DeepEqual(rec.Name, expectedRec.Name) {
			return fmt.Errorf(
				"'fqdn' does not match: got '%s', expected '%s'",
				safePtrValue(rec.Name), safePtrValue(expectedRec.Name))
		}
		if !reflect.DeepEqual(rec.Ipv4Addr, expectedRec.Ipv4Addr) {
			return fmt.Errorf(
				"'ipv4address' does not match: got '%s', expected '%s'",
				safePtrValue(rec.Ipv4Addr), safePtrValue(expectedRec.Ipv4Addr))
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

func TestAccResourceARecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "infoblox_a_record" "foo"{
						fqdn = "name1.test.com"
						ip_addr = "10.0.0.2"
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccARecordCompare(t, "infoblox_a_record.foo", &ibclient.RecordA{
						Ipv4Addr: utils.StringPtr("10.0.0.2"),
						Name:     utils.StringPtr("name1.test.com"),
						View:     "default",
						Ttl:      utils.Uint32Ptr(0),
						UseTtl:   utils.BoolPtr(false),
						Ea:       nil,
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_a_record" "foo2"{
						fqdn = "name2.test.com"
						ip_addr = "192.168.31.31"
						ttl = 10
						dns_view = "nondefault_view"
						comment = "test comment 1"
						ext_attrs = jsonencode({
						  "Location" = "New York"
						  "Site" = "HQ"
						})
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccARecordCompare(t, "infoblox_a_record.foo2", &ibclient.RecordA{
						Ipv4Addr: utils.StringPtr("192.168.31.31"),
						Name:     utils.StringPtr("name2.test.com"),
						View:     "nondefault_view",
						Ttl:      utils.Uint32Ptr(10),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("test comment 1"),
						Ea: ibclient.EA{
							"Location": "New York",
							"Site":     "HQ",
						},
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_a_record" "foo2"{
						fqdn = "name3.test.com"
						ip_addr = "10.10.0.1"
						ttl = 155
						dns_view = "nondefault_view"
						comment = "test comment 2"
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccARecordCompare(t, "infoblox_a_record.foo2", &ibclient.RecordA{
						Ipv4Addr: utils.StringPtr("10.10.0.1"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "nondefault_view",
						Ttl:      utils.Uint32Ptr(155),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("test comment 2"),
					}),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_a_record" "foo2"{
						fqdn = "name3.test.com"
						ip_addr = "10.10.0.1"
						dns_view = "nondefault_view"
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccARecordCompare(t, "infoblox_a_record.foo2", &ibclient.RecordA{
						Ipv4Addr: utils.StringPtr("10.10.0.1"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "nondefault_view",
						UseTtl:   utils.BoolPtr(false),
					}),
				),
			},
		},
	})
}
