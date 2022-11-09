package kubernetes

var tempService = `---
kind: Service
apiVersion: v1
metadata:
  name: {{.ServiceName}}
  namespace: pcs
  labels:
    name: {{.ServiceName}}
    service: pcs
spec:
  ports: {{range .Ports}}
    - name: {{.Name}}
      targetPort: {{.Port}}
      port: {{.Port}}
      protocol: TCP{{end}}
    - name: service-port
      protocol: TCP
      port: 8080
      targetPort: 8080
  selector:
    name: {{.ServiceName}}
    service: pcs
  type: ClusterIP
---`
