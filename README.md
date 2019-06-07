# nutcracker

a simple shell parsing and execution engine

### Features

nutcracker can parse a subset of shell commands as detailed below:

#### Environment variables

```bash
echo $HOME ${ENVVAR:-default value}
```

#### Command substitution

```bash
echo $(cat file.txt)
```

#### String variable interpolation

```bash
"$(echo hello) $HOME ${ENVVAR:-default value}"
```
