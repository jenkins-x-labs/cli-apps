package v1alpha_test

import (
	"github.com/jenkins-x-labs/cli-apps/pkg/api/apps/v1alpha"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ListAppsForEnvironment(t *testing.T) {
	appsService := v1alpha.AppsService{}
	appsList, err := appsService.ListAppsForEnvironment("git@github.com:jstrachan/environment-jstrachan-jx5-dev.git")
	assert.NoError(t, err)
	assert.True(t, len(appsList.Items) > 0)
}

func Test_AddAppForEnvironment(t *testing.T) {
	appsService := v1alpha.AppsService{}
	app := &v1alpha.App{}
	app.Name = "TestApp"
	app.Namespace = "TestNS"
	app.Version = "1.0"
	app.Description = "Description"
	app.Phase = "Phase"
	app.Repository = "Repo1"
	app.Values = []string{"v1", "v2"}
	appsService.AddAppForEnvironment(app, "git@github.com:jstrachan/environment-jstrachan-jx5-dev.git")
}

func Test_RemoveAppFromEnvironment(t *testing.T) {
	appsService := v1alpha.AppsService{}
	app := &v1alpha.App{}
	app.Name = "stable/nginx-ingress"
	appsService.RemoveAppFromEnvironment(app, "git@github.com:jstrachan/environment-jstrachan-jx5-dev.git")
}
