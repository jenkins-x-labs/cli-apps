package v1alpha

import (
	"fmt"
	"github.com/jenkins-x-labs/cli-apps/pkg/external/gitops/gitter"
	"github.com/jenkins-x-labs/cli-apps/pkg/external/schema_utils"
	v1 "github.com/jenkins-x/jx/pkg/apis/jenkins.io/v1"
	"github.com/jenkins-x/jx/pkg/helm"
	"github.com/jenkins-x/jx/pkg/util"
	"github.com/pkg/errors"
	"io/ioutil"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"strings"
)

// AppsService the Jenkins X apps service
type AppsService struct {
}

func cloneEnvironmentRepo(environmentRepoUrl string) (string, error) {
	gitter := gits.NewGitCLI()
	// TODO: Cleanup temp dir
	dir, err := ioutil.TempDir("", "helmboot-")
	if err != nil {
		return "", errors.Wrap(err, "failed to create temporary directory")
	}
	err = gitter.Clone(environmentRepoUrl, dir)
	return dir, err
}

func getAppsConfig(environmentRepoUrl string) (*AppConfig, string, error) {
	dir, err := cloneEnvironmentRepo(environmentRepoUrl)
	appsConfig, err := loadAppConfig(dir)
	if err != nil {
		return nil, "", errors.Wrapf(err, "failed to load applications config in environment clone of %s", ".")
	}
	return appsConfig, dir, nil
}

// ListAppsForEnvironment lists the apps installed in an environment
func (as *AppsService) ListAppsForEnvironment(environmentRepoUrl string) (*v1.AppList, error) {
	appsList := &v1.AppList{}
	appsConfig, _, err := getAppsConfig(environmentRepoUrl)
	if err != nil {
		return nil, err
	}
	for _, a := range appsConfig.Apps {
		app := v1.App{}
		// TODO: Sanitize the name
		// A namespace was also defaulted here, can re drop that now? We should know this due to the environment repo used
		app.Name = a.Name
		app.Labels = map[string]string{
			helm.LabelAppName:    a.Name,
			helm.LabelAppVersion: a.Version,
		}
		app.Annotations = map[string]string{
			helm.AnnotationAppDescription: a.Description,
			helm.AnnotationAppRepository:  a.Repository,
		}
		appsList.Items = append(appsList.Items, app)
	}
	return appsList, nil
}

// AddAppForEnvironment Adds an app to an environment
func (as *AppsService) AddAppForEnvironment(newApp *App, environmentRepoUrl string) error {
	appsConfig, envDir, err := getAppsConfig(environmentRepoUrl)
	if err != nil {
		return err
	}
	appFound := false
	for i, app := range appsConfig.Apps {
		if app.Name == newApp.Name {
			appsConfig.Apps[i] = *newApp
			appFound = true
		}
	}
	if !appFound {
		appsConfig.Apps = append(appsConfig.Apps, *newApp)
	}
	err = saveConfig(appsConfig, filepath.Join(envDir, AppConfigFileName))
	if err != nil {
		return err
	}
	fmt.Println(envDir)
	return nil
}

// RemoveAppFromEnvironment Removes an app from an environment
func (as *AppsService) RemoveAppFromEnvironment(newApp *App, environmentRepoUrl string) error {
	appsConfig, envDir, err := getAppsConfig(environmentRepoUrl)
	if err != nil {
		return err
	}
	appFound := false
	newAppList := make([]App, 0)
	for i, app := range appsConfig.Apps {
		if app.Name != newApp.Name {
			newAppList = append(newAppList, appsConfig.Apps[i])

		} else {
			appFound = true
		}
	}
	if appFound {
		appsConfig.Apps = newAppList
		err = saveConfig(appsConfig, filepath.Join(envDir, AppConfigFileName))
		if err != nil {
			return err
		}
		operations.CreatePull
	}
	fmt.Println(envDir)
	return nil
}

// LoadAppConfig loads the boot applications configuration file
// if there is not a file called `jx-apps.yml` in the given dir we will scan up the parent
// directories looking for the requirements file as we often run 'jx' steps in sub directories.
func loadAppConfig(dir string) (*AppConfig, error) {
	fileName := AppConfigFileName
	if dir != "" {
		fileName = filepath.Join(dir, fileName)
	}

	exists, err := util.FileExists(fileName)
	if err != nil {
		return nil, errors.Errorf("error looking up %s in directory %s", fileName, dir)
	}

	config := &AppConfig{}
	if !exists {
		return config, nil
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return config, fmt.Errorf("Failed to load file %s due to %s", fileName, err)
	}
	validationErrors, err := schema_utils.ValidateYaml(config, data)
	if err != nil {
		return config, fmt.Errorf("failed to validate YAML file %s due to %s", fileName, err)
	}
	if len(validationErrors) > 0 {
		return config, fmt.Errorf("Validation failures in YAML file %s:\n%s", fileName, strings.Join(validationErrors, "\n"))
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return config, fmt.Errorf("Failed to unmarshal YAML file %s due to %s", fileName, err)
	}

	// validate all phases are known types, default to apps if not specified
	for _, app := range config.Apps {
		if app.Phase != "" {
			if app.Phase != PhaseSystem && app.Phase != PhaseApps {
				return config, fmt.Errorf("failed to validate YAML file, invalid phase '%s', needed on of %v",
					string(app.Phase), PhaseValues)
			}
		}
	}
	return config, err
}

// SaveConfig saves the configuration file to the given project directory
func saveConfig(c *AppConfig, fileName string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileName, data, util.DefaultWritePermissions)
	if err != nil {
		return errors.Wrapf(err, "failed to save file %s", fileName)
	}
	return nil
}
