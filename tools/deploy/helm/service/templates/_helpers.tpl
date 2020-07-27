{{- define "service-name" -}}
{{/* Convert home-automation-example.service_name to example-service-name by removing namespace */}}
{{- trimPrefix (printf "%v-" .Release.Namespace) .Release.Name | replace "." "-" | replace "_" "-" }}
{{- end -}}
