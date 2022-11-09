package kubernetes

var tempDeployment = `---
kind: Deployment
apiVersion: apps/v1
metadata:
  name: {{.ServiceName}}
  namespace: pcs
  labels:
    service: pcs
    version: {{.Version}}
spec:
  replicas: {{.Replicas}}
  selector:
    matchLabels:
      name: {{.ServiceName}}
      service: pcs
  template:
    metadata:
      labels:
        name: {{.ServiceName}}
        service: pcs
    spec:
      imagePullSecrets:
        - name: regsecret
      containers:
        - name: {{.ServiceName}}
          image: {{.Image}}:{{.Version}}
          ports: {{range .Ports}}
            - name: {{.Name}}
              containerPort: {{.Port}}
              protocol: TCP{{end}}
            - name: service-port
              containerPort: 8080
              protocol: TCP
            - name: metrics
              containerPort: 9000
              protocol: TCP
            - name: pprof
              containerPort: 7777
              protocol: TCP{{if .EnvFromConfigMap}}
          envFrom:
            - configMapRef:
                name: {{.EnvFromConfigMap}}{{end}}
          env: {{range .Envs}}
            - name: {{.Name}}
              value: "{{.Value}}"{{end}}
            - name: SERVICE_ADDRESS_METRICS
              value: ':9000'
            - name: SERVICE_ADDRESS_PPROF
              value: ':7777'
            - name: PORT
              value: '8080'
---`
