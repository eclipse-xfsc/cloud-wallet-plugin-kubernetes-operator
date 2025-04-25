package main

import (
	"context"

	"github.com/spf13/viper"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/kong"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/kubernetes"
	"github.com/eclipse-xfsc/cloud-wallet-plugin-kubernetes-operator/logger"
)

func main() {
	ctx, f := context.WithCancel(context.Background())

	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	viper.AutomaticEnv()

	logger := logger.GetLogger()
	defer logger.Sync()

	err := kubernetes.InitializeKubernetes()

	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = kubernetes.StartPluginObserver(ctx, kong.SyncKongService, kong.SyncKongServices)

	if err != nil {
		logger.Error(err.Error())
		return
	}

	defer f()

}
