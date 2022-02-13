package clusterRecordModel

type ClusterRecord struct {
	Name        string `json:"name" bson:"name, omitempty"`
	Hostname    string `json:"hostname" bson:"hostname"`
	Environment string `json:"environment" bson:"environment"`
	Status      string `json:"status" bson:"status"`
}
