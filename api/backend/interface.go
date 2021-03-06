package backend

import (
	"github.com/quintilesims/layer0/common/models"
)

type Backend interface {
	CreateEnvironment(environmentName, instanceSize string, minClusterCount int, userData []byte) (*models.Environment, error)
	UpdateEnvironment(environmentID string, minClusterCount int) (*models.Environment, error)
	DeleteEnvironment(environmentID string) error
	GetEnvironment(environmentID string) (*models.Environment, error)
	ListEnvironments() ([]*models.Environment, error)

	ListDeploys() ([]*models.Deploy, error)
	GetDeploy(deployID string) (*models.Deploy, error)
	CreateDeploy(name string, body []byte) (*models.Deploy, error)
	DeleteDeploy(deployID string) error

	ListServices() ([]*models.Service, error)
	GetService(envID, serviceID string) (*models.Service, error)
	CreateService(serviceName, environmentID, deployID, loadBalancerID string) (*models.Service, error)
	DeleteService(environmentID, serviceID string) error
	ScaleService(environmentID, serviceID string, count int) (*models.Service, error)
	UpdateService(environmentID, serviceID, deployID string) (*models.Service, error)
	GetServiceLogs(environmentID, serviceID string, tail int) ([]*models.LogFile, error)

	CreateTask(envID, taskName, deployVersion string, copies int, overrides []models.ContainerOverride) (*models.Task, error)
	ListTasks() ([]*models.Task, error)
	GetTask(envID, taskID string) (*models.Task, error)
	DeleteTask(envID, taskID string) error
	GetTaskLogs(environmentID, taskID string, tail int) ([]*models.LogFile, error)

	CreateCertificate(name, public, private, chain string) (*models.Certificate, error)
	ListCertificates() ([]*models.Certificate, error)
	GetCertificate(certificateID string) (*models.Certificate, error)
	DeleteCertificate(certificateID string) error

	ListLoadBalancers() ([]*models.LoadBalancer, error)
	GetLoadBalancer(id string) (*models.LoadBalancer, error)
	DeleteLoadBalancer(id string) error
	CreateLoadBalancer(loadBalancerName, environmentID string, isPublic bool, ports []models.Port) (*models.LoadBalancer, error)
	UpdateLoadBalancer(loadBalancerID string, ports []models.Port) (*models.LoadBalancer, error)

	StartRightSizer()
	GetRightSizerHealth() (string, error)
}
