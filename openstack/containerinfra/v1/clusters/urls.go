package clusters

import "github.com/gophercloud/gophercloud"

func createURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

func listURL(c *gophercloud.ServiceClient) string {
	return c.ServiceURL("clusters")
}

func getURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}

func deleteURL(c *gophercloud.ServiceClient, id string) string {
	return c.ServiceURL("clusters", id)
}
