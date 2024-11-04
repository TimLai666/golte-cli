# golte-cli

Golte CLI is a command-line tool for building and running [Golte](https://github.com/nichady/golte) projects.

## Installation

```bash
go get github.com/nichady/golte-cli
```

## Usage

### Initialize(create) a project

```bash
golte new <project-name>
```

### Build the project

```bash
golte build
```

### Run the project

```bash
golte run
```

### Notes

- The executable file will be placed in the `dist` directory.
- The executable file name will be the same as the project name.
- On Windows, the executable file will have a `.exe` suffix.
