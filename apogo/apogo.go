package apogo

import (
	"errors"
	"log"
	"os"
)

var (
	defaultApogo              = &Apogo{}
	defaultLogger ApogoLogger = log.New(os.Stderr, "", log.LstdFlags)
)

type Apogo struct {
	Client *Client
}

func NewApogo(conf *Conf) *Apogo {
	return &Apogo{NewClient(conf)}
}

func (m *Apogo) Start() error {
	return m.Client.Start()
}

func (m *Apogo) StartWithConfFile(name string) error {
	conf, err := NewConf(name)
	if err != nil {
		return err
	}
	return m.StartWithConf(conf)
}

func (m *Apogo) StartWithConf(conf *Conf) error {
	m.Client = NewClient(conf)

	return m.Client.Start()
}

func (m *Apogo) Stop() error {
	return m.Client.Stop()
}

func (m *Apogo) StartWatchUpdate() {
	ceChan := m.Client.WatchUpdate()

	go func() {
		for {
			ce := <-ceChan

			for _, ob := range m.Client.getObservers() {
				ob.HandleChangeEvent(ce)
			}
		}
	}()
}

func (m *Apogo) RegisterObserver(observer ChangeEventObserver) (recall func()) {
	m.Client.registerObserver(observer)
	return func() {
		m.Client.recallObserver(observer)
	}
}

func (m *Apogo) SubscribeToNamespaces(namespaces ...string) error {
	return m.Client.SubscribeToNamespaces(namespaces...)
}

func (m *Apogo) GetStringWithNamespace(namespace, key string) (string, bool) {
	return m.Client.GetStringWithNamespace(namespace, key)
}

func (m *Apogo) GetString(key string) (string, bool) {
	return m.Client.GetString(key)
}

func (m *Apogo) GetIntWithNamespace(namespace, key string) (int, bool) {
	return m.Client.GetIntWithNamespace(namespace, key)
}

func (m *Apogo) GetInt(key string) (int, bool) {
	return m.Client.GetInt(key)
}

func (m *Apogo) GetFloat64WithNamespace(namespace, key string) (float64, bool) {
	return m.Client.GetFloat64WithNamespace(namespace, key)
}

func (m *Apogo) GetFloat64(key string) (float64, bool) {
	return m.Client.GetFloat64(key)
}

func (m *Apogo) GetBoolWithNamespace(namespace, key string) (bool, bool) {
	return m.Client.GetBoolWithNamespace(namespace, key)
}

func (m *Apogo) GetBool(key string) (bool, bool) {
	return m.Client.GetBool(key)
}

func (m *Apogo) GetNameSpaceContent(namespace string) (string, bool) {
	return m.Client.GetNamespaceContent(namespace)
}

func (m *Apogo) GetAllKeys(namespace string) []string {
	return m.Client.GetAllKeys(namespace)
}

func (m *Apogo) GetReleaseKey(namespace string) (string, bool) {
	return m.Client.GetReleaseKey(namespace)
}

// Start Apogo [Deprecated]
func Start() error {
	if defaultApogo.Client == nil {
		return errors.New("please use StartWithConfFile")
	}
	return defaultApogo.Start()
}

// StartWithConfFile run Apogo with conf file
func StartWithConfFile(name string) error {
	conf, err := NewConf(name)
	if err != nil {
		return err
	}
	return StartWithConf(conf)
}

// StartWithConf run Apogo with Conf
func StartWithConf(conf *Conf) error {
	return defaultApogo.StartWithConf(conf)
}

// Stop sync config
func Stop() error {
	return defaultApogo.Stop()
}

// StartWatchUpdate starts an infinite loop reading changeEvent from update channel
//
//	and calls HandleChangeEvent method of all observers
func StartWatchUpdate() {
	defaultApogo.StartWatchUpdate()
}

// RegisterObserver registers an observer that will be notified when change event happens
func RegisterObserver(observer ChangeEventObserver) (recall func()) {
	return defaultApogo.RegisterObserver(observer)
}

// SubscribeToNamespaces fetch namespace config to local and subscribe to updates
func SubscribeToNamespaces(namespaces ...string) error {
	return defaultApogo.SubscribeToNamespaces(namespaces...)
}

// GetStringWithNamespace get value from given namespace
func GetStringWithNamespace(namespace, key string) (string, bool) {
	return defaultApogo.GetStringWithNamespace(namespace, key)
}

// GetString from default namespace
func GetString(key string) (string, bool) {
	return GetStringWithNamespace(defaultNamespace, key)
}

func GetIntWithNamespace(namespace, key string) (int, bool) {
	return defaultApogo.GetIntWithNamespace(namespace, key)
}

func GetInt(key string) (int, bool) {
	return defaultApogo.GetInt(key)
}

func GetFloat64WithNamespace(namespace, key string) (float64, bool) {
	return defaultApogo.GetFloat64WithNamespace(namespace, key)
}

func GetFloat64(key string) (float64, bool) {
	return defaultApogo.GetFloat64(key)
}

func GetBoolWithNamespace(namespace, key string) (bool, bool) {
	return defaultApogo.GetBoolWithNamespace(namespace, key)
}

func GetBool(key string) (bool, bool) {
	return defaultApogo.GetBool(key)
}

// GetNamespaceContent get contents of namespace
func GetNameSpaceContent(namespace string) (string, bool) {
	return defaultApogo.GetNameSpaceContent(namespace)
}

// GetAllKeys return all config keys in given namespace
func GetAllKeys(namespace string) []string {
	return defaultApogo.GetAllKeys(namespace)
}

// GetReleaseKey return release key for namespace
func GetReleaseKey(namespace string) (string, bool) {
	return defaultApogo.GetReleaseKey(namespace)
}

func SetLogger(logger ApogoLogger) {
	defaultLogger = logger
}
