# VaultZ

```yaml
file:

- signed: author
- encrypt: none
- obj:

    - header: 
        - amount: 1
        - encrypt: none
        - obj:
            - keyHash: Salted hash of final key
            - recipients: Number of recipients
            - lock: Lock size

    - Pins:
        - amount: header.obj.recipients
        - encrypt: to the recipient
        - obj: recipient pin

    - Locks:
        - amount: combination (header.obj.recipients, header.obj.lock)
        - encrypt: none
        - obj: pgp sym enc, key = (pins[a].obj + pins[b].obj + pins[c].obj), plain = final key

    - Payload:
        - amount: 1
        - encrypt: symmetric with final key
        - obj: payload

```