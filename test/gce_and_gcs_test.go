package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestComputeInstance(t *testing.T) {
	t.Parallel()

	// input variables
	instanceName := fmt.Sprintf("tf-instance-%s", strings.ToLower(random.UniqueId()))
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	randomRegion := gcp.GetRandomRegion(t, projectId, nil, nil)
	randomZone := gcp.GetRandomZoneForRegion(t, projectId, randomRegion)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Where your TF code lives
		TerraformDir: "../terraform/",

		Vars: map[string]interface{}{
			"instance_name": instanceName,
			"project_id":    projectId,
			"region":        randomRegion,
			"zone":          randomZone,
		},
		EnvVars: map[string]string{
			"GOOGLE_CLOUD_PROJECT": projectId,
		},
	})

	// terraform init, terraform apply, terraform destroy
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Test: SSH into the new instance
	publicIp := terraform.Output(t, terraformOptions, "public_ip")
	// Variables to create keyPair
	sshUsername := "hmh-user"
	keySize := 2048

	instance := gcp.FetchInstance(t, projectId, instanceName)
	keyPair := ssh.GenerateRSAKeyPair(t, keySize) // returns private and public key
	instance.AddSshKey(t, sshUsername, keyPair.PublicKey)

	host := ssh.Host{
		Hostname:    publicIp,
		SshKeyPair:  keyPair,
		SshUserName: sshUsername,
	}

	// Retry loop to allow changes to load
	maxRetries := 20
	sleepBetweenRetries := 3 * time.Second
	outputText := "HMH Demo"

	retry.DoWithRetry(t, "Attempting to SSH", maxRetries, sleepBetweenRetries, func() (string, error) {
		output, err := ssh.CheckSshCommandE(t, host, fmt.Sprintf("echo '%s'", outputText))
		if err != nil {
			return "", err
		}

		if strings.TrimSpace(outputText) != strings.TrimSpace(output) {
			return "", fmt.Errorf("Expected: %s. Got: %s\n", outputText, output)
		}
		return "", nil
	})
}

func TestStorageBucket(t *testing.T) {
	t.Parallel()

	// input variables
	bucketName := fmt.Sprintf("tf-bucket-%s", strings.ToLower(random.UniqueId()))
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)
	randomRegion := gcp.GetRandomRegion(t, projectId, nil, nil)
	randomZone := gcp.GetRandomZoneForRegion(t, projectId, randomRegion)

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Where your TF code lives
		TerraformDir: "../terraform/",

		Vars: map[string]interface{}{
			"bucket_name": bucketName,
			"project_id":  projectId,
			"region":      randomRegion,
			"zone":        randomZone,
		},
		EnvVars: map[string]string{
			"GOOGLE_CLOUD_PROJECT": projectId,
		},
	})

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Test: Check that the bucket exists
	gcp.AssertStorageBucketExists(t, bucketName)

	// Test: Check that the bucket URL matches
	testUrl := fmt.Sprintf("gs://%s", bucketName)
	bucketUrl := terraform.Output(t, terraformOptions, "bucket_url")

	assert.Equal(t, bucketUrl, testUrl)
}
