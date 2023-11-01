# EnvLoader - CLI Tool 
EnvLoader is a powerful command-line interface (CLI) tool designed to enhance the management of environment variables in your projects. With the ability to encrypt and decrypt files using a maximum 32-bit key, EnvLoader also provides seamless injection of environment variables from .ini files into binary files.

## Features
- File Encryption and Decryption: Securely encrypt and decrypt files using a key of up to 32 bits to protect sensitive information.

- Environment Variable Injection into Binary Files: Easily load environment variables from .ini files and inject them into binary files, suitable for Docker Compose and shell scripts.

## Usage
Using **key prompt**:
```
envLoader encrypt -f config.ini
```

```
envLoader decrypt -f config.ini
```

Using **key string**:  
```
-k [your secure key]
```

```
--key [your secure key]
```

Using **key file**:  
```
--key-file [your file with secure key]
```

Using **output to a different file**:  
```
-o [file]
```

```
--output [file]
```

## Loading variables
EnvLoader allows you to seamlessly load environment variables from both encrypted and decrypted .ini configuration files into shell scripts. This is particularly useful for managing configurations across different environments such as development, testing, and production.

Assuming simple _config.ini_ file and binary script _run.sh_:
``` ini
message="Hello everyone!"

[github]
name=MyFancyNameOnGithub

[gitlab]
name=MyfancyNameOnGitlab
```

``` sh
#!/bin/sh

echo "Running! $MESSAGE" >> test.txt
echo "Running! $NAME" >> test.txt
```

### Encrypted configuration
Load Variables from **encrypted** .ini configuration:
```
envLoader -f config.ini -b run.sh -k [your secure key] -e github
```
In this example, the -e flag is used to choose the section (github in this case), and EnvLoader will include variables from both the global section and the specified section in the generated script.

### Decrypted configuration
Load Variables from **decrypted** .ini configuration:
```
envLoader -f config.ini -b run.sh -e gitlab
```

## Try It Yourself!

Clone the repository, build the application or download it.

```
make build
```

Place the binary into your **$PATH**.

### Encrypt it!

Encrypt the configuration file with provided password file and load it's variables into the binary scripts.

```
cd mock/  
envLoader encrypt -f config.ini --key-file password.txt  
envLoader load -f config.ini -b run.sh --key-file password.txt -e github
```

Check the created **test.txt**:

```
cat test.txt
```
