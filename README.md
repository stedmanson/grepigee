# Grepigee

Grepigee is a command-line tool for searching and managing Apigee configurations across different environments.

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
  - [Config File Location](#config-file-location)
  - [Creating/Editing the Config File](#creatingediting-the-config-file)
  - [Using Config Values](#using-config-values)
- [Usage](#usage)
  - [Authentication](#authentication)
  - [Search with regex](#search-with-regex)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Installation

[Add installation instructions here]

## Configuration

Grepigee uses a YAML configuration file to store default settings. You can create or edit this file to set default values for environment, organization, and other options.

### Config File Location

- Mac/Linux: `~/.grepigee.yaml`
- Windows: `%USERPROFILE%\.grepigee.yaml`

### Creating/Editing the Config File

1. Create a new file named `.grepigee.yaml` in your home directory.
2. Add your configuration settings:

```yaml
environment: prod
organisation: mycompany
```

3. Save the file.

You can also use the following command to save your current settings to the config file:

```
grepigee save-config
```


## Usage

### Authentication

To authenticate and set up your Apigee bearer token:

This command will provide instructions on how to set your APIGEE_BEARER_TOKEN. Follow the on-screen instructions to complete the authentication process.


### Search with regex

To search for a specific pattern in your Apigee configurations:

grepigee find -e <environment> -x <regex_pattern>

Example:

```
grepigee find -e dev -x "(?i)api-ecs"
```

This command searches for the case-insensitive pattern "api-ecs" in the "dev" environment.

Options:
- `-e, --env`: Specify the environment to search in (required)
- `-x, --regex`: The regular expression pattern to search for
- `-o, --org`: Specify the organization (optional if set in config)

### Examples

Search in a specific environment:
Copygrepigee find -e test -x "api-v1"

Use a saved environment from config:
Copygrepigee find -x "error-handler"

Save current settings to config:
Copygrepigee -e staging -o myorg save-config


### Contributing
[Add contribution guidelines here]

### License
[Add license information here]