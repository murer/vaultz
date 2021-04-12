# VaultZ


```
.vaultz/dextra/self.keypair

.vaultz/dextra/keyring/raoni.pubkey
.vaultz/dextra/keyring/lucas.pubkey
.vaultz/dextra/keyring/alvaro.pubkey
.vaultz/dextra/keyring/ed.pubkey
.vaultz/dextra/keyring/murer.pubkey

.vaultz/dextra/safe/dsa/writer/murer.id
.vaultz/dextra/safe/dsa/writer/fernando.id

.vaultz/dextra/safe/dsa/reader/leo.id
.vaultz/dextra/safe/dsa/reader/lucas.id
.vaultz/dextra/safe/dsa/reader/alvaro.id
```

```shell
vaultz dextra enc test.txt -r lucas alvaro -o test.txt.vaultz
vaultz dextra dec test.txt.vaultz -o test.txt

vaultz dextra safe dsa enc test.txt -o test.txt.vaultz 
vaultz dextra safe dsa dec test.txt.vaultz -o test.txt
```