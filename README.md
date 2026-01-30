# drone-env-inspector

- [drone-env-inspector](#drone-env-inspector)
  - [Synopsis](#synopsis)
  - [Parameters](#parameters)
  - [Notes](#notes)
  - [Plugin Image](#plugin-image)
  - [Examples](#examples)

## Synopsis

This plugin inspects environment variables and outputs their values for use in subsequent pipeline steps. It supports reading multiple environment variables at once and can optionally write to a secret output file for sensitive values.

To learn how to utilize Drone plugins in Harness CI, please consult the provided [documentation](https://developer.harness.io/docs/continuous-integration/use-ci/use-drone-plugins/run-a-drone-plugin-in-ci).

## Parameters

| Parameter | Choices/<span style="color:blue;">Defaults</span> | Comments |
| :--- | :--- | :--- |
| env_name <span style="font-size: 10px"><br/>`string`</span> <span style="color:red; font-size: 10px">`required`</span> | | Comma-separated list of environment variable names to inspect. |
| secret <span style="font-size: 10px"><br/>`boolean`</span> | Default: `false` | If true, writes output to HARNESS_OUTPUT_SECRET_FILE instead of DRONE_OUTPUT. Use this for sensitive values. |
| log_level <span style="font-size: 10px"><br/>`string`</span> | | Set to `debug` or `trace` for verbose logging. |

## Notes

- The plugin echoes the value of each specified environment variable to the console in `KEY=VALUE` format (e.g., `DRONE_OUTPUT=/path/to/output`).
- The plugin also writes the values to the output file for use in subsequent pipeline steps.
- If an environment variable does not exist, the plugin will log a warning and output an empty value for that variable.
- Multiple environment variables can be specified using comma-separated values (e.g., `VAR1,VAR2,VAR3`).
- Spaces around variable names are automatically trimmed.
- When `secret: true` is set, the output is written to `HARNESS_OUTPUT_SECRET_FILE` to prevent sensitive values from appearing in logs.

## Plugin Image

The plugin `plugins/env-inspector` is available for the following architectures:

| OS | Tag |
| --- | --- |
| latest | `linux-amd64/arm64` |
| linux/amd64 | `linux-amd64` |
| linux/arm64 | `linux-arm64` |

## Examples

```yaml
# Single environment variable
- step:
    type: Plugin
    name: inspect-env
    identifier: inspect_env
    spec:
      connectorRef: harness-docker-connector
      image: plugins/env-inspector
      settings:
        env_name: MY_VAR

# Multiple environment variables
- step:
    type: Plugin
    name: inspect-multiple-env
    identifier: inspect_multiple_env
    spec:
      connectorRef: harness-docker-connector
      image: plugins/env-inspector
      settings:
        env_name: VAR1,VAR2,VAR3

# Secret output (for sensitive values)
- step:
    type: Plugin
    name: inspect-secrets
    identifier: inspect_secrets
    spec:
      connectorRef: harness-docker-connector
      image: plugins/env-inspector
      settings:
        env_name: API_KEY,DB_PASSWORD
        secret: true

# With debug logging
- step:
    type: Plugin
    name: inspect-env-debug
    identifier: inspect_env_debug
    spec:
      connectorRef: harness-docker-connector
      image: plugins/env-inspector
      settings:
        env_name: MY_VAR
        log_level: debug

# Use output in subsequent step
- step:
    type: Run
    name: Use Env Value
    identifier: use_env_value
    spec:
      shell: Sh
      command: |
        echo "The value is: <+steps.inspect_env.output.outputVariables.MY_VAR>"
```

> <span style="font-size: 14px; margin-left:5px; background-color: #d3d3d3; padding: 4px; border-radius: 4px;">:information_source: If you notice any issues in this documentation, you can [edit this document](https://github.com/harness-community/drone-env-inspector/blob/main/README.md) to improve it.</span>

