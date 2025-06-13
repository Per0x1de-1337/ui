package common

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	// "github.com/kubestellar/ui/its/manual/handlers"

	// "sync"
	// "github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ContextInfo holds basic info for a kubeconfig context.
type ContextInfo struct {
	Name    string `json:"name"`
	Cluster string `json:"cluster"`
}

// ManagedClusterInfo holds details about a managed (imported) cluster.
type ManagedClusterInfo struct {
	Name         string            `json:"name"`
	Labels       map[string]string `json:"labels"`
	CreationTime time.Time         `json:"creationTime"`
	Context      string            `json:"context,omitempty"`
}

// ManagedClusterStatus represents the status of a managed cluster
type ManagedClusterStatus struct {
	Conditions []ManagedClusterCondition `json:"conditions,omitempty"`
	Version    map[string]string         `json:"version,omitempty"`
	Capacity   map[string]string         `json:"capacity,omitempty"`
}

// ManagedClusterCondition represents a condition of a managed cluster
type ManagedClusterCondition struct {
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	Message            string      `json:"message,omitempty"`
	Reason             string      `json:"reason,omitempty"`
	Status             string      `json:"status,omitempty"`
	Type               string      `json:"type,omitempty"`
}

type KubeInfoResult struct {
	Contexts       []string
	Clusters       []string
	CurrentContext string
	ItsData        interface{}
	Err            error
}

func GetKubeInfo() ([]ContextInfo, []string, string, []ManagedClusterInfo, error) {
	kubeconfig := kubeconfigPath()
	// Log which kubeconfig is being used.
	if os.Getenv("KUBECONFIG") == "" {
		log.Printf("Using default kubeconfig path: %s", kubeconfig)
	} else {
		log.Printf("Using kubeconfig from environment: %s", kubeconfig)
	}

	config, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, nil, "", nil, err
	}

	var contexts []ContextInfo
	clusterSet := make(map[string]bool)
	currentContext := config.CurrentContext
	var managedClusters []ManagedClusterInfo

	// Process ITS contexts (e.g., contexts starting with "its")
	for contextName := range config.Contexts {
		if strings.HasPrefix(contextName, "its") {
			log.Printf("Processing ITS context: %s", contextName)
			clientConfig := clientcmd.NewNonInteractiveClientConfig(
				*config,
				contextName,
				&clientcmd.ConfigOverrides{},
				clientcmd.NewDefaultClientConfigLoadingRules(),
			)
			restConfig, err := clientConfig.ClientConfig()
			if err != nil {
				log.Printf("Error creating REST config for context %s: %v", contextName, err)
				continue
			}
			clientset, err := kubernetes.NewForConfig(restConfig)
			if err != nil {
				log.Printf("Error creating clientset for context %s: %v", contextName, err)
				continue
			}
			clustersBytes, err := clientset.RESTClient().Get().
				AbsPath("/apis/cluster.open-cluster-management.io/v1").
				Resource("managedclusters").
				DoRaw(context.TODO())
			if err != nil {
				log.Printf("Error fetching managed clusters from context %s: %v", contextName, err)
				continue
			}
			var clusterList struct {
				Items []struct {
					Metadata struct {
						Name              string            `json:"name"`
						Labels            map[string]string `json:"labels"`
						CreationTimestamp string            `json:"creationTimestamp"`
					} `json:"metadata"`
				} `json:"items"`
			}
			if err := json.Unmarshal(clustersBytes, &clusterList); err != nil {
				log.Printf("Error unmarshaling clusters from context %s: %v", contextName, err)
				continue
			}
			for _, item := range clusterList.Items {
				creationTime, _ := time.Parse(time.RFC3339, item.Metadata.CreationTimestamp)
				managedClusters = append(managedClusters, ManagedClusterInfo{
					Name:         item.Metadata.Name,
					Labels:       item.Metadata.Labels,
					CreationTime: creationTime,
					Context:      contextName,
				})
			}
		}
	}

	// Process contexts with "-kubeflex" suffix.
	for contextName, ctx := range config.Contexts {
		if strings.HasSuffix(contextName, "-kubeflex") {
			contexts = append(contexts, ContextInfo{
				Name:    contextName,
				Cluster: ctx.Cluster,
			})
			clusterSet[ctx.Cluster] = true
		}
	}

	var clusters []string
	for clusterName := range clusterSet {
		clusters = append(clusters, clusterName)
	}

	return contexts, clusters, currentContext, managedClusters, nil
}

// func ImportClusterHandler(c *gin.Context) {
// 	file, err := c.FormFile("kubeconfig")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "kubeconfig file is required"})
// 		return
// 	}

// 	src, err := file.Open()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
// 		return
// 	}
// 	defer src.Close()

// 	data, err := io.ReadAll(src)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file contents"})
// 		return
// 	}

// 	// 2. Load kubeconfig
// 	cfg, err := clientcmd.Load(data)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid kubeconfig format"})
// 		return
// 	}

// 	adjustClusterServerEndpoints(cfg)

// 	tmpPath := filepath.Join(os.TempDir(), fmt.Sprintf("import-%d.kubeconfig", time.Now().UnixNano()))
// 	outData, err := clientcmd.Write(*cfg)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize kubeconfig"})
// 		return
// 	}
// 	if err := os.WriteFile(tmpPath, outData, 0600); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to write temp kubeconfig"})
// 		return
// 	}

// 	hubKube := os.Getenv("HUB_BOOTSTRAP_KUBECONFIG")
// 	if hubKube == "" {
// 		hubKube = filepath.Join(os.Getenv("HOME"), ".kube", "bootstrap.kubeconfig")
// 	}

// 	exec.Command("helm", "repo", "add", "ocm", "https://open-cluster-management.io/helm-charts").Run()
// 	exec.Command("helm", "repo", "update").Run()

// 	releaseName := fmt.Sprintf("klusterlet-%s", strings.ReplaceAll(cfg.CurrentContext, "_", "-"))

// 	cmd := exec.Command("helm", "upgrade", "--install",
// 		releaseName, "ocm/klusterlet",
// 		"--kubeconfig", tmpPath,
// 		"--namespace", "open-cluster-management",
// 		"--create-namespace",
// 		"--set", fmt.Sprintf("hubKubeconfig=%s", hubKube),
// 	)

// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("helm install failed: %s", string(output))})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message":     "Cluster import initiated",
// 		"release":     releaseName,
// 		"helm_output": string(output),
// 	})
// }

// func adjustClusterServerEndpoints(config *clientcmdapi.Config) {
// 	for name, cluster := range config.Clusters {

// 		if strings.Contains(cluster.Server, "localhost") {
// 			cluster.Server = strings.Replace(cluster.Server, "localhost", fmt.Sprintf("%s", name), 1)
// 		}
// 	}
// }

// kubeconfigPath returns the path to the kubeconfig file
func kubeconfigPath() string {
	if path := os.Getenv("KUBECONFIG"); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Unable to get user home directory: %v", err)
	}
	return fmt.Sprintf("%s/.kube/config", home)
}
