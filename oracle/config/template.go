package config

const (
	DefaultConfigTemplate = `

###############################################################################
###                                  Oracle                                 ###
###############################################################################
[oracle]
## Update Interval (in seconds) is the time between each time the oracle triggers providers to update price-data
update_interval = "{{ .Oracle.UpdateInterval }}"

## Timeout is the time that the vote-extension handler will wait for a response from the oracle (either running in / out-of process), generally this parameter should be 
## less than the timeout_prevote parameter in the consensus config
timeout = "{{ .Oracle.Timeout }}"

## InProcess specifies whether the oracle configured, is currently running as a remote grpc-server, or will be run in process
in_process = {{ .Oracle.InProcess }}

## RemoveAddress is the address of the remote oracle grpc-server, only used if in_process is set to false
remote_address= "{{ .Oracle.RemoteAddress }}"

# Providers
{{- range $provider := .Oracle.Providers }}

[[oracle.providers]]
name = "{{ $provider.Name }}"
apikey = "{{ $provider.Apikey }}"
provider_timeout = "{{ $provider.ProviderTimeout }}"

# Token Name to TokenMetadata
[[oracle.providers.token_name_to_metadata]]
{{- range $key, $value := $provider.TokenNameToMetadata }}
{{ $key }} = {
    symbol = "{{ $value.Symbol }}"
    decimals = {{ $value.Decimals }}
    is_twap = {{ $value.IsTWAP }}
}
{{- end }}
{{- end }}

# Currency Pairs
{{- range $pair := .Oracle.CurrencyPairs }}

[[oracle.currency_pairs]]
base = "{{ $pair.Base }}"
quote = "{{ $pair.Quote }}"
{{- end }}

[oracle.metrics]
## PrometheusServerAddress is the address of the prometheus server that the oracle will expose metrics to
prometheus_server_address = "{{ .Oracle.Metrics.PrometheusServerAddress }}"

## Oracle Metrics
[oracle.metrics.oracle_metrics]
## This will enable metrics to be reported for the oracle side-car, this will only be enabled if the oracle is running in the configured process
enabled = {{ .Oracle.OracleMetrics.Enabled }}

## App Metrics
[oracle.metrics.app_metrics]
## This will enable metrics to be reported for the oracle client in the application, this should always be enabled if oracle metrics are desired
enabled = {{ .Oracle.AppMetrics.Enabled }}
validator_cons_address = "{{ .Oracle.AppMetrics.ValidatorConsAddress }}"
`
)
