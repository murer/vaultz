# VaultZ

```yaml
file:

- header: 
    - encrypt: none
    - data: |
        xxxxx

- pins:

    - comment: pins must be sorted by KeyId

    - pin1:
        - encrypt: user 1 public key
        - data: generated pin 1

    - pin2:
        - encrypt: user 2 public key
        - data: generated pin 2
    
    - pinN:
        - encrypt: user N public key
        - data: generated pin 3


```