package types

import (
	"time"
)

// ManagedClusterInfo holds details about a managed (imported) cluster.
type ManagedClusterInfo struct {
	Name         string            `json:"name"`
	Labels       map[string]string `json:"labels"`
	CreationTime time.Time         `json:"creationTime"`
	Context      string            `json:"context,omitempty"`
}

// ContextInfo holds basic info for a kubeconfig context.
type ContextInfo struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}

// ClusterDetails holds detailed information about a cluster.
type ClusterDetails struct {
	ClusterName        string           `json:"clusterName"`
	Contexts           []ContextInfo    `json:"contexts"`
	ITSManagedClusters []ManagedClusterInfo `json:"itsManagedClusters"`
}
