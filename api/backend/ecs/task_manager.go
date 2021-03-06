package ecsbackend

import (
	log "github.com/Sirupsen/logrus"
	"github.com/quintilesims/layer0/api/backend"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/errors"
	"github.com/quintilesims/layer0/common/models"
	"strings"
)

var ClusterCapacityReason = "Waiting for cluster capacity to run"

type ECSTaskManager struct {
	ECS            ecs.Provider
	CloudWatchLogs cloudwatchlogs.Provider
	Backend        backend.Backend
	ClusterScaler  ClusterScaler
	Scheduler      TaskScheduler
}

func NewECSTaskManager(
	ecsProvider ecs.Provider,
	cloudWatchLogsProvider cloudwatchlogs.Provider,
	backend backend.Backend,
	clusterScaler ClusterScaler,
) *ECSTaskManager {
	return &ECSTaskManager{
		ECS:            ecsProvider,
		CloudWatchLogs: cloudWatchLogsProvider,
		Backend:        backend,
		ClusterScaler:  clusterScaler,
		Scheduler:      NewL0TaskScheduler(ecsProvider),
	}
}

func (this *ECSTaskManager) ListTasks() ([]*models.Task, error) {
	environments, err := this.Backend.ListEnvironments()
	if err != nil {
		return nil, err
	}

	taskCopies := map[string][]*ecs.Task{}
	for _, environment := range environments {
		ecsEnvironmentID := id.L0EnvironmentID(environment.EnvironmentID).ECSEnvironmentID()

		taskARNs, err := getTaskARNs(this.ECS, ecsEnvironmentID, nil)
		if err != nil {
			return nil, err
		}

		if len(taskARNs) > 0 {
			tasks, err := this.describeTasks(ecsEnvironmentID, taskARNs)
			if err != nil {
				return nil, err
			}

			for _, task := range tasks {
				startedBy := stringOrEmpty(task.StartedBy)

				if strings.HasPrefix(startedBy, id.PREFIX) {
					if _, ok := taskCopies[startedBy]; !ok {
						taskCopies[startedBy] = []*ecs.Task{}
					}

					taskCopies[startedBy] = append(taskCopies[startedBy], task)
				}
			}
		}
	}

	tasks := []*models.Task{}
	for _, copies := range taskCopies {
		model, err := modelFromTasks(copies)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, model)
	}

	scheduledTasks := this.Scheduler.ListTasks()
	for _, task := range scheduledTasks {
		model, err := modelFromTasks([]*ecs.Task{task})
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, model)
	}

	return tasks, nil
}

func (this *ECSTaskManager) GetTask(environmentID, taskID string) (*models.Task, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	tasks, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(ecsTaskID.String()))
	if err != nil {
		return nil, err
	}

	taskDescs := []*ecs.Task{}
	if len(tasks) > 0 {
		taskDescs, err = this.describeTasks(ecsEnvironmentID, tasks)
		if err != nil {
			return nil, err
		}
	}

	taskDescs = append(taskDescs, this.Scheduler.GetTask(ecsTaskID)...)

	return modelFromTasks(taskDescs)
}

func (this *ECSTaskManager) DeleteTask(environmentID, taskID string) error {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	taskARNs, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(ecsTaskID.String()))
	if err != nil {
		return err
	}

	// This stops the task, later reaping by AWS will prevent it from being returned.
	reason := "Task stopped by User"

	for _, taskARN := range taskARNs {
		if err := this.ECS.StopTask(ecsEnvironmentID.String(), reason, *taskARN); err != nil {
			return err
		}
	}

	this.Scheduler.DeleteTask(ecsTaskID)

	return nil
}

