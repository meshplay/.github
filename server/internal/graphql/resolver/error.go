package resolver

import (
	"github.com/khulnasoft/meshkit/errors"
)

// Please reference the following before contributing an error code:
// https://docs.meshplay.khulnasofy.com/project/contributing/contributing-error
// https://github.com/meshplay/meshkit/blob/master/errors/errors.go
const (
	ErrResolverInvalidRequestCode           = "meshplay-server-1184"
	ErrResolverNilClientCode                = "meshplay-server-1185"
	ErrResolverMeshsyncSubscriptionCode     = "meshplay-server-1186"
	ErrResolverOperatorSubscriptionCode     = "meshplay-server-1187"
	ErrAddonSubscriptionCode                = "meshplay-server-1188"
	ErrResolverControlPlaneSubscriptionCode = "meshplay-server-1189"
	ErrMeshplayClientCode                    = "meshplay-server-1190"
	ErrResolverPublishBrokerCode            = "meshplay-server-1191"
	ErrNoMeshSyncCode                       = "meshplay-server-1192"
	ErrDataPlaneSubscriptionCode            = "meshplay-server-1193"
	ErrBrokerNotConnectedCode               = "meshplay-server-1194"
	ErrGettingNamespaceCode                 = "meshplay-server-1195"
	ErrFetchingPatternsCode                 = "meshplay-server-1196"
	ErrInvalidOAMTypeCode                   = "meshplay-server-1197"
	ErrKubectlDescribeCode                  = "meshplay-server-1198"
	ErrEmptyCurrentK8sContextCode           = "meshplay-server-1199"
	ErrConfigurationPatternsCode            = "meshplay-server-1200"
	ErrConfigurationApplicationsCode        = "meshplay-server-1201"
	ErrConfigurationFiltersCode             = "meshplay-server-1202"
	ErrK8sContextCode                       = "meshplay-server-1203"
	ErrClusterResourcesSubscriptionCode     = "meshplay-server-1204"
	ErrGettingClusterResourcesCode          = "meshplay-server-1205"
	ErrMeshModelSummarySubscriptionCode     = "meshplay-server-1206"
	ErrGettingMeshModelSummaryCode          = "meshplay-server-1207"
	ErrGettingRegistryManagerCode           = "meshplay-server-1208"
	ErrGettingTelemetryComponentsCode       = "meshplay-server-1209"
	ErrAdapterInsufficientInformationCode   = "meshplay-server-1210"
	ErrPerformanceProfilesSubscriptionCode  = "meshplay-server-1211"
	ErrPerformanceResultSubscriptionCode    = "meshplay-server-1212"
	ErrGormDatabaseCode                     = "meshplay-server-1213"
)

var (
	ErrNilClient              = errors.New(ErrResolverNilClientCode, errors.Alert, []string{"Kubernetes client not initialized"}, []string{"Kubernetes config is not initialized with Meshplay"}, []string{}, []string{"Upload your kubernetes config via the settings dashboard. If uploaded, wait for a minute for it to get initialized"})
	ErrInvalidRequest         = errors.New(ErrResolverInvalidRequestCode, errors.Alert, []string{"Invalid query, please check syntax"}, []string{"The Graphql query requested is invalid"}, []string{}, []string{"Check the query parameters and syntax of the query being run"})
	ErrNoMeshSync             = errors.New(ErrNoMeshSyncCode, errors.Alert, []string{"MeshSync disabled"}, []string{"MeshSync custom controller is not running in your kubernetes cluster"}, []string{"Meshplay Operator is not running in your cluster or is crashing"}, []string{"Enable Meshplay Operator from the settings page in the UI", "Check for logs in the meshplay-operator pods from inside the application for more information"})
	ErrBrokerNotConnected     = errors.New(ErrBrokerNotConnectedCode, errors.Alert, []string{"Broker not connected"}, []string{"Meshplay Broker is not connected to Meshplay Server"}, []string{"Meshplay Broker is crashed or not reachable"}, []string{"Restart Meshplay Server", "Please check if Meshplay server has access to the Broker"})
	ErrEmptyCurrentK8sContext = errors.New(ErrEmptyCurrentK8sContextCode, errors.Alert, []string{"Current kubernetes context is empty"}, []string{"Meshplay is not able to get the current kubernetes context"}, []string{"Meshplay is crashed or not reachable"}, []string{"Restart Meshplay Server", "Please check if Meshplay server is accessible to the Database"})
)

