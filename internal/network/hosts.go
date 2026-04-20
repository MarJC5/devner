// Package network manages /etc/hosts entries and Caddy site configuration
// for devner-managed domains.
//
// RFC 6761 reserves *.localhost for 127.0.0.1 and modern OS resolvers
// (macOS, Linux, Windows 10+) honor this natively. Devner therefore DOES
// NOT write *.localhost entries to /etc/hosts — only custom TLDs (.dev,
// .test, etc.) require it.
package network

import (
	"fmt"
	"strings"

	"github.com/devner/devner/internal/platform"
	"github.com/txn2/txeh"
)

const marker = "devner managed"

type HostsManager struct {
	hosts *txeh.Hosts
}

func NewHostsManager() (*HostsManager, error) {
	h, err := txeh.NewHosts(&txeh.HostsConfig{
		ReadFilePath:  platform.HostsFilePath(),
		WriteFilePath: platform.HostsFilePath(),
	})
	if err != nil {
		return nil, fmt.Errorf("open hosts file: %w", err)
	}
	return &HostsManager{hosts: h}, nil
}

// NeedsHostsEntry returns false for *.localhost domains (RFC 6761) and true
// otherwise. Callers should skip hosts mutation when this returns false.
func NeedsHostsEntry(domain string) bool {
	return !strings.HasSuffix(strings.ToLower(domain), ".localhost")
}

// Add inserts a host entry if the domain is not a .localhost TLD.
// Returns (added bool, err). added=false means no-op (either already present
// or .localhost domain which resolves natively).
func (m *HostsManager) Add(domain, target string) (bool, error) {
	if !NeedsHostsEntry(domain) {
		return false, nil
	}
	if target == "" {
		target = "127.0.0.1"
	}
	if found, _, _ := m.hosts.HostAddressLookup(domain, txeh.IPFamilyV4); found {
		return false, nil
	}
	m.hosts.AddHost(target, domain)
	return true, nil
}

func (m *HostsManager) Remove(domain string) (bool, error) {
	if !NeedsHostsEntry(domain) {
		return false, nil
	}
	if found, _, _ := m.hosts.HostAddressLookup(domain, txeh.IPFamilyV4); !found {
		return false, nil
	}
	m.hosts.RemoveHost(domain)
	return true, nil
}

func (m *HostsManager) Save() error {
	if err := m.hosts.Save(); err != nil {
		if !platform.IsElevated() {
			return fmt.Errorf("save hosts file (need elevated privileges): %w", err)
		}
		return fmt.Errorf("save hosts file: %w", err)
	}
	return nil
}

// Marker is emitted as a comment around devner entries so we can distinguish
// our entries from others. txeh does not support block comments; we add
// marker lines via WriteHostFile separately when rebuilding from scratch.
func Marker() string { return marker }
