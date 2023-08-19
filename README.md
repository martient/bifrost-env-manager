# Bifrost-env-manager

## Logic

If the variable is random or custom then it will be write one time and never update. If the variable is static then we will update it each time.

For example the database password is random and will be store and not edited but the target api version is static so can be change by the user in time.

## Key value definition

### High level config

version is not required but higly recommanded, by default 1.0.0

software_target is not required but higly recommanded

filename is not required but recommanded, by default .env

### Static variables

key -> value, pretty simple.

### Random value variables

The key with all is params

required:
    - key
    - lenght

optional:
    - as_upper_case, default True
    - as_lower_case, default True
    - as_diggit, default True
    - as_special_character, default True
    - available_character, default
        - if asUpperCase "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        - if asLowerCase "abcdefghijklmnopqrstuvwxyz"
        - if asDigit 0123456789"
        - if asSpecialCharacter "!#$%^&*()-_=+[]{}|;:,.<>/?"

### Custom value variables

The key with is line who need to be built

Required:
    - key
    - line
    - values array

    if value between {{ }} is not define here we will check in random value, static value and env. if the value still not found, the CLI will raise an error.
