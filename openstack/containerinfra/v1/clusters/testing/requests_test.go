package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
	fake "github.com/gophercloud/gophercloud/openstack/containerinfra/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "clusters": [
    {
      "status": "CREATE_COMPLETE",
      "uuid": "a56a6cd8-0779-461b-b1eb-26cec904284a",
      "links": [
        {
          "href": "http://65.61.151.130:9511/v1/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a",
          "rel": "self"
        },
        {
          "href": "http://65.61.151.130:9511/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a",
          "rel": "bookmark"
        }
      ],
      "stack_id": "f8ef771f-1ffa-4ad5-99b8-651bf7669f80",
      "master_count": 1,
      "clustertemplate_id": "5b793604-fc76-4886-a834-ed522812cdcb",
      "node_count": 1,
      "cluster_create_timeout": 0,
      "name": "k8scluster"
    }
  ]
}
			`)
	})

	client := fake.ServiceClient()
	count := 0

	results := clusters.List(client, clusters.ListOpts{})

	err := results.EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := clusters.ExtractClusters(page)
		if err != nil {
			t.Errorf("Failed to extract clusters: %v", err)
			return false, err
		}

		expected := []clusters.Cluster{
			{
				Status:            "CREATE_COMPLETE",
				Name:              "k8scluster",
				ID:                "a56a6cd8-0779-461b-b1eb-26cec904284a",
				Masters:           1,
				Nodes:             1,
				ClusterTemplateID: "5b793604-fc76-4886-a834-ed522812cdcb",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "status": "CREATE_COMPLETE",
  "container_version": "1.9.1",
  "uuid": "a56a6cd8-0779-461b-b1eb-26cec904284a",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a",
      "rel": "bookmark"
    }
  ],
  "stack_id": "f8ef771f-1ffa-4ad5-99b8-651bf7669f80",
  "created_at": "2016-07-14T23:58:50+00:00",
  "api_address": "https://172.29.248.18:6443",
  "discovery_url": "https://discovery.etcd.io/ac7f669ebe467d061c59bfe5b6a5f6fe",
  "updated_at": "2016-07-15T00:02:53+00:00",
  "master_count": 1,
  "clustertemplate_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "master_addresses": [
    "172.29.248.18"
  ],
  "node_count": 1,
  "node_addresses": [
    "172.29.248.19"
  ],
  "status_reason": "Stack CREATE completed successfully",
  "cluster_create_timeout": 0,
  "name": "k8scluster"
}
			`)
	})

	b, err := clusters.Get(fake.ServiceClient(), "a56a6cd8-0779-461b-b1eb-26cec904284a").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "CREATE_COMPLETE", b.Status)
	th.AssertEquals(t, "1.9.1", b.ContainerVersion)
	th.AssertEquals(t, "k8scluster", b.Name)
	th.AssertEquals(t, "5b793604-fc76-4886-a834-ed522812cdcb", b.ClusterTemplateID)
	th.AssertEquals(t, 1, b.Nodes)
	th.AssertEquals(t, "a56a6cd8-0779-461b-b1eb-26cec904284a", b.ID)
	th.AssertEquals(t, "https://172.29.248.18:6443", b.COEEndpoint)
	th.AssertEquals(t, "172.29.248.18", b.MasterAddresses[0])
	th.AssertEquals(t, "172.29.248.19", b.NodeAddresses[0])
}

func TestGetFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters/duplicatename", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 409,
      "code": "client",
      "links": [],
      "title": "Multiple clusters exist with same name",
      "detail": "Multiple clusters exist with same name. Please use the cluster uuid instead.",
      "request_id": ""
    }
  ]
}
		`)
	})

	res := clusters.Get(fake.ServiceClient(), "duplicatename")

	th.AssertEquals(t, "Multiple clusters exist with same name. Please use the cluster uuid instead.", res.Err.Error())

	er, ok := res.Err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusConflict, er.Actual)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
  "node_count": 2,
  "clustertemplate_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "name": "mycluster"
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
  "status": "CREATE_IN_PROGRESS",
  "container_version": "1.9.1",
  "uuid": "39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
  "links": [
    {
      "href": "http://65.61.151.130:9511/v1/clusters/39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
      "rel": "self"
    },
    {
      "href": "http://65.61.151.130:9511/clusters/39109e8a-516e-41a4-8b1d-22e9a56e4aa2",
      "rel": "bookmark"
    }
  ],
  "stack_id": "f27f1581-fb2e-4033-af93-d2cf19bb8462",
  "created_at": "2016-08-08T16:45:18+00:00",
  "api_address": null,
  "discovery_url": "https://discovery.etcd.io/32ccf12b42e75b6822ac18c2c0391e5f",
  "updated_at": null,
  "master_count": 1,
  "clustertemplate_id": "5b793604-fc76-4886-a834-ed522812cdcb",
  "master_addresses": null,
  "node_count": 2,
  "node_addresses": null,
  "status_reason": null,
  "cluster_create_timeout": 60,
  "name": "mycluster"
}
		`)
	})

	options := clusters.CreateOpts{Name: "mycluster", Nodes: 2, ClusterTemplateID: "5b793604-fc76-4886-a834-ed522812cdcb"}
	b, err := clusters.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "CREATE_IN_PROGRESS", b.Status)
	th.AssertEquals(t, "mycluster", b.Name)
	th.AssertEquals(t, "5b793604-fc76-4886-a834-ed522812cdcb", b.ClusterTemplateID)
	th.AssertEquals(t, "39109e8a-516e-41a4-8b1d-22e9a56e4aa2", b.ID)
	th.AssertEquals(t, 1, b.Masters)
	th.AssertEquals(t, 2, b.Nodes)
}

func TestCreateFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 500,
      "code": "client",
      "links": [],
      "title": "Nova is down",
      "detail": "Nova is down. Try again later.",
      "request_id": ""
    }
  ]
}
		`)
	})

	options := clusters.CreateOpts{Name: "mycluster", Nodes: 2, ClusterTemplateID: "5b793604-fc76-4886-a834-ed522812cdcb"}

	res := clusters.Create(fake.ServiceClient(), options)

	th.AssertEquals(t, "Nova is down. Try again later.", res.Err.Error())

	er, ok := res.Err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusInternalServerError, er.Actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := clusters.Delete(fake.ServiceClient(), "a56a6cd8-0779-461b-b1eb-26cec904284a")
	th.AssertNoErr(t, res.Err)
}

func TestDeleteFailed(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/clusters/a56a6cd8-0779-461b-b1eb-26cec904284a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `
{
  "errors": [
    {
      "status": 400,
      "code": "client",
      "links": [],
      "title": "Cluster k8scluster already has an operation in progress",
      "detail": "Cluster k8scluster already has an operation in progress.",
      "request_id": ""
    }
  ]
}
		`)
	})

	res := clusters.Delete(fake.ServiceClient(), "a56a6cd8-0779-461b-b1eb-26cec904284a")

	th.AssertEquals(t, "Cluster k8scluster already has an operation in progress.", res.Err.Error())

	er, ok := res.Err.(*fake.ErrorResponse)
	th.AssertEquals(t, true, ok)
	th.AssertEquals(t, http.StatusBadRequest, er.Actual)
}
