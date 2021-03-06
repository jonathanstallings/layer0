package main

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"github.com/quintilesims/layer0/api/handlers"
	"github.com/quintilesims/layer0/api/logic"
	"github.com/quintilesims/layer0/common/config"
	"github.com/quintilesims/layer0/common/logutils"
	"github.com/quintilesims/layer0/common/startup"
	"net/http"
	"strings"
)

func setupRestful(lgc logic.Logic) {
	adminLogic := logic.NewL0AdminLogic(lgc)
	certificateLogic := logic.NewL0CertificateLogic(lgc)
	deployLogic := logic.NewL0DeployLogic(lgc)
	environmentLogic := logic.NewL0EnvironmentLogic(lgc)
	loadBalancerLogic := logic.NewL0LoadBalancerLogic(lgc)
	serviceLogic := logic.NewL0ServiceLogic(lgc, deployLogic)
	taskLogic := logic.NewL0TaskLogic(lgc, deployLogic)
	jobLogic := logic.NewL0JobLogic(lgc, taskLogic, deployLogic)

	adminHandler := handlers.NewAdminHandler(adminLogic)
	certificateHandler := handlers.NewCertificateHandler(certificateLogic)
	deployHandler := handlers.NewDeployHandler(deployLogic)
	environmentHandler := handlers.NewEnvironmentHandler(environmentLogic, jobLogic)
	jobHandler := handlers.NewJobHandler(jobLogic)
	loadBalancerHandler := handlers.NewLoadBalancerHandler(loadBalancerLogic, jobLogic)
	serviceHandler := handlers.NewServiceHandler(serviceLogic, jobLogic)
	tagHandler := handlers.NewTagHandler(lgc.TagData)
	taskHandler := handlers.NewTaskHandler(taskLogic)

	restful.SetLogger(logutils.SilentLogger{})
	restful.Add(deployHandler.Routes())
	restful.Add(serviceHandler.Routes())
	restful.Add(environmentHandler.Routes())
	restful.Add(certificateHandler.Routes())
	restful.Add(tagHandler.Routes())
	restful.Add(adminHandler.Routes())
	restful.Add(loadBalancerHandler.Routes())
	restful.Add(taskHandler.Routes())
	restful.Add(jobHandler.Routes())

	restful.Filter(handlers.LogRequest)
	restful.Filter(handlers.AddVersionHeader)
	restful.Filter(handlers.EnableCORS)
	restful.Filter(restful.OPTIONSFilter())
	restful.DefaultContainer.Filter(handlers.HttpsRedirect)

	config := swagger.Config{
		WebServices:     restful.RegisteredWebServices(),
		ApiPath:         "/apidocs.json",
		SwaggerPath:     swaggerPath,
		SwaggerFilePath: swaggerFilePath,
		StaticHandler:   new(SwaggerRedirectHandler),
	}

	swagger.InstallSwaggerService(config)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "custom 404")
			return
		}

		http.Redirect(w, r, config.SwaggerPath, 302)
	})
}

type SwaggerRedirectHandler struct{}

var swaggerPath = "/apidocs/"
var swaggerFilePath = "api/external/swagger-ui/dist"

func (*SwaggerRedirectHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	proto := req.Header.Get("X-Forwarded-Proto")
	if proto == "http" {
		url := fmt.Sprintf("https://%v%v", req.Host, req.URL)
		http.Redirect(writer, req, url, 301)
	} else {
		http.StripPrefix(swaggerPath, http.FileServer(http.Dir(swaggerFilePath))).ServeHTTP(writer, req)
	}
}

var Version string

func main() {
	if err := config.Validate(config.RequiredAPIVariables); err != nil {
		logrus.Fatal(err)
	}

	switch strings.ToLower(config.APILogLevel()) {
	case "0", "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "1", "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "2", "warning":
		logrus.SetLevel(logrus.WarnLevel)
	case "3", "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "4", "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logger := logutils.NewStackTraceLogger("Main")
	logutils.SetGlobalLogger(logger)

	if Version == "" {
		Version = "unset/developer"
	}

	config.SetAPIVersion(Version)
	logrus.Printf("l0-api %v", Version)

	port := ":" + config.APIPort()
	region := config.AWSRegion()
	credProvider := config.NewConfigCredProvider()

	var connection string
	if config.MySQLConnection() != "" {
		connection = config.MySQLConnection()
	}

	var adminConnection string
	if config.MySQLConnection() != "" {
		adminConnection = config.MySQLAdminConnection()
	}

	backend, err := startup.GetBackend(credProvider, region, connection, adminConnection)
	if err != nil {
		logrus.Fatal(err)
	}

	lgc, err := startup.GetLogic(backend, connection, adminConnection)
	if err != nil {
		logrus.Fatal(err)
	}

	setupRestful(*lgc)

	// Since this is ECS, we run a right sizer for the clusters to kill off
	// unused clusters
	logrus.Infof("Starting RightSizer")
	backend.StartRightSizer()

	deployLogic := logic.NewL0DeployLogic(*lgc)
	taskLogic := logic.NewL0TaskLogic(*lgc, deployLogic)
	jobLogic := logic.NewL0JobLogic(*lgc, taskLogic, deployLogic)
	jobJanitor := logic.NewJobJanitor(jobLogic)

	logrus.Infof("Starting Job Janitor")
	jobJanitor.Run()

	logrus.Print("Service on localhost" + port)
	logrus.Fatal(http.ListenAndServe(port, nil))
}
