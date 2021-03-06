package config

import (
	"fmt"
	"os"
	"strings"
)

// IMPORTANT!
// The environment variables represented as constants here should
// always line up with the environment variables in setup/container_definitions.json
const (
	AWS_ACCOUNT_ID                  = "LAYER0_AWS_ACCOUNT_ID"
	AWS_ACCESS_KEY_ID               = "LAYER0_AWS_ACCESS_KEY_ID"
	AWS_SECRET_ACCESS_KEY           = "LAYER0_AWS_SECRET_ACCESS_KEY"
	AWS_VPC_ID                      = "LAYER0_AWS_VPC_ID"
	AWS_PRIVATE_SUBNETS             = "LAYER0_AWS_PRIVATE_SUBNETS"
	AWS_PUBLIC_SUBNETS              = "LAYER0_AWS_PUBLIC_SUBNETS"
	AWS_ECS_ROLE                    = "LAYER0_AWS_ECS_ROLE"
	AWS_KEY_PAIR                    = "LAYER0_AWS_KEY_PAIR"
	AWS_S3_BUCKET                   = "LAYER0_AWS_S3_BUCKET"
	AWS_ECS_AGENT_SECURITY_GROUP_ID = "LAYER0_AWS_ECS_AGENT_SECURITY_GROUP_ID"
	AWS_ECS_INSTANCE_PROFILE        = "LAYER0_AWS_ECS_INSTANCE_PROFILE"
	JOB_ID                          = "LAYER0_JOB_ID"
	MYSQL_CONNECTION                = "LAYER0_MYSQL_CONNECTION"
	MYSQL_ADMIN_CONNECTION          = "LAYER0_MYSQL_ADMIN_CONNECTION"
	AWS_SERVICE_AMI                 = "LAYER0_AWS_SERVICE_AMI"
	AWS_REGION                      = "LAYER0_AWS_REGION"
	API_AUTH_TOKEN                  = "LAYER0_API_AUTH_TOKEN"
	API_ENDPOINT                    = "LAYER0_API_ENDPOINT"
	API_PORT                        = "LAYER0_API_PORT"
	API_LOG_LEVEL                   = "LAYER0_API_LOG_LEVEL"
	CLI_AUTH                        = "LAYER0_CLI_AUTH"
	PREFIX                          = "LAYER0_PREFIX"
	RUNNER_LOG_LEVEL                = "LAYER0_RUNNER_LOG_LEVEL"
	RUNNER_VERSION_TAG              = "LAYER0_RUNNER_VERSION_TAG"
	SKIP_SSL_VERIFY                 = "LAYER0_SKIP_SSL_VERIFY"
	SKIP_VERSION_VERIFY             = "LAYER0_SKIP_VERSION_VERIFY"
	SQLITE_DB_PATH                  = "LAYER0_SQLITE_DB_PATH"
)

// non environment variable constants
const (
	API_CERTIFICATE_ID     = "api"
	API_CERTIFICATE_NAME   = "api"
	API_ENVIRONMENT_ID     = "api"
	API_ENVIRONMENT_NAME   = "api"
	API_LOAD_BALANCER_ID   = "api"
	API_LOAD_BALANCER_NAME = "api"
	API_SERVICE_ID         = "api"
	API_SERVICE_NAME       = "api"
)

var RequiredAPIVariables = []string{
	AWS_ACCOUNT_ID,
	AWS_ACCESS_KEY_ID,
	AWS_SECRET_ACCESS_KEY,
	AWS_VPC_ID,
	AWS_PRIVATE_SUBNETS,
	AWS_PUBLIC_SUBNETS,
	AWS_ECS_ROLE,
	AWS_KEY_PAIR,
	AWS_S3_BUCKET,
	AWS_ECS_AGENT_SECURITY_GROUP_ID,
	AWS_ECS_INSTANCE_PROFILE,
	MYSQL_CONNECTION,
	MYSQL_ADMIN_CONNECTION,
}

var RequiredCLIVariables = []string{}

