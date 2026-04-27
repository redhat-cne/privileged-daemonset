package privilegeddaemonset

import (
	"sync"
	"time"
)

// Config holds optional overrides for package behavior. When passed to SetConfig,
// fields set to their zero value keep the built-in default from DefaultConfig.
type Config struct {
	// RoleServiceAccountName is the name of the Role, RoleBinding, and ServiceAccount used for the privileged DaemonSet and RBAC. Default: "privileged-ds".
	RoleServiceAccountName string
	// WaitPollInterval is the sleep between polls when waiting for DaemonSet deletion or readiness. Default: 5s.
	WaitPollInterval time.Duration
	// NamespaceDeleteTimeout is how long to wait for a namespace to be removed after Delete. Default: 2m.
	NamespaceDeleteTimeout time.Duration
	// NamespaceDeletionPollInterval is the poll interval in namespaceWaitForDeletion. Default: 1s.
	NamespaceDeletionPollInterval time.Duration
	// TolerationPeriodSeconds is tolerationSeconds for node.kubernetes.io not-ready and unreachable taints. Default: 300.
	TolerationPeriodSeconds int64
	// DaemonSetDeleteTimeout is how long to wait for a DaemonSet to disappear after delete. Default: 5m.
	DaemonSetDeleteTimeout time.Duration
}

// DefaultConfig returns the package defaults matching the previous hard-coded constants.
func DefaultConfig() Config {
	return Config{
		RoleServiceAccountName:        "privileged-ds",
		WaitPollInterval:              5 * time.Second,
		NamespaceDeleteTimeout:        2 * time.Minute,
		NamespaceDeletionPollInterval: time.Second,
		TolerationPeriodSeconds:       300,
		DaemonSetDeleteTimeout:        5 * time.Minute,
	}
}

var (
	pkgCfgMu sync.RWMutex
	pkgCfg   = DefaultConfig()
)

// SetConfig merges override into the defaults from DefaultConfig: each field in override that is
// the zero value for its type is replaced by the default. Call before other package functions if
// you need non-default timeouts or names.
func SetConfig(override Config) {
	pkgCfgMu.Lock()
	defer pkgCfgMu.Unlock()
	pkgCfg = mergeConfig(override)
}

func mergeConfig(override Config) Config {
	b := DefaultConfig()
	if override.RoleServiceAccountName != "" {
		b.RoleServiceAccountName = override.RoleServiceAccountName
	}
	if override.WaitPollInterval > 0 {
		b.WaitPollInterval = override.WaitPollInterval
	}
	if override.NamespaceDeleteTimeout > 0 {
		b.NamespaceDeleteTimeout = override.NamespaceDeleteTimeout
	}
	if override.NamespaceDeletionPollInterval > 0 {
		b.NamespaceDeletionPollInterval = override.NamespaceDeletionPollInterval
	}
	if override.TolerationPeriodSeconds > 0 {
		b.TolerationPeriodSeconds = override.TolerationPeriodSeconds
	}
	if override.DaemonSetDeleteTimeout > 0 {
		b.DaemonSetDeleteTimeout = override.DaemonSetDeleteTimeout
	}
	return b
}

func getConfig() Config {
	pkgCfgMu.RLock()
	defer pkgCfgMu.RUnlock()
	return pkgCfg
}
