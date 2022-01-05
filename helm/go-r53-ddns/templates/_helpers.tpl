{{/*
Expand the name of the chart.
*/}}
{{- define "go-r53-ddns.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 50 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "go-r53-ddns.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 50 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "go-r53-ddns.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "go-r53-ddns.labels" -}}
helm.sh/chart: {{ include "go-r53-ddns.chart" . }}
{{ include "go-r53-ddns.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "go-r53-ddns.selectorLabels" -}}
app.kubernetes.io/name: {{ include "go-r53-ddns.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "go-r53-ddns.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "go-r53-ddns.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create value key-pairs for Secrets in the JobTemplate
*/}}
{{- define "go-r53-ddns.envFromSecretGen" -}}
{{- range $_,$value := .Values.cronConf.envFromSecret }}
- name: {{ $value.key | upper | quote }}
  valueFrom:
    secretKeyRef:
      name: {{ required (cat "Secret Reference is require for " $value.key ) $value.from.secret | quote }}
      key: {{ required (cat "Key is require for " $value.key ) $value.from.key | quote }}
{{- end }}
{{- end }}

{{/*
Create value key-pairs for environment variables in the JobTemplate
*/}}
{{- define "go-r53-ddns.envGen" -}}
{{- range $_,$value := . }}
- name: {{ $value.key | upper | quote }}
  value: {{ required "All environment values must be set." $value.value | quote }}
{{- end }}
{{- end }}