package kubernetes

import (
	"flag"
	"fmt"
	"gitlab.com/dipper-iot/shared/util"
	"html/template"
	"log"
	"os"
)

func SetTemplateDeployment(temp string) {
	tempDeployment = temp
}

func Gen(services []*Service) {
	var version string
	flag.StringVar(&version, "v", "latest", "version build")
	flag.Parse()

	if version == "" {
		version = os.Getenv("IMAGE_VERSION")
	}

	tD, err := template.New("deployment template").Parse(string(tempDeployment)) // Create a template.
	if err != nil {
		log.Print(err)
		return
	}

	tS, err := template.New("service template").Parse(string(tempService)) // Create a template.
	if err != nil {
		log.Print(err)
		return
	}

	util.CreateFolder("./kubernetes/services/")

	for _, service := range services {
		service.Version = version

		if service.Replicas == 0 {
			service.Replicas = 1
		}

		f, err := os.Create(fmt.Sprintf("./kubernetes/services/%s-deployment.yaml", service.ServiceName))
		if err != nil {
			log.Println("create file: ", err)
			return
		}
		err = tD.Execute(f, &service)
		if err != nil {
			log.Print("Can't execute ", err)
		}
		f.Close()

		if service.Service {

			f, err = os.Create(fmt.Sprintf("./kubernetes/services/%s-svc.yaml", service.ServiceName))
			if err != nil {
				log.Println("create file: ", err)
				return
			}
			err = tS.Execute(f, &service)
			if err != nil {
				log.Print("Can't execute ", err)
			}
			f.Close()
		}

	}
}
