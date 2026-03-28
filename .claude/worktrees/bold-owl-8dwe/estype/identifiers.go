package estype

// DocumentID represents a Document identifier.
type DocumentID string

// String returns the string representation of the DocumentID.
func (d DocumentID) String() string { return string(d) }

// DataStream represents an Elasticsearch Data Stream name.
type DataStream string

// String returns the string representation of the DataStream.
func (d DataStream) String() string { return string(d) }

// Policy represents an ILM or Snapshot Lifecycle Policy name.
type Policy string

// String returns the string representation of the Policy.
func (p Policy) String() string { return string(p) }

// Pipeline represents an Ingest Pipeline ID.
type Pipeline string

// String returns the string representation of the Pipeline.
func (p Pipeline) String() string { return string(p) }

// Template represents an Index Template name.
type Template string

// String returns the string representation of the Template.
func (t Template) String() string { return string(t) }

// Repository represents a Snapshot Repository name.
type Repository string

// String returns the string representation of the Repository.
func (r Repository) String() string { return string(r) }

// Snapshot represents a Snapshot name.
type Snapshot string

// String returns the string representation of the Snapshot.
func (s Snapshot) String() string { return string(s) }

// TaskID represents a Task identifier.
type TaskID string

// String returns the string representation of the TaskID.
func (t TaskID) String() string { return string(t) }

// InferenceID represents an ML Inference Endpoint ID.
type InferenceID string

// String returns the string representation of the InferenceID.
func (i InferenceID) String() string { return string(i) }

// MLJobID represents an Anomaly Detection Job ID.
type MLJobID string

// String returns the string representation of the MLJobID.
func (m MLJobID) String() string { return string(m) }

// DatafeedID represents an ML Datafeed ID.
type DatafeedID string

// String returns the string representation of the DatafeedID.
func (d DatafeedID) String() string { return string(d) }

// TransformID represents a Transform Job ID.
type TransformID string

// String returns the string representation of the TransformID.
func (t TransformID) String() string { return string(t) }

// DataFrameAnalyticsID represents an ML Data Frame Analytics Job ID.
type DataFrameAnalyticsID string

// String returns the string representation of the DataFrameAnalyticsID.
func (d DataFrameAnalyticsID) String() string { return string(d) }

// TrainedModelID represents an ML Trained Model ID.
type TrainedModelID string

// String returns the string representation of the TrainedModelID.
func (t TrainedModelID) String() string { return string(t) }

// KeepAlive represents a time duration string for Point in Time (e.g. "1m").
type KeepAlive string

// String returns the string representation of the KeepAlive.
func (k KeepAlive) String() string { return string(k) }

// ESQLQuery represents an ES|QL query string.
type ESQLQuery string

// String returns the string representation of the ESQLQuery.
func (e ESQLQuery) String() string { return string(e) }
