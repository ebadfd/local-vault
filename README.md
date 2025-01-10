# Local Vault 

Local Vault is a command-line tool designed to securely store and manage sensitive environment variables and files. Using GPG encryption under the hood, 
it allows you to encrypt your sensitive data and access it when needed. This tool is ideal for managing environment variables across multiple projects, 
or any other files you want to encrypt, and can be used in combination with cloud storage for backups with trust.

The tool's structure is organized by project, and allows for flexible management of environment files.

## Directory Structure

The directory hierarchy follows this structure:

```
project
    └── app
        └── env
            └── file
```

Here, an **app** represents a repository, and each project can contain multiple applications with their own respective environment files.

## Installation and Setup

### Initialize the Vault

To get started, initialize the vault by running:

```bash
local-vault init
```

This will set up the necessary database schema and create a home folder for storing your encrypted files.

### Create a New Project

To create a new project, run:

```bash
local-vault project -n <project-name>
```

You'll be prompted to fill in details for the project. This includes providing a recipient email address, which will be used for decrypting the files.

## Managing Encrypted Files

### Import Environment Variables

To import environment variables into your project, use:

```bash
local-vault import -a <app-name> -e <environment> -p <project-name> <file-path>
```

Example:

```bash
local-vault import -a backend-api -e dev -p test .env
```

This will import the contents of the `.env` file into the specified project, environment, and profile.

### Dump Environment Variables to a File

To dump the stored environment variables into a file, run:

```bash
local-vault dump -a <app-name> -e <environment> -p <project-name> <file-path>
```

Example:

```bash
local-vault dump -a backend-api -e dev -p test .env
```

This will output the decrypted environment variables into the specified file.

## Additional Features

- **Archiving**: If you update or overwrite an existing record, the previous file is not deleted. It will be stored in an archive folder for potential recovery in the future.

## Notes

- The tool is currently being used to test how it functions, and your feedback on the experience will help guide further improvements.
- I will improve the tool further in the next release, still exploring the usability of the tool.
