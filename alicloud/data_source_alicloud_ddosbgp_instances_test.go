package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDdosbgpInstanceDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.alicloud_ddosbgp_instances.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAcc%d", rand),
		dataSourceDdosbgpInstanceConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alicloud_ddosbgp_instance.default.id}-fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}",
			"ids":        []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alicloud_ddosbgp_instance.default.name}-fake",
			"ids":        []string{"${alicloud_ddosbgp_instance.default.id}"},
		}),
	}

	var existDdosbgpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"ids.0":                      CHECKSET,
			"names.#":                    "1",
			"names.0":                    fmt.Sprintf("tf_testAcc%d", rand),
			"instances.#":                "1",
			"instances.0.name":           fmt.Sprintf("tf_testAcc%d", rand),
			"instances.0.type":           string(Enterprise),
			"instances.0.base_bandwidth": "20",
			"instances.0.bandwidth":      "101",
			"instances.0.ip_count":       "100",
			"instances.0.ip_type":        "IPv4",
			"instances.0.region":         "cn-hangzhou",
		}
	}

	var fakeDdosbgpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"names.#":     "0",
			"instances.#": "0",
		}
	}

	var ddosbgpInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDdosbgpInstanceMapFunc,
		fakeMapFunc:  fakeDdosbgpInstanceMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.DdosbgpSupportedRegions)

	}
	ddosbgpInstanceCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, allConf)
}

func dataSourceDdosbgpInstanceConfigDependence(name string) string {
	return fmt.Sprintf(`
    provider "alicloud" {
        endpoints {
            bssopenapi = "business.aliyuncs.com"
        }
    }

	resource "alicloud_ddosbgp_instance" "default" {
      name                    = "%s"
      base_bandwidth          = "20"	
      bandwidth               = "101"
      ip_count                = "100"
      ip_type                 = "IPv4"
      type                    = "Enterprise"
	}`, name)
}