package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a cluster resource.
func (r commonResult) Extract() (*Cluster, error) {
	var s *Cluster
	err := r.ExtractInto(&s)
	return s, err
}

// CreateResult temporarily contains the response from a Create call.
type CreateResult struct {
	commonResult
}

// GetResult represents the result of a get operation.
type GetResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Represents a Container Orchestration Engine Cluster, i.e. a cluster
type Cluster struct {
	// UUID for the cluster
	ID string `json:"uuid"`

	// Human-readable name for the cluster. Might not be unique.
	Name string `json:"name"`

	// Indicates whether cluster is currently operational. Possible values include:
	// CREATE_IN_PROGRESS, CREATE_FAILED, CREATE_COMPLETE, UPDATE_IN_PROGRESS, UPDATE_FAILED, UPDATE_COMPLETE,
	// DELETE_IN_PROGRESS, DELETE_FAILED, DELETE_COMPLETE, RESUME_COMPLETE, RESTORE_COMPLETE, ROLLBACK_COMPLETE,
	// SNAPSHOT_COMPLETE, CHECK_COMPLETE, ADOPT_COMPLETE.
	Status string `json:"status"`

	// Additional information on the cluster status, such as why the cluster is in a failed state.
	StatusReason string `json:"status_reason"`

	// The number of master nodes in the cluster.
	Masters int `json:"master_count"`

	// The number of host nodes in the cluster.
	Nodes int `json:"node_count"`

	// The UUID of the clustertemplate used to generate the cluster.
	ClusterTemplateID string `json:"clustertemplate_id"`

	// The URL of the COE API.
	COEEndpoint string `json:"api_address"`

	// The IP addresses of the master nodes.
	MasterAddresses []string `json:"master_addresses"`

	// The IP addresses of the host nodes.
	NodeAddresses []string `json:"node_addresses"`

	// The version of the Docker client compatible with the cluster.
	ContainerVersion string `json:"container_version"`
}

// ClusterPage is the page returned by a pager when traversing over a
// collection of clusters.
type ClusterPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of clusters has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (r ClusterPage) NextPageURL() (string, error) {
	var s struct {
		Next string `json:"next"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Next, nil
}

// IsEmpty checks whether a ClusterPage struct is empty.
func (r ClusterPage) IsEmpty() (bool, error) {
	is, err := ExtractClusters(r)
	return len(is) == 0, err
}

// ExtractClusters accepts a Page struct, specifically a ClusterPage struct,
// and extracts the elements into a slice of Cluster structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractClusters(r pagination.Page) ([]Cluster, error) {
	var s struct {
		Clusters []Cluster `json:"clusters"`
	}
	err := (r.(ClusterPage)).ExtractInto(&s)
	return s.Clusters, err
}
