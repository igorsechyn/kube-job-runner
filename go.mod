module kube-job-runner

go 1.12

require (
	github.com/aws/aws-sdk-go v1.23.8
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/lib/pq v1.2.0
	github.com/rs/zerolog v1.15.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	k8s.io/api v0.17.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
)
