package databricks

import (
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDatabricksGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabricksGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabricksScimGroup(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDatabricksGroupExists("databricks_group.group"),
					resource.TestCheckResourceAttr(
						"databricks_group.group", "display_name", "TFGroupName"),
				),
			},
			{
				Config: testAccDatabricksScimGroupConfigUpdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckDatabricksGroupExists("databricks_group.group"),
					resource.TestCheckResourceAttr(
						"databricks_cluster.cluster", "display_name", "TFUpdatedGroupName"),
				),
			},
		},
	})
}

func testAccCheckDatabricksGroupExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set")
		}

		conn := testAccProvider.Meta().(*Client).scim

		_, err := conn.GetGroup(rs.Primary.ID)
		if err != nil {
			return nil
		}

		return nil
	}
}

func testAccCheckDatabricksGroupsDestroy(s *terraform.State) error {
	fmt.Println("[DEBUG] Running destroy test")
	endpoint := testAccProvider.Meta().(*Client).scim

	groupId := s.RootModule().Resources["databricks_group.group"].Primary.ID

	_, err := endpoint.GetGroup(groupId)

	if err == nil {
		return errors.New("cluster still exists")
	}

	if !resourceDatabricksClusterNotExistsError(err) {
		return err
	}

	return nil
}

func testAccDatabricksScimGroupConfigUpdate() string {
	return `
resource "databricks_group" "group" {
	display_name            = "TFGroupName"
} 
`
}

func testAccDatabricksScimGroup() string {
	return `
resource "databricks_group" "group" {
	display_name            = "TFUpdatedGroupName"
}
`
}
