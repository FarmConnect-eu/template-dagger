package main

// WithAwsBackend configure le backend S3 + credentials AWS (compatible MinIO)
func (m *Terraform) WithAwsBackend(
	accessKeyID string,
	secretAccessKey string,
	bucket string,
	key string,
	// +optional
	// +default="us-east-1"
	region string,
) *Terraform {
	newVariables := make([]Variable, len(m.Variables), len(m.Variables)+2)
	copy(newVariables, m.Variables)

	newVariables = append(newVariables,
		Variable{Key: "AWS_ACCESS_KEY_ID", Value: accessKeyID, IsSecret: true, TfVar: false},
		Variable{Key: "AWS_SECRET_ACCESS_KEY", Value: secretAccessKey, IsSecret: true, TfVar: false},
	)

	return &Terraform{
		Variables:        newVariables,
		State:            &StateConfig{Backend: "s3", Bucket: bucket, Key: key, Region: region},
		TerraformVersion: m.TerraformVersion,
	}
}
