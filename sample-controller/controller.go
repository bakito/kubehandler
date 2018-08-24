package main

import (
	"log"
	"time"

	"github.com/bakito/kubehandler"
	"github.com/bakito/kubehandler/pkg/signals"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type LoggerHandler struct {
	kubehandler.DefaultHandler
}

func (l *LoggerHandler) GetName() string {
	return "LoggerHandler"
}

func (l *LoggerHandler) AddFunc(namespace string, name string) error {
	log.Printf("Pod %s added in namespace %s", name, namespace)
	return nil
}

func (l *LoggerHandler) UpdateFunc(namespace string, name string) error {
	log.Printf("Pod %s updated in namespace %s", name, namespace)
	return nil
}

func (l *LoggerHandler) DeleteFunc(namespace string, name string) error {
	log.Printf("Pod %s deleted in namespace %s", name, namespace)
	return nil
}

func main() {

	// Ignoring some errors for brevity
	cfg, _ := clientcmd.BuildConfigFromFlags("", "")
	kubeClient, _ := kubernetes.NewForConfig(cfg)

	// Get a pod informer
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
	informer := kubeInformerFactory.Core().V1().Pods().Informer()

	loggingHandler := &LoggerHandler{
		DefaultHandler: kubehandler.DefaultHandler{
			Informer: informer,
			Synced:   informer.HasSynced,
		},
	}

	// This name is used as the workqueue name
	loop := kubehandler.NewEventLoop("logger_queue")

	// Register all your handlers
	loop.Register(loggingHandler)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	// Start the k8s informer so you get events
	go kubeInformerFactory.Start(stopCh)

	// Start processing events. This can run in a go routine if you want to
	// continue doing something else.
	loop.Run(2, stopCh)
}
