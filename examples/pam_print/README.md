# pam_print.so

`pam_print.so` is a PAM module that prints the user, flags, arguments, and environment variables provided to the module.
The `--output` flag can be provided in the service configuration file to specify if the output should be printed as human-readable text (`--output=text`) or JSON (`--output=json`).

## Example

```pam
account required /usr/lib/security/pam_print.so --output=json foo bar
auth required /usr/lib/security/pam_print.so --output=text foo bar
```
