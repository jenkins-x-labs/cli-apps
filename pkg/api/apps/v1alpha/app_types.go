package v1alpha

const (
	// AppConfigFileName is the name of the applications configuration file
	AppConfigFileName = "jx-apps.yml"
	// PhaseSystem is installed before the apps phase
	PhaseSystem Phase = "system"
	// PhaseApps is installed after the system phase
	PhaseApps Phase = "apps"
)

// PhaseValues the string values for Phases
var PhaseValues = []string{"system", "apps"}

// AppConfig contains the apps to install during boot for helmfile / helm 3
type AppConfig struct {
	// Apps of applications
	Apps []App `json:"apps"`
	// Repositories list of helm repositories
	Repositories []RepositorySpec `json:"repositories,omitempty"`
	// DefaultNamespace the default namespace to install applications into
	DefaultNamespace string `json:"defaultNamespace,omitempty"`
}

// RepositorySpec that defines values for a helm repo
type RepositorySpec struct {
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	CaFile   string `json:"caFile,omitempty"`
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

// App is the configuration of an app used during boot for helmfile / helm 3
type App struct {
	// Name of the application / helm chart
	Name string `json:"name,omitempty"`
	// Repository the helm repository
	Repository string `json:"repository,omitempty"`
	// Namespace to install the application into
	Namespace string `json:"namespace,omitempty"`
	// Phase of the pipeline to install application
	Phase Phase `json:"phase,omitempty"`
	// Version the version to install if you want to override the version from the Version Stream.
	// Note we recommend using the version stream for app versions
	Version string `json:"version,omitempty"`
	// Description an optional description of the app
	Description string `json:"description,omitempty"`
	// Values any explicit value files to be used
	Values []string `json:"values,omitempty"`
}

// Phase of the pipeline to install application
type Phase string
