package ecsbackend

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/quintilesims/layer0/api/backend/ecs/id"
	"github.com/quintilesims/layer0/common/aws/cloudwatchlogs"
	"github.com/quintilesims/layer0/common/aws/ecs"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/models"
	"strings"
)

const MAX_TASK_IDS = 100

func stringp(s string) *string {
	return &s
}

func int64p(i int64) *int64 {
	return &i
}

func boolp(b bool) *bool {
	return &b
}

func pstring(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func int64OrZero(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}

func ContainsErrCode(err error, code string) bool {
	if err == nil {
		return false
	}

	awsErr, ok := err.(awserr.Error)
	if !ok {
		return false
	}

	return strings.ToLower(awsErr.Code()) == strings.ToLower(code)
}

func ContainsErrMsg(err error, msg string) bool {
	if err == nil {
		return false
	}

	return strings.Contains(
		strings.ToLower(err.Error()),
		strings.ToLower(msg))
}

// IteratePages performs a do-while loop on a paginatedf
// until nextToken is nil or an error is returned
type paginatedf func(*string) (*string, error)

func IteratePages(fn paginatedf) error {
	var err error
	var nextToken *string

	nextToken, err = fn(nextToken)
	if err != nil {
		return err
	}

	for nextToken != nil {
		nextToken, err = fn(nextToken)
		if err != nil {
			return err
		}
	}

	return nil
}

var CreateRenderedDeploy = func(body []byte) (*deploy, error) {
	deploy, err := marshalDeploy(body)
	if err != nil {
		return nil, err
	}

	for _, container := range deploy.ContainerDefinitions {
		if container.LogConfiguration == nil {
			container.LogConfiguration = &awsecs.LogConfiguration{
				LogDriver: stringp("awslogs"),
				Options: map[string]*string{
					"awslogs-group":         stringp(config.AWSLogGroupID()),
					"awslogs-region":        stringp(config.AWSRegion()),
					"awslogs-stream-prefix": stringp("l0"),
				},
			}
		}
	}

	return deploy, nil
}

var getTaskARNs = func(ecs ecs.Provider, ecsEnvironmentID id.ECSEnvironmentID, startedBy *string) ([]*string, error) {
	// we can only check each of the states individually, thus we must issue 3 API calls

	running := "RUNNING"
	tasks, err := ecs.ListTasks(ecsEnvironmentID.String(), nil, &running, startedBy, nil)
	if err != nil {
		return nil, err
	}

	stopped := "STOPPED"
	stoppedTasks, err := ecs.ListTasks(ecsEnvironmentID.String(), nil, &stopped, startedBy, nil)
	if err != nil {
		return nil, err
	}

	pending := "PENDING"
	pendingTasks, err := ecs.ListTasks(ecsEnvironmentID.String(), nil, &pending, startedBy, nil)
	if err != nil {
		return nil, err
	}

	tasks = append(tasks, stoppedTasks...)
	tasks = append(tasks, pendingTasks...)

	return tasks, nil
}

var GetLogs = func(cloudWatchLogs cloudwatchlogs.Provider, taskARNs []*string, tail int) ([]*models.LogFile, error) {
	taskIDCatalog := generateTaskIDCatalog(taskARNs)

	orderBy := "LogStreamName"
	logStreams, err := cloudWatchLogs.DescribeLogStreams(config.AWSLogGroupID(), orderBy)
	if err != nil {
		return nil, err
	}

	logFiles := []*models.LogFile{}
	for _, logStream := range logStreams {
		// filter by streams that have <prefix>/<container name>/<stream task id>
		streamNameSplit := strings.Split(*logStream.LogStreamName, "/")
		if len(streamNameSplit) != 3 {
			continue
		}

		streamTaskID := streamNameSplit[2]
		if _, ok := taskIDCatalog[streamTaskID]; !ok {
			continue
		}

		logFile := &models.LogFile{
			Name:  streamNameSplit[1],
			Lines: []string{},
		}

		// since the time range is exclusive, expand the range to get first/last events
		logEvents, err := cloudWatchLogs.GetLogEvents(
			config.AWSLogGroupID(),
			*logStream.LogStreamName,
			*logStream.FirstEventTimestamp-1,
			*logStream.LastEventTimestamp+1,
			int64(tail))
		if err != nil {
			return nil, err
		}

		for _, logEvent := range logEvents {
			logFile.Lines = append(logFile.Lines, *logEvent.Message)
		}

		logFiles = append(logFiles, logFile)
	}

	return logFiles, nil
}

func generateTaskIDCatalog(taskARNs []*string) map[string]bool {
	catalog := map[string]bool{}
	for _, taskARN := range taskARNs {
		taskID := strings.Split(*taskARN, "/")[1]
		catalog[taskID] = true
	}

	return catalog
}
