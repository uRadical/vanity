# vanity
Generate the HTML necessary to host Go modules and programs on vanity domains.

## Installation

### MacOS
```shell
brew tap uradical/tap
brew install vanity
```


## Usage
One two mandatory arguments are required to specify the package import path and the URL to
the repository containing the code. A third optional argument specifies the output directory
where the generated `index.html` will be created, if this is omitted the file will be written
to the current working directory.

```
vaniety example.com/pkg/mylib https://github.com/example-ord/mylib
```

To specify a directory to write `index.html` to use:

```
vaniety \
    example.com/pkg/mylib \
    https://github.com/example-org/mylib \
   ~/mysite/public/pkg/mylib
```
