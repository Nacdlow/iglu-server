package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"gitlab.com/group-nacdlow/nacdlow-server/models"
	macaron "gopkg.in/macaron.v1"
)

// PluginManifest is used to describe the plugin's id, name, author, version, etc.
type PluginManifest struct {
	Id, Name, Author, Version string
}

// DeviceRegistration represents a device to be registered to a plugin. This is
// used to inform a plugin about a device.
type DeviceRegistration struct {
	DeviceID    int
	Description string
	Type        models.DeviceType
}

// ExtensionType specifies the type of web extension.
type ExtensionType int

const (
	CSS = iota + 1
	JavaScript
)

// WebExtension represents an addon to the web page.
type WebExtension struct {
	Type           ExtensionType
	PathMatchRegex string
	Source         string
}

// Iglu is the interface that we're exposing as a plugin.
type Iglu interface {
	OnLoad()
	Middleware() macaron.Handler
	// TODO these are interfaces yet to be implemented
	GetManifest() PluginManifest
	RegisterDevice(reg DeviceRegistration) error
	OnDeviceToggle(id int, status bool) error
	GetWebExtensions() []WebExtension
}

// IgluRPC is what the server is using to communicate to the plugin over RPC
type IgluRPC struct{ client *rpc.Client }

func (i *IgluRPC) OnLoad() {
	err := i.client.Call("Plugin.OnLoad", new(interface{}), nil)
	if err != nil {
		panic(err)
	}
}

func (i *IgluRPC) Middleware() (handler macaron.Handler) {
	err := i.client.Call("Plugin.Middleware", new(interface{}), &handler)
	if err != nil {
		panic(err)
	}
	return
}

// IgluRPCServer is the RPC server which IgluRPC talks to.
type IgluRPCServer struct {
	Impl Iglu
}

func (s *IgluRPCServer) OnLoad(args interface{}) error {
	s.Impl.OnLoad()
	return nil
}

func (s *IgluRPCServer) Middleware(args interface{}) macaron.Handler {
	return s.Impl.Middleware()
}

// This is the implementation of plugin.Plugin.
type IgluPlugin struct {
	Impl Iglu
}

func (p *IgluPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &IgluRPCServer{Impl: p.Impl}, nil
}

func (IgluPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &IgluRPC{client: c}, nil
}
