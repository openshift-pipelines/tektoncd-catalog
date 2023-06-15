## Workspaces

| Workspace      | Optional                           | Description                |
| :------------- | :--------------------------------: | :------------------------- |
{{- range .workspaces }}
| `{{ .name }}`  | `{{ .optional | formatOptional }}` | {{ .description | chomp }} |
{{- end }}

## Params

| Param         | Type                       | Default                      | Description                |
| :------------ | :------------------------: | :--------------------------- | :------------------------- |
{{- range .params }}
| `{{ .name }}` | `{{ .type | formatType }}` | {{ .default | formatValue }} | {{ .description | chomp }} |
{{- end }}

## Results

| Result        | Description                |
| :------------ | :------------------------- |
{{- range .results }}
| `{{ .name }}` | {{ .description | chomp }} |
{{- end }}
