package diff

import (
	"strings"
)

func CheckCI() CIDetected {
	return checkCI(&systemEnvReader{})
}

func checkCI(reader EnvReader) CIDetected {
	// Appveyor
	// https://www.appveyor.com/docs/environment-variables/
	// Travis
	// https://docs.travis-ci.com/user/environment-variables/#default-environment-variables
	if lookupBuildEnvironment(reader, Appveyor) == "true" {
		return Detected
	}

	// Jenkins
	// https://wiki.jenkins.io/display/JENKINS/Building+a+software+project#Buildingasoftwareproject-belowJenkinsSetEnvironmentVariables
	if len(lookupBuildEnvironment(reader, Jenkins)) > 0 {
		return Detected
	}

	// GitHub Action
	// https://help.github.com/en/actions/automating-your-workflow-with-github-actions/using-environment-variables#default-environment-variables
	if len(lookupBuildEnvironment(reader, GitHub)) > 0 {
		return Detected
	}

	// AzureDevops
	// https://docs.microsoft.com/en-us/azure/devops/pipelines/build/variables?view=azure-devops&tabs=yaml#agent-variables
	// Variable Name is 'Agent.Id' to detect if this is a Azure Pipelines agent.
	// Note that variables are upper-cased and '.' is replaced with '_' on Azure Pipelines.
	// https://docs.microsoft.com/en-us/azure/devops/pipelines/process/variables?view=azure-devops&tabs=yaml%2Cbatch#access-variables-through-the-environment
	if len(lookupBuildEnvironment(reader, AzureDevOps)) > 0 {
		return Detected
	}

	// TeamCity
	// https://www.jetbrains.com/help/teamcity/predefined-build-parameters.html#PredefinedBuildParameters-ServerBuildProperties
	if len(lookupBuildEnvironment(reader, TeamCity)) > 0 {
		return Detected
	}

	// MyGet
	// https://docs.myget.org/docs/reference/build-services#Available_Environment_Variables
	if lookupBuildEnvironment(reader, MyGet) == "myget" {
		return Detected
	}

	// GitLab
	// https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
	if len(lookupBuildEnvironment(reader, GitLab)) > 0 {
		return Detected
	}

	return NotDetected
}

func lookupBuildEnvironment(reader EnvReader, env BuildServerEnvironment) string {
	var value string
	envValue, found := reader.LookupEnv(string(env))
	if found {
		value = strings.ToLower(envValue)
	}
	return value
}
