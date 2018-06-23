package clusters

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts allows the sorting of paginated collections through
// the API. SortKey allows you to sort by a particular cluster attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for pagination.
type ListOpts struct {
	Marker  string `q:"marker"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// clusters. It accepts a ListOptsBuilder, which allows you to sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific cluster based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Get(getURL(c, id), &r.Body, ro)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts satisfies the CreateOptsBuilder interface
type CreateOpts struct {
	ClusterTemplateID string `json:"cluster_template_id" required:"true"`
	Name              string `json:"name,omitempty"`
	Masters           int    `json:"master_count,omitempty"`
	Nodes             int    `json:"node_count,omitempty"`
	DiscoveryURL      string `json:"discovery_url,omitempty"`
}

// ToClusterCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new cluster using the values
// provided.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Post(createURL(c), b, &r.Body, ro)
	return
}

// Delete accepts a unique ID and deletes the cluster associated with it.
func Delete(c *gophercloud.ServiceClient, clusterID string) (r DeleteResult) {
	ro := &gophercloud.RequestOpts{ErrorContext: &common.ErrorResponse{}}
	_, r.Err = c.Delete(deleteURL(c, clusterID), ro)
	return
}
