package models

type CreateLoadBalancerRequest struct {
	LoadBalancerName string `json:"load_balancer_name"`
	EnvironmentID    string `json:"environment_id"`
	IsPublic         bool   `json:"is_public"`
	Ports            []Port `json:"ports"`
}
