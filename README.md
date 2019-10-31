# pomf-go

Basic implementation of Pomf.

**Note:** Only the `/upload` path is implemented so far, and it doesn't
handle errors.

## Install

```sh
go get -u github.com/warkio/pomf-go/cmd/pomf-server

cd "${GOPATH}/bin"

# Copy the `pomf-server` executable to some directory in your $PATH.
#
# This assumes you want to put it in `${HOME}/bin`.
cp pomf-server "${HOME}/bin"
```

## Usage

### Start server

```sh
export POMF_LISTEN_ADDRESS='127.0.0.1:8000'

# The resulting file URL will just concatenate this with the file name,
# so you might want to make sure there's a trailing slash.
export POMF_URL_PREFIX='https://example.com/'

# Where to put the files.
export POMF_UPLOAD_DIRECTORY='./files/'

pomf-server
```

### Upload file

```sh
curl -i -sSF 'files[]=@waifu.png' 'https://127.0.0.1:8000/upload'
```
