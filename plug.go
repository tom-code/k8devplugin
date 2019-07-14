package main

import (
	"fmt"
  "net"
  "time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pluginapi "k8s.io/kubernetes/pkg/kubelet/apis/deviceplugin/v1beta1"
	registerapi "k8s.io/kubernetes/pkg/kubelet/apis/pluginregistration/v1"
)

const (
	socket = "/var/lib/kubelet/plugins_registry/my.sock"
)

type regServer struct {
}

func (rs regServer) GetInfo(ctx context.Context, rqt *registerapi.InfoRequest) (*registerapi.PluginInfo, error) {
	info := &registerapi.PluginInfo{
		Type:              registerapi.DevicePlugin,
		Name:              "coolant/water",
		Endpoint:          socket,
		SupportedVersions: []string{"v1alpha1", "v1beta1"},
	}
	fmt.Println("get info invoked!")
	return info, nil
}

func (rs regServer) NotifyRegistrationStatus(ctx context.Context, regstat *registerapi.RegistrationStatus) (*registerapi.RegistrationStatusResponse, error) {
	if regstat.PluginRegistered {
		fmt.Println("registered!")
	} else {
		fmt.Printf("reg error %s\n", regstat.Error)
	}
	return &registerapi.RegistrationStatusResponse{}, nil
}

type plugServer struct {
}


func (rs plugServer) Allocate(ctx context.Context, rqt *pluginapi.AllocateRequest) (*pluginapi.AllocateResponse, error) {
  fmt.Println("allcate!")
	resp := new(pluginapi.AllocateResponse)
	return resp, nil
}

func (rs plugServer) GetDevicePluginOptions(ctx context.Context, empty *pluginapi.Empty) (*pluginapi.DevicePluginOptions, error) {
  fmt.Println("get options!")
	return &pluginapi.DevicePluginOptions{
		PreStartRequired: false,
	}, nil
}

func (rs plugServer) ListAndWatch(emtpy *pluginapi.Empty, stream pluginapi.DevicePlugin_ListAndWatchServer) error {
  fmt.Println("list")

  resp := new(pluginapi.ListAndWatchResponse)

  resp.Devices = []*pluginapi.Device{
    &pluginapi.Device{ID: "unique_id-00", Health: pluginapi.Healthy},
    &pluginapi.Device{ID: "unique_id-01", Health: pluginapi.Healthy},
    &pluginapi.Device{ID: "unique_id-02", Health: pluginapi.Unhealthy},
  }

  stream.Send(resp)

	for {
		time.Sleep(1*time.Second)
	}
}

func (rs plugServer) PreStartContainer(ctx context.Context, psRqt *pluginapi.PreStartContainerRequest) (*pluginapi.PreStartContainerResponse, error) {
  fmt.Println("prestartcontainer")
	return &pluginapi.PreStartContainerResponse{}, nil
}

func main() {
	fmt.Println("start")

	grpcsrv := grpc.NewServer()
	registrator := regServer{}
	plugin := plugServer{}

	registerapi.RegisterRegistrationServer(grpcsrv, registrator)
	pluginapi.RegisterDevicePluginServer(grpcsrv, plugin)

	listener, err := net.Listen("unix", socket)
	if err != nil {
		panic(err)
	}
	grpcsrv.Serve(listener)
}