func (this *ECSTaskManager) CreateTask(
	environmentID string,
	taskName string,
	deployID string,
	copies int,
	overrides []models.ContainerOverride,
) (*models.Task, error) {

	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsDeployID := id.L0DeployID(deployID).ECSDeployID()

	taskID := id.GenerateHashedEntityID(taskName)
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	// trigger the scaling algorithm first or the task we are about to create gets
	// included in the pending count of the cluster
	if _, _, err := this.ClusterScaler.TriggerScalingAlgorithm(ecsEnvironmentID, &ecsDeployID, copies); err != nil {
		return nil, err
	}

	ecsOverrides := []*ecs.ContainerOverride{}
	for _, override := range overrides {
		o := ecs.NewContainerOverride(override.ContainerName, override.EnvironmentOverrides)
		ecsOverrides = append(ecsOverrides, o)
	}

	tasks, failed, err := this.ECS.RunTask(ecsEnvironmentID.String(), ecsDeployID.TaskDefinition(), int64(copies), stringp(ecsTaskID.String()), ecsOverrides)
	if err != nil {
		if !ContainsErrMsg(err, "No Container Instances were found in your cluster") {
			return nil, err
		}

		log.Debugf("Not enough room in cluster. Adding task '%s' to scheduler", ecsTaskID)
		this.Scheduler.AddTask(ecsTaskID, ecsDeployID, ecsEnvironmentID, copies, overrides)

		dummyTask := ecsPendingTask(ecsTaskID, ecsDeployID, ecsEnvironmentID)
		tasks = []*ecs.Task{dummyTask}
	}

	if numFailed := len(failed); numFailed > 0 {
		log.Debug("RunTask failed to start %d tasks: %v", numFailed, failed)
		log.Debug("Adding task '%s' to scheduler", ecsTaskID)

		this.Scheduler.AddTask(ecsTaskID, ecsDeployID, ecsEnvironmentID, numFailed, overrides)

		dummyTask := ecsPendingTask(ecsTaskID, ecsDeployID, ecsEnvironmentID)
		tasks = []*ecs.Task{dummyTask}
	}

	return modelFromTasks(tasks)
}

func (this *ECSTaskManager) GetTaskLogs(environmentID, taskID string, tail int) ([]*models.LogFile, error) {
	ecsEnvironmentID := id.L0EnvironmentID(environmentID).ECSEnvironmentID()
	ecsTaskID := id.L0TaskID(taskID).ECSTaskID()

	taskARNs, err := getTaskARNs(this.ECS, ecsEnvironmentID, stringp(ecsTaskID.String()))
	if err != nil {
		return nil, err
	}

	return GetLogs(this.CloudWatchLogs, taskARNs, tail)
}

// Assumes the tasks are all of the same type
func modelFromTasks(tasks []*ecs.Task) (*models.Task, error) {
	if len(tasks) == 0 {
		return nil, errors.Newf(errors.InvalidTaskID, "The specified task does not exist")
	}

	var pendingCount, runningCount int64
	copies := []models.TaskCopy{}
	for _, task := range tasks {
		if *task.LastStatus == "RUNNING" {
			runningCount = runningCount + 1
		} else if *task.LastStatus == "PENDING" {
			pendingCount = pendingCount + 1
		}

		details := []models.TaskDetail{}
		for _, container := range task.Containers {
			detail := models.TaskDetail{
				ContainerName: *container.Name,
				LastStatus:    *container.LastStatus,
				Reason:        stringOrEmpty(container.Reason),
				ExitCode:      int64OrZero(container.ExitCode),
			}

			details = append(details, detail)
		}

		copy := models.TaskCopy{
			Details:    details,
			Reason:     stringOrEmpty(task.StoppedReason),
			TaskCopyID: stringOrEmpty(task.TaskArn),
		}

		copies = append(copies, copy)
	}

	model := &models.Task{
		EnvironmentID: id.ClusterARNToECSEnvironmentID(*tasks[0].ClusterArn).L0EnvironmentID(),
		PendingCount:  pendingCount,
		RunningCount:  runningCount,
		DesiredCount:  int64(len(tasks)),
		TaskID:        id.ECSTaskID(*tasks[0].StartedBy).L0TaskID(),
		Copies:        copies,
		DeployID:      id.TaskDefinitionARNToECSDeployID(*tasks[0].TaskDefinitionArn).L0DeployID(),
	}

	return model, nil
}

func (this *ECSTaskManager) describeTasks(ecsEnvironmentID id.ECSEnvironmentID, taskARNs []*string) ([]*ecs.Task, error) {
	ret := []*ecs.Task{}
	for i := len(taskARNs); i > 0; i = len(taskARNs) {
		if i > MAX_TASK_IDS {
			i = MAX_TASK_IDS
		}

		output, err := this.ECS.DescribeTasks(ecsEnvironmentID.String(), taskARNs[:i])
		if err != nil {
			return nil, err
		}

		ret = append(ret, output...)
		taskARNs = taskARNs[i:]
	}

	return ret, nil
}