func ErrMeshsyncSubscription(err error) error {
	return errors.New(ErrResolverMeshsyncSubscriptionCode, errors.Alert, []string{"MeshSync Subscription failed", err.Error()}, []string{"GraphQL subscription for MeshSync stopped"}, []string{"Could be a network issue"}, []string{"Check if meshplay server is reachable from the browser"})
}

func ErrOperatorSubscription(err error) error {
	return errors.New(ErrResolverOperatorSubscriptionCode, errors.Alert, []string{"Operator Subscription failed", err.Error()}, []string{"GraphQL subscription for Operator stopped"}, []string{"Could be a network issue"}, []string{"Check if meshplay server is reachable from the browser"})
}

func ErrAddonSubscription(err error) error {
	return errors.New(ErrAddonSubscriptionCode, errors.Alert, []string{"Addons Subscription failed", err.Error()}, []string{"GraphQL subscription for Addons stopped"}, []string{"Could be a network issue"}, []string{"Check if meshplay server is reachable from the browser"})
}

func ErrPerformanceProfilesSubscription(err error) error {
	return errors.New(ErrPerformanceProfilesSubscriptionCode, errors.Alert, []string{"PerformanceProfiles Subscription failed", err.Error()}, []string{"GraphQL subscription for PerformanceProfiles stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrPerformanceResultSubscription(err error) error {
	return errors.New(ErrPerformanceResultSubscriptionCode, errors.Alert, []string{"PerformanceResult Subscription failed", err.Error()}, []string{"GraphQL subscription for PerformanceResult stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrGormDatabase(err error) error {
	return errors.New(ErrGormDatabaseCode, errors.Alert, []string{"Database operation failed", err.Error()}, []string{"Database operation failed. Please ensure that the database isn't corrupted"}, []string{"Could be a db issue"}, []string{"Confirm that database connection is working as expected."})
}

func ErrControlPlaneSubscription(err error) error {
	return errors.New(ErrResolverControlPlaneSubscriptionCode, errors.Alert, []string{"Control Plane Subscription failed", err.Error()}, []string{"GraphQL subscription for Control Plane stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrDataPlaneSubscription(err error) error {
	return errors.New(ErrDataPlaneSubscriptionCode, errors.Alert, []string{"Data Plane Subscription failed", err.Error()}, []string{"GraphQL subscription for Data Plane stopped"}, []string{"Could be a network issue"}, []string{"Check if meshplay server is reachable from the browser"})
}

func ErrPublishBroker(err error) error {
	return errors.New(ErrResolverPublishBrokerCode, errors.Alert, []string{"Unable to publish to broker", err.Error()}, []string{"Unable to create a broker publisher"}, []string{"Could be a network issue", "Meshplay Broker could have crashed"}, []string{"Check if Meshplay Broker is reachable from Meshplay Server", "Check if Meshplay Broker is up and running inside the configured cluster"})
}

func ErrGettingNamespace(err error) error {
	return errors.New(ErrGettingNamespaceCode, errors.Alert, []string{"Cannot get available namespaces"}, []string{err.Error()}, []string{"The table in the database might not exist"}, []string{})
}

func ErrFetchingPatterns(err error) error {
	return errors.New(ErrFetchingPatternsCode, errors.Alert, []string{"Cannot fetch designs"}, []string{err.Error()}, []string{"There might be something wrong with the Meshplay or Meshplay Cloud"}, []string{"Try again, if still exist, please post an issue on Meshplay repository"})
}

func ErrInvalidOAMType() error {
	return errors.New(
		ErrInvalidOAMTypeCode,
		errors.Alert,
		[]string{"invalid oam type is requested"},
		[]string{"invalid oam type requested, supported types are workload, scope, trait "},
		[]string{}, []string{},
	)
}

func ErrKubectlDescribe(err error) error {
	return errors.New(
		ErrKubectlDescribeCode,
		errors.Alert,
		[]string{"failed to find the resource", "invalid resource type", "resource doens't exists"},
		[]string{err.Error(), "invalid kubernetes resource type or couldn't find the specified resource"},
		[]string{}, []string{},
	)
}

func ErrPatternsSubscription(err error) error {
	return errors.New(ErrConfigurationPatternsCode, errors.Alert, []string{"Configuration Subscription failed", err.Error()}, []string{"GraphQL subscription for designs stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrApplicationsSubscription(err error) error {
	return errors.New(ErrConfigurationApplicationsCode, errors.Alert, []string{"Configuration Subscription failed", err.Error()}, []string{"GraphQL subscription for Applications stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrFiltersSubscription(err error) error {
	return errors.New(ErrConfigurationFiltersCode, errors.Alert, []string{"Configuration Subscription failed", err.Error()}, []string{"GraphQL subscription for Filters stopped"}, []string{"Could be a network issue"}, []string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrClusterResourcesSubscription(err error) error {
	return errors.New(
		ErrClusterResourcesSubscriptionCode,
		errors.Alert,
		[]string{"ClusterResources Subscription failed", err.Error()},
		[]string{"GraphQL subscription for ClusterResources Subscription stopped"},
		[]string{"Could be a network issue"},
		[]string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrGettingClusterResources(err error) error {
	return errors.New(
		ErrGettingClusterResourcesCode,
		errors.Alert,
		[]string{"Unable to retrieve cluster resources"},
		[]string{err.Error()},
		[]string{"Table in the database might not exists"},
		[]string{""},
	)
}

func ErrMeshModelSummarySubscription(err error) error {
	return errors.New(
		ErrMeshModelSummarySubscriptionCode,
		errors.Alert,
		[]string{"MeshModelSummary Subscription failed", err.Error()},
		[]string{"GraphQL subscription for MeshModelSummary Subscription stopped"},
		[]string{"Could be a network issue"},
		[]string{"Confirm that Meshplay Server is reachable from your browser."})
}

func ErrGettingMeshModelSummary(err error) error {
	return errors.New(
		ErrGettingMeshModelSummaryCode,
		errors.Alert,
		[]string{"Unable to retrieve MeshModel Summary"},
		[]string{err.Error()},
		[]string{"Table in the database might not exists"},
		[]string{""},
	)
}

func ErrGettingRegistryManager(err error) error {
	return errors.New(
		ErrGettingRegistryManagerCode,
		errors.Alert,
		[]string{"Unable to retrieve Registry Manager"},
		[]string{err.Error()},
		[]string{"Registry Manager might not exists"},
		[]string{""},
	)
}

func ErrK8sContextSubscription(err error) error {
	return errors.New(ErrK8sContextCode, errors.Alert, []string{"Failed to get k8s context from remote provider", err.Error()}, []string{"There might be something wrong with the Meshplay or Meshplay Cloud"}, []string{"Could be a network issue"}, []string{""})
}

func ErrGettingTelemetryComponents(err error) error {
	return errors.New(ErrGettingTelemetryComponentsCode, errors.Alert, []string{"unable to retrieve telemetry components"}, []string{err.Error()}, []string{"table in the database might be corrupted"}, []string{"try resetting database from settings"})
}

func ErrAdapterInsufficientInformation(err error) error {
	return errors.New(ErrAdapterInsufficientInformationCode, errors.Critical, []string{"Unable to process adapter request, incomplete request"}, []string{err.Error()}, []string{}, []string{})
}
