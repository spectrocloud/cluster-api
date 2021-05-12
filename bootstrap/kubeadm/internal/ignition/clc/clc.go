/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package clc generates Ignition using Container Linux Config Transpiler.
package clc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	ignition "github.com/coreos/ignition/config/v2_3"
	ignitionTypes "github.com/coreos/ignition/config/v2_3/types"
	clct "github.com/flatcar-linux/container-linux-config-transpiler/config"
	"github.com/pkg/errors"

	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha4"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/internal/cloudinit"
)

const (
	clcTemplate = `---
{{- if .Users }}
passwd:
  users:
    {{- range .Users }}
    - name: {{ .Name }}
      {{- if .Gecos }}
      gecos: {{ .Gecos }}
      {{- end }}
      {{- if .Groups }}
      groups:
        {{- range Split .Groups ", " }}
        - {{ . }}
        {{- end }}
      {{- end }}
      {{- if .HomeDir }}
      home_dir: {{ .HomeDir }}
      {{- end }}
      {{- if .Shell }}
      shell: {{ .Shell }}
      {{- end }}
      {{- if .Passwd }}
      password_hash: {{ .Passwd }}
      {{- end }}
      {{- if .PrimaryGroup }}
      primary_group: {{ .PrimaryGroup }}
      {{- end }}
      {{- if .SSHAuthorizedKeys }}
      ssh_authorized_keys:
        {{- range .SSHAuthorizedKeys }}
        - {{ . }}
        {{- end }}
      {{- end }}
    {{- end }}
{{- end }}
systemd:
  units:
    - name: kubeadm.service
      enabled: true
      contents: |
        [Unit]
        Description=kubeadm
        # Run only once. After successful run, this file is moved to /tmp/.
        ConditionPathExists=/etc/kubeadm.yml
        [Service]
        # To not restart the unit when it exits, as it is expected.
        Type=oneshot
        ExecStart=/etc/kubeadm.sh
        [Install]
        WantedBy=multi-user.target
    {{- if .NTP }}{{ if .NTP.Enabled }}
    - name: ntpd.service
      enabled: true
    {{- end }}{{- end }}
    {{- range .Mounts }}
    {{- $label := index . 0 }}
    {{- $mountpoint := index . 1 }}
    {{- $disk := index $.FilesystemDevicesByLabel $label }}
    {{- $mountOptions := slice . 2 }}
    - name: {{ $mountpoint | MountpointName }}.mount
      enabled: true
      contents: |
        [Unit]
        Description = Mount {{ $label }}

        [Mount]
        What={{ $disk }}
        Where={{ $mountpoint }}
        Options={{ Join $mountOptions "," }}

        [Install]
        WantedBy=multi-user.target
    {{- end }}
storage:
  {{- if .DiskSetup }}{{- if .DiskSetup.Partitions }}
  disks:
    {{- range .DiskSetup.Partitions }}
    - device: {{ .Device }}
      {{- if .Overwrite }}
      wipe_table: {{ .Overwrite }}
      {{- end }}
      {{- if .Layout }}
      partitions:
      - {}
      {{- end }}
    {{- end }}
  {{- end }}{{- end }}
  {{- if .DiskSetup }}{{- if .DiskSetup.Filesystems }}
  filesystems:
    {{- range .DiskSetup.Filesystems }}
    - name: {{ .Label }}
      mount:
        device: {{ .Device }}
        format: {{ .Filesystem }}
        label: {{ .Label }}
        {{- if .ExtraOpts }}
        options:
          {{- range .ExtraOpts }}
          - {{ . }}
          {{- end }}
        {{- end }}
    {{- end }}
  {{- end }}{{- end }}
  files:
    {{- range .Users }}
    {{- if .Sudo }}
    - path: /etc/sudoers.d/{{ .Name }}
      mode: 0600
      contents:
        inline: |
          {{ .Name }} {{ .Sudo }}
    {{- end }}
    {{- end }}
    {{- if .UsersWithPasswordAuth }}
    - path: /etc/ssh/sshd_config
      mode: 0600
      contents:
        inline: |
          # Use most defaults for sshd configuration.
          Subsystem sftp internal-sftp
          ClientAliveInterval 180
          UseDNS no
          UsePAM yes
          PrintLastLog no # handled by PAM
          PrintMotd no # handled by PAM

          Match User {{ .UsersWithPasswordAuth }}
            PasswordAuthentication yes
    {{- end }}
    {{- range .WriteFiles }}
    - path: {{ .Path }}
      # Owner
      #
      # If Encoding == gzip+base64 || Encoding == gzip
      # compression: true
      #
      # If Encoding == gzip+base64 || Encoding == "base64"
      # Put "!!binary" notation before the content to let YAML decoder treat data as
      # base64 data.
      #
      {{ if ne .Permissions "" -}}
      mode: {{ .Permissions }}
      {{ end -}}
      contents:
        inline: |
          {{ .Content | Indent 10 }}
    {{- end }}
    - path: /etc/kubeadm.sh
      mode: 0700
      contents:
        inline: |
          #!/bin/bash
          set -e
          {{ range .PreKubeadmCommands }}
          {{ . }}
          {{- end }}

          {{ .KubeadmCommand }}
          mkdir -p /run/cluster-api && echo success > /run/cluster-api/bootstrap-success.complete
          mv /etc/kubeadm.yml /tmp/
          {{range .PostKubeadmCommands }}
          {{ . }}
          {{- end }}
    - path: /etc/kubeadm.yml
      mode: 0600
      contents:
        inline: |
          ---
          {{ .KubeadmConfig | Indent 10 }}
    {{- if .NTP }}{{- if and .NTP.Enabled .NTP.Servers }}
    - path: /etc/ntp.conf
      mode: 0644
      contents:
        inline: |
          # Common pool
          {{- range  .NTP.Servers }}
          server {{ . }}
          {{- end }}

          # Warning: Using default NTP settings will leave your NTP
          # server accessible to all hosts on the Internet.

          # If you want to deny all machines (including your own)
          # from accessing the NTP server, uncomment:
          #restrict default ignore

          # Default configuration:
          # - Allow only time queries, at a limited rate, sending KoD when in excess.
          # - Allow all local queries (IPv4, IPv6)
          restrict default nomodify nopeer noquery notrap limited kod
          restrict 127.0.0.1
          restrict [::1]
    {{- end }}{{- end }}
`
)

