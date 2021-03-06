package ecsbackend

import (
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/golang/mock/gomock"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/api/backend/mock_backend"
	"github.com/quintilesims/layer0/common/aws/autoscaling"
	"github.com/quintilesims/layer0/common/aws/autoscaling/mock_autoscaling"
	"github.com/quintilesims/layer0/common/aws/ec2"
	"github.com/quintilesims/layer0/common/aws/ec2/mock_ec2"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/aws/ecs/mock_ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"github.com/quintilesims/layer0/common/testutils"
	"testing"
)

type MockECSEnvironmentManager struct {
	EC2         *mock_ec2.MockProvider
	ECS         *mock_ecs.MockProvider
	AutoScaling *mock_autoscaling.MockProvider
	Backend     *mock_backend.MockBackend
}

func NewMockECSEnvironmentManager(ctrl *gomock.Controller) *MockECSEnvironmentManager {
	return &MockECSEnvironmentManager{
		EC2:         mock_ec2.NewMockProvider(ctrl),
		ECS:         mock_ecs.NewMockProvider(ctrl),
		AutoScaling: mock_autoscaling.NewMockProvider(ctrl),
		Backend:     mock_backend.NewMockBackend(ctrl),
	}
}

func (this *MockECSEnvironmentManager) Environment() *ECSEnvironmentManager {
	manager := NewECSEnvironmentManager(this.ECS, this.EC2, this.AutoScaling, this.Backend)
	manager.Clock = &testutils.StubClock{}
	return manager
}

