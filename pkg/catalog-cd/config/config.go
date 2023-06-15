package config

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	tkncli "github.com/tektoncd/cli/pkg/cli"
)

type Config struct {
	Stream         *tkncli.Stream
	kubeConfigPath string
	kubeContext    string
	namespace      string
	tp             *tkncli.TektonParams
}

func (c *Config) Infof(format string, a ...any) {
	if _, err := fmt.Fprintf(c.Stream.Out, format, a...); err != nil {
		panic(err)
	}
}

func (c *Config) Errorf(format string, a ...any) {
	if _, err := fmt.Fprintf(c.Stream.Err, format, a...); err != nil {
		panic(err)
	}
}

func (c *Config) GetTektonParams() *tkncli.TektonParams {
	if c.tp != nil {
		return c.tp
	}

	c.tp = &tkncli.TektonParams{}
	if c.kubeConfigPath != "" {
		c.tp.SetKubeConfigPath(c.kubeConfigPath)
	}
	if c.kubeContext != "" {
		c.tp.SetKubeContext(c.kubeContext)
	}
	if c.namespace != "" {
		c.tp.SetNamespace(c.namespace)
	}
	return c.tp
}

func (c *Config) GetNamespace() string {
	return c.GetTektonParams().Namespace()
}

func (c *Config) GetClientsOrPanic() *tkncli.Clients {
	cs, err := c.GetTektonParams().Clients()
	if err != nil {
		panic(err)
	}
	return cs
}

// NewConfigWithFlags sets up a new Config instance adding command-line flags to set its attributes.
func NewConfigWithFlags(stream *tkncli.Stream, flags *pflag.FlagSet) *Config {
	cfg := &Config{Stream: stream}

	flags.StringVarP(&cfg.kubeConfigPath, "kubeconfig", "k", cfg.kubeConfigPath,
		"kubectl config file")
	flags.StringVarP(&cfg.kubeContext, "context", "c", cfg.kubeContext,
		"kubernetes context name")
	flags.StringVarP(&cfg.namespace, "namespace", "n", cfg.namespace,
		"kubernetes namespace name")
	return cfg
}

// NewConfig sets up a new Config using the default STDIN, STDOUT and STDERR.
func NewConfig() *Config {
	return &Config{Stream: &tkncli.Stream{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}}
}
