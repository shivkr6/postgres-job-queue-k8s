{{/*
Return the chart name.
*/}}
{{- define "postgres-job-queue.name" -}}
{{- .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return a release-scoped full name for future resources that need it.
Current K8 resources intentionally keep their existing Kubernetes names.
*/}}
{{- define "postgres-job-queue.fullname" -}}
{{- printf "%s-%s" .Release.Name (include "postgres-job-queue.name" .) | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Return the chart label value.
*/}}
{{- define "postgres-job-queue.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels for Helm-managed objects.
*/}}
{{- define "postgres-job-queue.labels" -}}
helm.sh/chart: {{ include "postgres-job-queue.chart" . }}
app.kubernetes.io/name: {{ include "postgres-job-queue.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Stable selector labels for the Postgres Pod.
These preserve the raw YAML's app=postgres routing behavior.
*/}}
{{- define "postgres-job-queue.postgresSelectorLabels" -}}
app: postgres
{{- end -}}

{{/*
Stable selector labels for the migration Job Pod.
*/}}
{{- define "postgres-job-queue.migrateSelectorLabels" -}}
app: queue-migrate
{{- end -}}
