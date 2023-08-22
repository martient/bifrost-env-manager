# 🌉 Bifrost-env-manager

## 📚 Overview

Bifrost-env-manager is a tool for managing environment variables with flexibility. It distinguishes between static, random, and custom variables, allowing you to maintain a coherent environment setup.

## 🧠 Logic

Environment variables are classified into three categories: static, random, and custom. This classification determines how variables are treated over time.

- **Static variables:** These are straightforward key-value pairs that remain unchanged.
- **Random value variables:** These variables are generated once and stored securely.
- **Custom value variables:** These variables are constructed using predefined rules and can include references to other variables.

## 🔑 Key Value Definitions

### High-level Configuration

- `version` (optional): Specifies the version of the configuration (default: 1.0.0).
- `software_target` (optional): Indicates the target software for the environment setup.
- `filename` (optional): Specifies the name of the environment file (default: .env).

### Static Variables

Key-value pairs that remain constant.

### Random Value Variables

These variables are generated and stored securely.

Required:

- `key`
- `length`

Optional:

- `as_upper_case` (default: true)
- `as_lower_case` (default: true)
- `as_digit` (default: true)
- `as_special_character` (default: true)
- `available_characters` (default: specific to the chosen options)

### Custom Value Variables

Variables with dynamic values constructed from predefined rules.

Required:

- `key`
- `line`
- `values` (array)

If a value enclosed in {{ }} is not defined here, the tool will search in random values, static values, and environment variables. An error will be raised if the value is not found.

## 🚀 Getting Started

### Usage

``` shell
❯ bifrost-env-manager
Software environement files manager

Usage:
  Bifrost-env-manager [command]

Available Commands:
  generate    Generate a new version of the env file
  help        Help about any command

Flags:
      --config string          config file for this software environement
      --disable-update-check   Disable auto update checking before execution
  -h, --help                   help for Bifrost-env-manager
      --path string            Path for the new file folder, ex: /home/ubuntu/code/
  -t, --toggle                 Help message for toggle

Use "Bifrost-env-manager [command] --help" for more information about a command.
```

#### Command

##### Basic

``` shell
❯ bifrost-env-manager generate
Using config file: config.json
.example.env file generated successfully!
```

##### Specific config

``` shell
❯ bifrost-env-manager generate --path /Users/martient/Documents/bifrost-env-manager/bin
Using config file: config.json
/somewhere/else/example.env file generated successfully!
```

##### Specific file path

``` json
❯ bifrost-env-manager generate
Using config file: config.json
.example.env file generated successfully!
```

### Config

Config example:

``` json
{
  "version": "1.0.0",
  "software_target": "example",
  "filename": ".example.env",
  "static_variables": [
    {
      "smpt": "x.e.f.x"
    },
    {
      "s3Secret": "dhsd;asdkas;dkasdasda"
    },
    {
      "somethingElse": "wead"
    }
  ],
  "random_value_variables": [
    {
      "key": "password_psql",
      "lenght": 128,
      "as_upper_case": true,
      "as_lower_case": false,
      "as_diggit": true,
      "as_special_character": true,
      "available_character": "qwertyuiopasfdghjkl;zxcvbnm123234567890.,?!"
    }
  ],
  "custom_value_variables": [
    {
      "key": "postgresql_url",
      "line": "postgres://{{ db_user_name }}:{{ password_psql }}@{{ host }}/{{ db_name }}",
      "values": [
        {
          "db_user_name": "xxx"
        },
        {
          "db_name": "xxx"
        }
      ]
    }
  ]
}
```

Or can be

``` json
{
  "static_variables": [
    {
      "smpt": "x.e.f.x"
    },
    {
      "s3Secret": "dhsd;asdkas;dkasdasda"
    },
    {
      "somethingElse": "wead"
    }
  ],
  "random_value_variables": [
    {
      "key": "password_psql",
    }
  ],
  "custom_value_variables": [
    {
      "key": "postgresql_url",
      "line": "postgres://{{ db_user_name }}:{{ password_psql }}@{{ host }}/{{ db_name }}",
      "values": [
        {
          "db_user_name": "xxx"
        },
        {
          "db_name": "xxx"
        }
      ]
    }
  ]
}
```

## 🤝 Contributing

🎉 We welcome contributions from everyone! Here's how you can contribute:

- Fork the repository.
- Create a new branch for your feature/fix.
- Make your changes and commit them.
- Push your changes to your forked repository.
- Create a pull request to the 'develop' branch of the main repository.

See the [CONTRIBUTING](CONTRIBUTING) file for more details.

## 📝 License

Bifrost-env-manager is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0). See the [LICENSE](LICENSE) file for more details.