func TestGetEnvironment(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)
				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
				clusterName := ecsEnvironmentID.String()
				securityGroupName := ecsEnvironmentID.SecurityGroupName()

				mockEnvironment.ECS.EXPECT().
					DescribeCluster(clusterName).
					Return(ecs.NewCluster(clusterName), nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(autoScalingGroupName).
					Return(asg, nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(clusterName).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(securityGroupName).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.GetEnvironment("envid")
			},
		},
		{
			Name: "Should return layer0-formatted environment id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)
				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterName := ecsEnvironmentID.String()

				mockEnvironment.ECS.EXPECT().
					DescribeCluster(gomock.Any()).
					Return(ecs.NewCluster(clusterName), nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(gomock.Any()).
					Return(asg, nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(gomock.Any()).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				environment, err := manager.GetEnvironment("envid")
				if err != nil {
					reporter.Fatal(err)
				}
				reporter.AssertEqual("envid", environment.EnvironmentID)
			},
		},
		{
			Name: "Should propagate ecs.DescribeCluster error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				mockEnvironment.ECS.EXPECT().
					DescribeCluster(gomock.Any()).
					Return(nil, fmt.Errorf("some_error"))

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				if _, err := manager.GetEnvironment("envid"); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestListEnvironments(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should return layer0-formatted environment ids and use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
				clusterName := ecsEnvironmentID.String()
				securityGroupName := ecsEnvironmentID.SecurityGroupName()

				mockEnvironment.ECS.EXPECT().
					Helper_DescribeClusters().
					Return([]*ecs.Cluster{ecs.NewCluster(clusterName)}, nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(autoScalingGroupName).
					Return(asg, nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(clusterName).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(securityGroupName).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				environments, err := manager.ListEnvironments()
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual(len(environments), 1)
				reporter.AssertEqual(environments[0].EnvironmentID, "envid")
			},
		},
		{
			Name: "Should propagate ecs.Helper_DescribeClusters error",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				mockEnvironment.ECS.EXPECT().
					Helper_DescribeClusters().
					Return(nil, fmt.Errorf("some_error"))

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				if _, err := manager.ListEnvironments(); err == nil {
					reporter.Errorf("Error was nil!")
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestDeleteEnvironment(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name: "Should delete autoscaling group, launch configuration, security group, and cluster correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
				//autoScalingGroup := autoscaling.NewGroup()
				launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()
				securityGroupName := ecsEnvironmentID.SecurityGroupName()
				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				clusterName := ecsEnvironmentID.String()

				mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMinSize(autoScalingGroupName, 0).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMaxSize(autoScalingGroupName, 0).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					DeleteAutoScalingGroup(&autoScalingGroupName).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					DeleteLaunchConfiguration(&launchConfigurationName).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(autoScalingGroupName).
					Return(nil, awserr.New("GroupNotFoundException", "group not found", nil))

				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(securityGroupName).
					Return(securityGroup, nil)

				mockEnvironment.EC2.EXPECT().
					DeleteSecurityGroup(securityGroup).
					Return(nil)

				mockEnvironment.ECS.EXPECT().
					DeleteCluster(clusterName).
					Return(nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.DeleteEnvironment("envid")
			},
		},
		{
			Name: "Should pass through idempotent aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMinSize(gomock.Any(), gomock.Any()).
					Return(awserr.New("ValidationError", "name not found", nil))

				mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMaxSize(gomock.Any(), gomock.Any()).
					Return(awserr.New("ValidationError", "name not found", nil))

				mockEnvironment.AutoScaling.EXPECT().
					DeleteAutoScalingGroup(gomock.Any()).
					Return(awserr.New("ValidationError", "name not found", nil))

				mockEnvironment.AutoScaling.EXPECT().
					DeleteLaunchConfiguration(gomock.Any()).
					Return(awserr.New("ValidationError", "name not found", nil))

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(gomock.Any()).
					Return(nil, awserr.New("GroupNotFoundException", "group not found", nil))

				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(nil, nil)

				mockEnvironment.ECS.EXPECT().
					DeleteCluster(gomock.Any()).
					Return(awserr.New("ClusterNotFoundException", "name not found", nil))

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				if err := manager.DeleteEnvironment(""); err != nil {
					reporter.Fatal(err)
				}
			},
		},
		{
			Name: "Should propagate unexpected aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSEnvironmentManager {
					mockEnvironment := NewMockECSEnvironmentManager(ctrl)

					autoScalingGroup := autoscaling.NewGroup()
					securityGroup := ec2.NewSecurityGroup("some_sg_id")

					mockEnvironment.AutoScaling.EXPECT().
						UpdateAutoScalingGroupMinSize(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						UpdateAutoScalingGroupMaxSize(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						DeleteAutoScalingGroup(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						DeleteLaunchConfiguration(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						DescribeAutoScalingGroup(gomock.Any()).
						Return(autoScalingGroup, g.Error()).
						AnyTimes()

					mockEnvironment.EC2.EXPECT().
						DescribeSecurityGroup(gomock.Any()).
						Return(securityGroup, g.Error()).
						AnyTimes()

					mockEnvironment.EC2.EXPECT().
						DeleteSecurityGroup(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.ECS.EXPECT().
						DeleteCluster(gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					return mockEnvironment.Environment()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSEnvironmentManager)

				for i := 0; i < 8; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some eror"))

					manager := setup(g)
					if err := manager.DeleteEnvironment(""); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestCreateEnvironment(t *testing.T) {
	defer id.StubIDGeneration("envid")()

	testCases := []testutils.TestCase{
		{
			Name: "Should create autoscaling group, launch configuration, security group, and cluster correctly",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()
				launchConfigurationName := ecsEnvironmentID.LaunchConfigurationName()
				securityGroupName := ecsEnvironmentID.SecurityGroupName()
				securityGroupID := "some_sg_id"
				clusterName := ecsEnvironmentID.String()

				mockEnvironment.ECS.EXPECT().
					CreateCluster(clusterName).
					Return(ecs.NewCluster(clusterName), nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(clusterName).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				mockEnvironment.EC2.EXPECT().
					CreateSecurityGroup(securityGroupName, gomock.Any(), config.TEST_AWS_VPC_ID).
					Return(&securityGroupID, nil)

				mockEnvironment.EC2.EXPECT().
					AuthorizeSecurityGroupIngressFromGroup(&securityGroupID, &securityGroupID).
					Return(nil)

				var checkLaunchConfig = func(name, amiID, iamInstanceProfile, instanceType, keyName, userData *string, securityGroups []*string) error {
					agentSGID := config.TEST_AWS_ECS_AGENT_SECURITY_GROUP_ID

					reporter.AssertEqualf(launchConfigurationName, *name, "LaunchConfigurationName")
					reporter.AssertEqualf(config.TEST_AWS_SERVICE_AMI, *amiID, "AMI ID")
					reporter.AssertEqualf(config.TEST_AWS_ECS_INSTANCE_PROFILE, *iamInstanceProfile, "InstanceProfile")
					reporter.AssertEqualf("m3.medium", *instanceType, "Instance Type")
					reporter.AssertEqualf(config.TEST_AWS_KEY_PAIR, *keyName, "KeyPair")
					reporter.AssertEqualf(securityGroupID, *securityGroups[0], "SecurityGroupID 0")
					reporter.AssertEqualf(agentSGID, *securityGroups[1], "SecurityGroupID 1")

					return nil
				}

				mockEnvironment.AutoScaling.EXPECT().
					CreateLaunchConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Do(checkLaunchConfig)

				minCount := 2
				maxCount := 2
				subnets := config.TEST_AWS_PRIVATE_SUBNETS
				mockEnvironment.AutoScaling.EXPECT().
					CreateAutoScalingGroup(
						autoScalingGroupName,
						launchConfigurationName,
						subnets,
						minCount,
						maxCount).
					Return(nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)
				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(autoScalingGroupName).
					Return(asg, nil)

				securityGroup := ec2.NewSecurityGroup(securityGroupID)
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(securityGroupName).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.CreateEnvironment("env_name", "m3.medium", 2, nil)
			},
		},
		{
			Name: "Should render base64 encoded useDataTemplate if specified",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterName := ecsEnvironmentID.String()
				securityGroupID := "some_sg_id"

				mockEnvironment.ECS.EXPECT().
					CreateCluster(gomock.Any()).
					Return(ecs.NewCluster(clusterName), nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(gomock.Any()).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				mockEnvironment.EC2.EXPECT().
					CreateSecurityGroup(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&securityGroupID, nil)

				mockEnvironment.EC2.EXPECT().
					AuthorizeSecurityGroupIngressFromGroup(gomock.Any(), gomock.Any()).
					Return(nil)

				userData := base64.StdEncoding.EncodeToString([]byte("user data"))
				mockEnvironment.AutoScaling.EXPECT().
					CreateLaunchConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), &userData, gomock.Any()).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					CreateAutoScalingGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(gomock.Any()).
					Return(asg, nil)

				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.CreateEnvironment("env_name", "m3.medium", 0, []byte("user data"))
			},
		},
		{
			Name: "Should return layer0-formatted environment id",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
				clusterName := ecsEnvironmentID.String()
				securityGroupID := "some_sg_id"

				mockEnvironment.ECS.EXPECT().
					CreateCluster(gomock.Any()).
					Return(ecs.NewCluster(clusterName), nil)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeLaunchConfiguration(gomock.Any()).
					Return(autoscaling.NewLaunchConfiguration("m3.medium"), nil)

				mockEnvironment.EC2.EXPECT().
					CreateSecurityGroup(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(&securityGroupID, nil)

				mockEnvironment.EC2.EXPECT().
					AuthorizeSecurityGroupIngressFromGroup(gomock.Any(), gomock.Any()).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					CreateLaunchConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				mockEnvironment.AutoScaling.EXPECT().
					CreateAutoScalingGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				asg := autoscaling.NewGroup()
				asg.LaunchConfigurationName = stringp(clusterName)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(gomock.Any()).
					Return(asg, nil)

				securityGroup := ec2.NewSecurityGroup("some_sg_id")
				mockEnvironment.EC2.EXPECT().
					DescribeSecurityGroup(gomock.Any()).
					Return(securityGroup, nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)

				environment, err := manager.CreateEnvironment("env_name", "m3.medium", 0, nil)
				if err != nil {
					reporter.Fatal(err)
				}

				reporter.AssertEqual("envid", environment.EnvironmentID)
			},
		},
		{
			Name: "Should propagate unexpected aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) interface{} {
					mockEnvironment := NewMockECSEnvironmentManager(ctrl)

					ecsEnvironmentID := id.L0EnvironmentID("envid").ECSEnvironmentID()
					clusterName := ecsEnvironmentID.String()
					securityGroupID := "some_sg_id"

					mockEnvironment.ECS.EXPECT().
						CreateCluster(gomock.Any()).
						Return(ecs.NewCluster(clusterName), g.Error())

					mockEnvironment.EC2.EXPECT().
						CreateSecurityGroup(gomock.Any(), gomock.Any(), gomock.Any()).
						Return(&securityGroupID, g.Error()).
						AnyTimes()

					mockEnvironment.EC2.EXPECT().
						AuthorizeSecurityGroupIngressFromGroup(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						CreateLaunchConfiguration(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						CreateAutoScalingGroup(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						DescribeAutoScalingGroup(gomock.Any()).
						Return(autoscaling.NewGroup(), g.Error()).
						AnyTimes()

					securityGroup := ec2.NewSecurityGroup("some_sg_id")
					mockEnvironment.EC2.EXPECT().
						DescribeSecurityGroup(gomock.Any()).
						Return(securityGroup, g.Error()).
						AnyTimes()

					return mockEnvironment.Environment()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) interface{})

				for i := 0; i < 5; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some error"))

					manager := setup(g).(*ECSEnvironmentManager)
					if _, err := manager.CreateEnvironment("some_name", "m3.medium", 0, nil); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}

func TestUpdateEnvironmentMinCount(t *testing.T) {
	testModel := &models.Environment{
		EnvironmentID: "some_id",
	}

	testCases := []testutils.TestCase{
		{
			Name: "Should use proper params in aws calls",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)
				ecsEnvironmentID := id.L0EnvironmentID(testModel.EnvironmentID).ECSEnvironmentID()
				autoScalingGroupName := ecsEnvironmentID.AutoScalingGroupName()

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(autoScalingGroupName).
					Return(autoscaling.NewGroup(), nil)

				mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMinSize(autoScalingGroupName, 0).
					Return(nil)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.updateEnvironmentMinCount(testModel, 0)
			},
		},
		{
			Name: "Should set maxSize first if minClusterCount is greater",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				mockEnvironment := NewMockECSEnvironmentManager(ctrl)

				mockEnvironment.AutoScaling.EXPECT().
					DescribeAutoScalingGroup(gomock.Any()).
					Return(autoscaling.NewGroup(), nil)

				maxSizeCall := mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMaxSize(gomock.Any(), gomock.Any()).
					Return(nil)

				minSizeCall := mockEnvironment.AutoScaling.EXPECT().
					UpdateAutoScalingGroupMinSize(gomock.Any(), gomock.Any()).
					Return(nil)

				gomock.InOrder(maxSizeCall, minSizeCall)

				return mockEnvironment.Environment()
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				manager := target.(*ECSEnvironmentManager)
				manager.updateEnvironmentMinCount(testModel, 2)
			},
		},
		{
			Name: "Should propagate aws errors",
			Setup: func(reporter *testutils.Reporter, ctrl *gomock.Controller) interface{} {
				return func(g testutils.ErrorGenerator) *ECSEnvironmentManager {
					mockEnvironment := NewMockECSEnvironmentManager(ctrl)

					mockEnvironment.AutoScaling.EXPECT().
						DescribeAutoScalingGroup(gomock.Any()).
						Return(autoscaling.NewGroup(), g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						UpdateAutoScalingGroupMaxSize(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					mockEnvironment.AutoScaling.EXPECT().
						UpdateAutoScalingGroupMinSize(gomock.Any(), gomock.Any()).
						Return(g.Error()).
						AnyTimes()

					return mockEnvironment.Environment()
				}
			},
			Run: func(reporter *testutils.Reporter, target interface{}) {
				setup := target.(func(testutils.ErrorGenerator) *ECSEnvironmentManager)

				for i := 0; i < 3; i++ {
					var g testutils.ErrorGenerator
					g.Set(i+1, fmt.Errorf("some error"))

					manager := setup(g)
					if err := manager.updateEnvironmentMinCount(testModel, 2); err == nil {
						reporter.Errorf("Error on variation %d, Error was nil!", i)
					}
				}
			},
		},
	}

	testutils.RunTests(t, testCases)
}
