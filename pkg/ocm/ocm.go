package ocm

import (
	"fmt"
	"time"

	sdk "github.com/openshift-online/ocm-sdk-go"
	clv1 "github.com/openshift-online/ocm-sdk-go/servicelogs/v1"
)

// ServiceLog is what is received from OCM
type ServiceLog struct {
	ClusterId     string
	ClusterUuid   string
	CreatedAt     time.Time
	CreatedBy     string
	Desc          string
	EventStreamId string
	Href          string
	Id            string
	InternalOnly  bool
	Kind          string
	LogType       string
	ServiceName   string
	Severity      string
	Summary       string
	Timestamp     time.Time
	Username      string
}

type Client struct {
	// conn sdk.Connection undecided if necessary here or not
	// Ocm Cluster Client
	// clusterClient *cm v1.ClustersClient
	clusterLogsClient    *clv1.ClustersClusterLogsClient
	clusterLogsAddClient *clv1.ClusterLogsClient
}

// NewClient Not sure if I want this to be part of the ocmClient Struct yet.
// Any ways it needs to be exposed to the user for them to close the connection
func NewClient(conn *sdk.Connection) Client {
	return Client{
		clusterLogsClient:    conn.ServiceLogs().V1().Clusters().ClusterLogs(),
		clusterLogsAddClient: conn.ServiceLogs().V1().ClusterLogs(),
	}
}

func NewConnectionWithTemporaryToken(url, token string) (*sdk.Connection, error) {
	connection, err := sdk.NewConnectionBuilder().URL(url).Tokens(token).Build()
	if err != nil {
		return nil, fmt.Errorf("error building ocm sdk connection :: %q \n", err)
	}

	return connection, nil
}

func (c Client) PostInternalServiceLog(clusterId string, description string) error {
	logEntry, err := clv1.NewLogEntry().
		InternalOnly(true).
		ClusterID(clusterId).
		Severity("Info").
		ServiceName("SREManualAction").
		Summary("INTERNAL ONLY, DO NOT SHARE WITH CUSTOMER").
		Description(description).
		Build()
	if err != nil {
		return err
	}
	clusterLogsAddResponse, err := c.clusterLogsAddClient.Add().Body(logEntry).Send()
	if err != nil {
		return err
	} else if clusterLogsAddResponse.Status() != 201 {
		return fmt.Errorf("expected 201 when adding service log but got %d", clusterLogsAddResponse.Status())
	}
	return nil
}

func (c Client) ListServiceLogs(clusterID string, query ...string) ([]ServiceLog, error) {
	queryString := ""
	for i, s := range query {
		if i != 0 {
			queryString += fmt.Sprintf(" and %s", s)
		}
	}

	list := make([]ServiceLog, 0)
	page := 1
	size := 1000
	for {
		resp, err := c.clusterLogsClient.List().
			ClusterID(clusterID).
			Search(queryString).
			Size(size).
			Page(page).
			Send()
		if err != nil {
			return []ServiceLog{}, err
		}

		resp.Items().Each(func(logEntry *clv1.LogEntry) bool {
			list = append(list, ServiceLog{
				logEntry.ClusterID(),
				logEntry.ClusterUUID(),
				logEntry.CreatedAt(),
				logEntry.CreatedBy(),
				logEntry.Description(),
				logEntry.EventStreamID(),
				logEntry.HREF(),
				logEntry.ID(),
				logEntry.InternalOnly(),
				logEntry.Kind(),
				string(logEntry.LogType()),
				logEntry.ServiceName(),
				string(logEntry.Severity()),
				logEntry.Summary(),
				logEntry.Timestamp(),
				logEntry.Username(),
			})
			return true
		})

		if resp.Size() < size {
			break
		}
		page++
	}

	return list, nil
}
