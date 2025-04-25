package kubernetes

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/spf13/viper"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/logger"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/types"
	v1types "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type notifyServiceChange func(watch.EventType, *v1types.Service, *types.Metadata)
type notifyServiceList func(*[]types.Metadata)

var log = logger.GetLogger()

func discoverPlugins(namespace string, notifier notifyServiceChange, listnotifier notifyServiceList) error {

	plugins, err := client.CoreV1().Services(namespace).List(context.Background(), v1.ListOptions{LabelSelector: types.PluginLabelSelector})
	if err == nil {
		metadataList := make([]types.Metadata, 0)
		for _, svc := range plugins.Items {
			metadata, err := extractServiceMetdata(&svc)
			if err != nil {
				log.Sugar().Errorln(err)
				continue
			}
			metadataList = append(metadataList, *metadata)
			if err != nil {
				log.Sugar().Errorln(err)
				continue
			}
			notifier(watch.Added, &svc, metadata)
		}
		listnotifier(&metadataList)
		return nil
	}
	return err
}

func extractServiceMetdata(service *v1types.Service) (*types.Metadata, error) {

	annotation := service.ObjectMeta.Annotations[types.PluginMetadataAnnotation]

	if annotation == "" {
		log.Sugar().Errorw("Plugin found, but there was no annotation for metadata", service)
	}
	var metadata *types.Metadata
	err := json.Unmarshal([]byte(annotation), &metadata)
	if err != nil {
		log.Sugar().Error(err)
		return nil, err
	}

	if metadata.Name == "" || metadata.Route == "" || metadata.Version == "" || metadata.ServiceGuid == "" || metadata.RouteGuid == "" {
		log.Sugar().Error(err)
		return nil, errors.New("Metadata has missing values.")
	}

	return metadata, err
}

func StartPluginObserver(context context.Context, notifier notifyServiceChange, listnotifier notifyServiceList) error {
	namespace := viper.GetString("NAMESPACE")
	err := discoverPlugins(namespace, notifier, listnotifier)

	if err != nil {
		return err
	}

	// Start Live Event Watcher
	timeOut := int64(60)
	watcher, err := client.CoreV1().Services(namespace).Watch(context, v1.ListOptions{TimeoutSeconds: &timeOut})
	if err != nil {
		log.Sugar().Errorln(err)
	}

	for event := range watcher.ResultChan() {
		svc := event.Object.(*v1types.Service)
		metadata, err := extractServiceMetdata(svc)

		if err != nil && metadata == nil {
			log.Sugar().Errorln(err)
			continue
		}

		notifier(event.Type, svc, metadata)
	}

	return nil
}
