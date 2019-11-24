module kube-job-runner

go 1.12

require (
	github.com/aws/aws-sdk-go v1.23.8
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/lib/pq v1.2.0
	github.com/rs/zerolog v1.15.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.3.0
	k8s.io/api v0.0.0-20190820101039-d651a1528133
	k8s.io/apimachinery v0.0.0-20190823012420-8ca64af22337
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
)