var RequiredRunnerVariables = []string{
	AWS_ACCESS_KEY_ID,
	AWS_SECRET_ACCESS_KEY,
	AWS_VPC_ID,
	AWS_PRIVATE_SUBNETS,
	AWS_PUBLIC_SUBNETS,
	MYSQL_CONNECTION,
}

func Validate(required []string) error {
	for _, key := range required {
		if os.Getenv(key) == "" {
			return fmt.Errorf("Required environment variable '%s' not set", key)
		}
	}

	return nil
}

func get(key string) string {
	return os.Getenv(key)
}

func getOr(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

var apiVersion string

func SetAPIVersion(version string) {
	apiVersion = version
}

func APIVersion() string {
	return apiVersion
}

var cliVersion string

func CLIVersion() string {
	return cliVersion
}

func SetCLIVersion(version string) {
	cliVersion = version
}

var serviceAMI = map[string]string{
	"us-west-2": "ami-6cb9ac0d",
	"us-east-1": "ami-804130ea",
	"eu-west-1": "ami-e563bf96",
}

func AWSLogGroupID() string {
	return fmt.Sprintf("l0-%s", Prefix())
}

func AWSServiceAMI() string {
	if ami := get(AWS_SERVICE_AMI); ami != "" {
		return ami
	}

	return serviceAMI[AWSRegion()]
}

func AWSAccountID() string {
	return get(AWS_ACCOUNT_ID)
}

func AWSAccessKey() string {
	return get(AWS_ACCESS_KEY_ID)
}

func AWSSecretKey() string {
	return get(AWS_SECRET_ACCESS_KEY)
}

func AWSRegion() string {
	return getOr(AWS_REGION, "us-west-2")
}

func AWSVPCID() string {
	return get(AWS_VPC_ID)
}

func AWSPrivateSubnets() string {
	return get(AWS_PRIVATE_SUBNETS)
}

func AWSPublicSubnets() string {
	return get(AWS_PUBLIC_SUBNETS)
}

func AWSECSRole() string {
	return get(AWS_ECS_ROLE)
}

func AWSKeyPair() string {
	return get(AWS_KEY_PAIR)
}

func AWSS3Bucket() string {
	return get(AWS_S3_BUCKET)
}

func APIAuthToken() string {
	// usr:pwd = layer0:nohaxplz, base64 encoded (basic http auth)
	return getOr(API_AUTH_TOKEN, "bGF5ZXIwOm5vaGF4cGx6")
}

func APIEndpoint() string {
	return getOr(API_ENDPOINT, "http://localhost:9090/")
}

func APIPort() string {
	return getOr(API_PORT, "9090")
}

func APILogLevel() string {
	return getOr(API_LOG_LEVEL, "1")
}

func CLIAuth() string {
	// usr:pwd = layer0:nohaxplz, base64 encoded (basic http auth
	token := getOr(CLI_AUTH, "bGF5ZXIwOm5vaGF4cGx6")
	return fmt.Sprintf("Basic %v", token)
}

func MySQLConnection() string {
	return get(MYSQL_CONNECTION)
}

func MySQLAdminConnection() string {
	return get(MYSQL_ADMIN_CONNECTION)
}

func Prefix() string {
	return getOr(PREFIX, "l0")
}

func RunnerLogLevel() string {
	return getOr(RUNNER_LOG_LEVEL, "1")
}

func RunnerVersionTag() string {
	return getOr(RUNNER_VERSION_TAG, "latest")
}

func SQLiteDbPath() string {
	return getOr(SQLITE_DB_PATH, "")
}

func AWSAgentGroupID() string {
	return get(AWS_ECS_AGENT_SECURITY_GROUP_ID)
}

func AWSECSInstanceProfile() string {
	return get(AWS_ECS_INSTANCE_PROFILE)
}

func ShouldVerifySSL() bool {
	val := strings.ToLower(getOr(SKIP_SSL_VERIFY, ""))
	if val == "1" || val == "true" {
		return false
	}

	return true
}

func ShouldVerifyVersion() bool {
	val := strings.ToLower(getOr(SKIP_VERSION_VERIFY, ""))
	if val == "1" || val == "true" {
		return false
	}

	return true
}
