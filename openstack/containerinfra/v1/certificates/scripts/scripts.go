package scripts

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clusters"
)

func Generate(coe string, cluster *clusters.Cluster) map[string][]byte {
	w, err := getScriptWriter(coe)
	if err != nil {
		return nil
	}
	return w.Generate(cluster)
}

type scriptWriter interface {
	Generate(cluster *clusters.Cluster) map[string][]byte
}

func getScriptWriter(coe string) (scriptWriter, error) {
	switch coe {
	case "swarm":
		return &swarmWriter{}, nil
	case "kubernetes":
		return &kubernetesWriter{}, nil
	default:
		return nil, fmt.Errorf("Unsupported COE: %s", coe)
	}
}
