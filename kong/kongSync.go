package kong

import (
	"fmt"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/common"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/logger"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/types"
	v1types "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

const protocol = "http"

var log = logger.GetLogger()

func SyncKongServices(services *[]types.Metadata) {
	adminApi := viper.GetString("KONG_ADMIN_API")

	tags, err := common.ViperGetStringSlice("PLUGIN_TAGS")
	if err != nil {
		log.Sugar().Errorln(err)
		return
	}
	kr, err := kongListRoutes(adminApi, "", tags)

	if err != nil || kr == nil {
		log.Sugar().Error(err)
		return
	}

	for _, r := range kr {
		route := r.([]interface{})

		for _, ro := range route {
			data := ro.(map[string]interface{})

			routeId := data["id"].(string)
			serviceId := data["service"].(map[string]interface{})["id"].(string)
			if routeId != "" && serviceId != "" {
				var found = false
				for _, svc := range *services {
					if svc.Route == routeId || svc.ServiceGuid == serviceId {
						found = true
					}
				}

				if !found {
					err := kongDeleteRoute(routeId, adminApi)
					if err == nil {
						err = kongDeleteService(serviceId, adminApi)
					}
					if err != nil {
						log.Sugar().Errorln("Deletion not successfull.")
					}
				}
			}
		}
	}
}

func SyncKongService(event watch.EventType, svc *v1types.Service, metadata *types.Metadata) {
	adminApi := viper.GetString("KONG_ADMIN_API")
	str := []string{metadata.Name}
	name := strings.ReplaceAll(strings.Join(str, "-"), " ", "-")

	service, err := kongListService(adminApi, metadata.ServiceGuid)

	tags, err := common.ViperGetStringSlice("PLUGIN_TAGS")
	if err != nil {
		log.Sugar().Errorln(err)
		return
	}
	methods, err := common.ViperGetStringSlice("PLUGIN_HTTP_METHODS")
	if err != nil {
		log.Sugar().Errorln(err)
		return
	}
	port, err := getPort(svc)
	if err != nil {
		log.Sugar().Errorln(err)
		return
	}

	host := svc.Name + "." + svc.Namespace + ".svc.cluster.local"

	if err == nil {
		if len(service) > 0 {

			if event == watch.Deleted {
				err = kongDeleteRoute(metadata.RouteGuid, adminApi)
				if err == nil {
					err = kongDeleteService(metadata.ServiceGuid, adminApi)
				}
				if err != nil {
					log.Sugar().Errorln("Deletion not successful.")
					return
				}
			}

			if event == watch.Modified || event == watch.Added {
				err = kongCreateService(metadata.ServiceGuid, name, protocol, host, "", port, adminApi, "PATCH", tags)
				if err == nil {
					err = kongCreateRoute(metadata.ServiceGuid, metadata.RouteGuid, name, metadata.Route, adminApi, "PATCH", tags, methods)
				}
				if err != nil {
					log.Sugar().Errorln(err, "Modification not successful.")
					return
				}
			}

		} else {
			err := kongCreateService(metadata.ServiceGuid, name, protocol, host, "", port, adminApi, "POST", tags)

			if err == nil {
				err = kongCreateRoute(metadata.ServiceGuid, metadata.RouteGuid, name, metadata.Route, adminApi, "POST", tags, methods)
			}
			if err != nil {
				log.Sugar().Errorln("Creation not successful.")
				return
			}
		}
	}
}

func getPort(svc *v1types.Service) (string, error) {
	postfix := viper.GetString("PLUGIN_PORT_NAME_POSTFIX")
	for _, port := range svc.Spec.Ports {
		if strings.HasSuffix(port.Name, postfix) {
			return strconv.Itoa(int(port.Port)), nil
		}

	}
	// If no port with the postfix is found, return an error
	err := fmt.Errorf("port with postfix `%s` not found for service: %s", postfix, svc.Name)
	return "", err
}
