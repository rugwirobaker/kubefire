package config

import (
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Cluster struct {
	Name         string `json:"name"`
	Bootstrapper string `json:"bootstrapper"`
	Pubkey       string `json:"pubkey"`
	Prikey       string `json:"prikey"`

	Image       string `json:"image"`
	KernelImage string `json:"kernel_image,omitempty"`
	KernelArgs  string `json:"kernel_args,omitempty"`

	Admin  Node `json:"admin"`
	Master Node `json:"master"`
	Worker Node `json:"worker"`

	ExtraOptions string `json:"extra_options"`

	Deployed bool `json:"deployed"` // the only status property
}

func NewCluster() *Cluster {
	c := Cluster{Admin: Node{}, Master: Node{}, Worker: Node{}}

	c.Admin.Cluster = &c
	c.Master.Cluster = &c
	c.Worker.Cluster = &c

	return &c
}

type Node struct {
	Count    int      `json:"count"`
	Memory   string   `json:"memory,omitempty"`
	Cpus     int      `json:"cpus,omitempty"`
	DiskSize string   `json:"disk_size,omitempty"`
	Cluster  *Cluster `json:"-"`
}

func (c *Cluster) ParseExtraOptions(obj interface{}) interface{} {
	value := reflect.ValueOf(obj).Elem()

	optionList := strings.Split(c.ExtraOptions, ",")

	for _, option := range optionList {
		values := strings.SplitN(option, "=", 2)
		if len(values) == 2 {
			field := value.FieldByName(values[0])

			switch field.Kind() {
			case reflect.String:
				pattern := regexp.MustCompile(`^["'](.+)["']$`)
				values[1] = pattern.ReplaceAllString(values[1], "$1")

				field.SetString(values[1])

			case reflect.Int:
				v, _ := strconv.Atoi(values[1])
				field.SetInt(int64(v))

			case reflect.Bool:
				b, _ := strconv.ParseBool(values[1])
				field.SetBool(b)
			}
		}
	}

	return value.Interface()
}

func (c *Cluster) LocalClusterDir() string {
	return path.Join(ClusterRootDir, c.Name)
}

func (c *Cluster) LocalKubeConfig() string {
	return path.Join(c.LocalClusterDir(), "admin.conf")
}

func (c *Cluster) LocalClusterConfigFile() string {
	return path.Join(c.LocalClusterDir(), "cluster.yaml")
}

func (c *Cluster) LocalClusterKeyFiles() (string, string) {
	return path.Join(c.LocalClusterDir(), "key"), path.Join(c.LocalClusterDir(), "key.pub")
}
