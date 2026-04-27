package privilegeddaemonset

import (
	"sync"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	d := DefaultConfig()
	if d.RoleServiceAccountName != "privileged-ds" {
		t.Fatalf("RoleServiceAccountName: got %q", d.RoleServiceAccountName)
	}
	if d.WaitPollInterval != 5*time.Second {
		t.Fatalf("WaitPollInterval: got %v", d.WaitPollInterval)
	}
	if d.NamespaceDeleteTimeout != 2*time.Minute {
		t.Fatalf("NamespaceDeleteTimeout: got %v", d.NamespaceDeleteTimeout)
	}
	if d.NamespaceDeletionPollInterval != time.Second {
		t.Fatalf("NamespaceDeletionPollInterval: got %v", d.NamespaceDeletionPollInterval)
	}
	if d.TolerationPeriodSeconds != 300 {
		t.Fatalf("TolerationPeriodSeconds: got %d", d.TolerationPeriodSeconds)
	}
	if d.DaemonSetDeleteTimeout != 5*time.Minute {
		t.Fatalf("DaemonSetDeleteTimeout: got %v", d.DaemonSetDeleteTimeout)
	}
}

func TestMergeConfigPartialOverride(t *testing.T) {
	got := mergeConfig(Config{
		WaitPollInterval:        100 * time.Millisecond,
		TolerationPeriodSeconds: 120,
	})
	if got.RoleServiceAccountName != "privileged-ds" {
		t.Fatalf("RoleServiceAccountName should stay default, got %q", got.RoleServiceAccountName)
	}
	if got.WaitPollInterval != 100*time.Millisecond {
		t.Fatalf("WaitPollInterval: got %v", got.WaitPollInterval)
	}
	if got.TolerationPeriodSeconds != 120 {
		t.Fatalf("TolerationPeriodSeconds: got %d", got.TolerationPeriodSeconds)
	}
	if got.NamespaceDeleteTimeout != 2*time.Minute {
		t.Fatalf("NamespaceDeleteTimeout should stay default, got %v", got.NamespaceDeleteTimeout)
	}
}

func TestSetConfigConcurrentRead(t *testing.T) {
	SetConfig(DefaultConfig())
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = getConfig().WaitPollInterval
		}()
	}
	SetConfig(Config{WaitPollInterval: 10 * time.Millisecond})
	wg.Wait()
}