type render struct {
	*cloudinit.BaseUserData

	KubeadmConfig            string
	UsersWithPasswordAuth    string
	FilesystemDevicesByLabel map[string]string
}

func defaultTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"Indent":         templateYAMLIndent,
		"Split":          strings.Split,
		"Join":           strings.Join,
		"MountpointName": mountpointName,
	}
}

func mountpointName(name string) string {
	return strings.TrimPrefix(strings.ReplaceAll(name, "/", "-"), "-")
}

func templateYAMLIndent(i int, input string) string {
	split := strings.Split(input, "\n")
	ident := "\n" + strings.Repeat(" ", i)
	return strings.Join(split, ident)
}

func renderCLC(input *cloudinit.BaseUserData, kubeadmConfig string) ([]byte, error) {
	if input == nil {
		return nil, errors.New("empty base user data")
	}

	t := template.Must(template.New("template").Funcs(defaultTemplateFuncMap()).Parse(clcTemplate))

	usersWithPasswordAuth := []string{}
	for _, user := range input.Users {
		if user.LockPassword != nil && !*user.LockPassword {
			usersWithPasswordAuth = append(usersWithPasswordAuth, user.Name)
		}
	}

	filesystemDevicesByLabel := map[string]string{}
	if input.DiskSetup != nil {
		for _, filesystem := range input.DiskSetup.Filesystems {
			filesystemDevicesByLabel[filesystem.Label] = filesystem.Device
		}
	}

	data := render{
		BaseUserData:             input,
		KubeadmConfig:            kubeadmConfig,
		UsersWithPasswordAuth:    strings.Join(usersWithPasswordAuth, ","),
		FilesystemDevicesByLabel: filesystemDevicesByLabel,
	}

	var out bytes.Buffer
	if err := t.Execute(&out, data); err != nil {
		return nil, errors.Wrapf(err, "failed to render template")
	}

	return out.Bytes(), nil
}

func Render(input *cloudinit.BaseUserData, clc *bootstrapv1.ContainerLinuxConfig, kubeadmConfig string) ([]byte, string, error) {
	if clc == nil {
		return nil, "", errors.New("get empty CLC config")
	}

	clcBytes, err := renderCLC(input, kubeadmConfig)
	if err != nil {
		return nil, "", errors.Wrapf(err, "rendering CLC configuration")
	}

	userData, warnings, err := buildIgnitionConfig(clcBytes, clc)
	if err != nil {
		return nil, "", errors.Wrapf(err, "building Ignition config")
	}

	return userData, warnings, nil
}

func buildIgnitionConfig(baseCLC []byte, clc *bootstrapv1.ContainerLinuxConfig) ([]byte, string, error) {
	// We control baseCLC config, so treat it as strict.
	ign, _, err := clcToIgnition(baseCLC, true)
	if err != nil {
		return nil, "", errors.Wrapf(err, "converting generated CLC to Ignition")
	}

	var clcWarnings string

	if clc.AdditionalConfig != "" {
		additionalIgn, warnings, err := clcToIgnition([]byte(clc.AdditionalConfig), clc.Strict)
		if err != nil {
			return nil, "", errors.Wrapf(err, "converting additional CLC to Ignition")
		}

		clcWarnings = warnings

		ign = ignition.Append(ign, additionalIgn)
	}

	userData, err := json.Marshal(&ign)
	if err != nil {
		return nil, "", errors.Wrapf(err, "marshaling generated Ignition config into JSON")
	}

	return userData, clcWarnings, nil
}

func clcToIgnition(data []byte, strict bool) (ignitionTypes.Config, string, error) {
	clc, ast, reports := clct.Parse(data)

	if (len(reports.Entries) > 0 && strict) || reports.IsFatal() {
		return ignitionTypes.Config{}, "", fmt.Errorf("error parsing Container Linux Config: %v", reports.String())
	}

	ign, report := clct.Convert(clc, "", ast)
	if (len(report.Entries) > 0 && strict) || report.IsFatal() {
		return ignitionTypes.Config{}, "", fmt.Errorf("error converting to Ignition: %v", report.String())
	}

	reports.Merge(report)

	return ign, reports.String(), nil
}
