// +build acceptance containerinfra

package v1

import (
	"strconv"
	"testing"

	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/certificates"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestClusterCRUDOperations(t *testing.T) {
	// Create a cluster
	client, err := clients.NewContainerInfraV1Client()
	c, err := clusters.Create(client, clusters.CreateOpts{ClusterTemplateID: "k8s"}).Extract()
	th.AssertNoErr(t, err)
	defer clusters.Delete(client, c.ID)
	th.AssertEquals(t, "CREATE_IN_PROGRESS", c.Status)
	th.AssertEquals(t, 1, c.Masters)
	th.AssertEquals(t, 1, c.Nodes)
	clusterID := c.ID
	clusterName := c.Name

	// List clusters
	pager := clusters.List(client, clusters.ListOpts{Limit: 1})
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		t.Logf("--- Page ---")

		clusterList, err := clusters.ExtractClusters(page)
		th.AssertNoErr(t, err)

		for _, c := range clusterList {
			t.Logf("Cluster: ID [%s] Name [%s] Status [%s] Nodes [%s]",
				c.ID, c.Name, c.Status, strconv.Itoa(c.Nodes))
		}

		return true, nil
	})
	th.CheckNoErr(t, err)

	// Get a cluster
	if clusterID == "" {
		t.Fatalf("In order to retrieve a cluster, the ClusterID must be set")
	}
	c, err = clusters.Get(client, clusterID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, clusterName, c.Name)
	th.AssertEquals(t, 1, c.Masters)
	th.AssertEquals(t, 1, c.Nodes)

	// Generate cluster credentials bundle
	c, err = waitForStatus(client, c, "CREATE_COMPLETE")
	th.AssertNoErr(t, err)
	bundle, err := certificates.CreateCredentialsBundle(client, clusterID)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, clusterID, bundle.ClusterID)
	th.AssertEquals(t, c.COEEndpoint, bundle.COEEndpoint)
	th.AssertEquals(t, true, bundle.PrivateKey.Bytes != nil)
	th.AssertEquals(t, true, bundle.Certificate.Bytes != nil)
	th.AssertEquals(t, true, bundle.CACertificate.Bytes != nil)
}

func waitForStatus(client *gophercloud.ServiceClient, cluster *clusters.Cluster, status string) (latest *clusters.Cluster, err error) {
	err = tools.WaitFor(func() (bool, error) {
		latest, err = clusters.Get(client, cluster.ID).Extract()
		if err != nil {
			return false, err
		}

		if latest.Status == status {
			// Success!
			return true, nil
		}

		if strings.HasSuffix(latest.Status, "FAILED") {
			return false, fmt.Errorf("The cluster is in the failed status. %s", latest.StatusReason)
		}

		return false, nil
	})
	return latest, err
}
