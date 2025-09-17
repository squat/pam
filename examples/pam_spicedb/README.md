# pam_spicedb.so

`pam_spicedb.so` is a PAM module that validates a user's account by checking a permission in SpiceDB.

## Usage

[embedmd]:# (help.txt)
```txt
Usage of pam_spicedb:
  -endpoint string
    	The SpiceDB URL
  -insecure-skip-verify
    	Whether to skip TLS verification
  -permission string
    	SpiceDB permission
  -resource-id string
    	SpiceDB resource ID
  -resource-type string
    	SpiceDB resource type
  -subject-type string
    	SpiceDB subject type
  -tls
    	Whether to enable TLS (default true)
  -token-file string
    	Path to a file containing the SpiceDB token
```

## Example

```pam
account required /usr/lib/security/pam_spicedb.so --tls=false --endpoint=localhost:50051 --token-file=/var/lib/spicedb/token --permission=ssh --resource-id=foo --resource-type=server --subject-type=user
```
