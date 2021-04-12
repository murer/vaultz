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

vaultz dextra enc test.txt -r lucas alvaro ed raoni -l 3 -o test.txt.vaultz
ed falha: vaultz dextra dec test.txt.vaultz -o test.txt

# raoni: vaultz dextra part test.txt.vaultz -r ed -o raoni.part
# alvaro: vaultz dextra part test.txt.vaultz -r ed -o alvaro.part
# ed: vaultz dextra part test.txt.vaultz -r ed -o ed.part
# ed success: vaultz dextra dec test.txt.vaultz -p ed.part raoni.part alvaro.part -o test.txt

vaultz dextra dec test.txt.vaultz -request -r raoni alvaro lucas ed -o request.txt
ed, raoni, lucas: vaultz dextra part test.txt.vaultz -request request.txt -o response.txt
vaultz dextra dec test.txt.vaultz -response raoni.resp alvaro.resp ed.resp -o test.txt

vaultz dextra safe dsa enc test.txt -o test.txt.vaultz 
vaultz dextra safe dsa dec test.txt.vaultz -o test.txt

vaultz dextra keyring import dextra.keyring
vaultz dextra self export


```

```
a := generate AES256 key
encrypt payload AES256 key

encrypt(a, b)
encrypt(a, c)
encrypt(a, d)
```