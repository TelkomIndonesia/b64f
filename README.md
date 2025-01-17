# B64F

Base64 Encode/Decode multiple files.

## How to use

1. create `.b64f` file under working directory with list of files that should be encoded to base64 file. Each entry support [glob pattern](https://github.com/gobwas/glob), for example:

    ```.
    pkg/util/testdata/*
    **/*.yml
    ```

1. run the following command to encode the listed files. This will create a new sibling file for each listed file with additional `.b64` extension, e.g. `pkg/util/testdata/a.key` will produce `pkg/util/testdata/a.key.b64`. Note that it will overwride the `.b64` file if it already existed.

    ```bash
    go run -mod=mod github.com/telkomindonesia/b64f encode
    ```

1. to decode back the `.b64` files, run the following command on the same directory. Note that it will overwride the decoded file if it already existed.

    ```bash
    go run -mod=mod github.com/telkomindonesia/b64f decode
    ```

### Run Methods

Instead of using `go run` you can download pre-compiled binary from [Github releases page](https://github.com/TelkomIndonesia/b64f/releases).

Another alternative is to use docker, e.g.

```bash
docker run --rm \
    -v "$PWD:$PWD" -w "$PWD" \
    ghcr.io/telkomindonesia/b64f:latest \
    [decode|encode]
```

## Tips

You can use `makefie` to make the command simpler, for example:

```makefile
b64f: mode ?= "decode"
b64f:
    go run -mod=mod github.com/telkomindonesia/b64f $(mode)        
```

which can then be used using `make b64f` or `make b64f mode=encode`

see [testdata](./testdata/) folder for example.
