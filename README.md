# VaultZ


```
.vaultz/dextra/self.priv
.vaultz/dextra/self.pub

.vaultz/dextra/keyring/raoni.pubkey
.vaultz/dextra/keyring/lucas.pubkey
.vaultz/dextra/keyring/alvaro.pubkey
.vaultz/dextra/keyring/ed.pubkey
.vaultz/dextra/keyring/murer.pubkey

.vaultz/dextra/padlock/dsa/writer/murer.id
.vaultz/dextra/padlock/dsa/writer/fernando.id
.vaultz/dextra/padlock/dsa/reader/leo.id
.vaultz/dextra/padlock/dsa/reader/lucas.id
.vaultz/dextra/padlock/dsa/reader/alvaro.id
```

```shell
murer: vaultz dextra enc  test.txt            -o test.txt.enc.vaultz       -r ed lucas alvaro raoni murer leo -l 3
ed:    vaultz dextra req  test.txt.enc.vaultz -o ed.request.padlock.vaultz -r ed lucas alvaro raoni murer leo
ed:    vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o ed.pin.padlock.vaultz
raoni: vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o raoni.pin.padlock.vaultz
murer: vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o murer.pin.padlock.vaultz
ed:    vaultz dextra dec  test.txt.enc.vaultz -p ed.pin.padlock.vaultz raoni.pin.padlock.vaultz murer.pin.padlock.vaultz
ed:    vaultz dextra dec  test.txt.enc.vaultz -p                       raoni.pin.padlock.vaultz murer.pin.padlock.vaultz -o test.txt

murer: vaultz dextra enc  test.txt            -o test.txt.enc.vaultz       -r ed lucas alvaro raoni murer leo    # -l 1
ed:    vaultz dextra dec  test.txt.enc.vaultz -o test.txt

murer: vaultz dextra enc  test.txt            -o test.txt.enc.vaultz       -g dsa    -l 3
ed:    vaultz dextra req  test.txt.enc.vaultz -o ed.request.padlock.vaultz -g dsa
ed:    vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o ed.pin.padlock.vaultz    -r ed
raoni: vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o raoni.pin.padlock.vaultz # confirm user interactive
murer: vaultz dextra resp test.txt.enc.vaultz ed.request.unlock.vaultz -o murer.pin.padlock.vaultz -r ed
ed:    vaultz dextra dec  test.txt.enc.vaultz -p ed.pin.padlock.vaultz raoni.pin.padlock.vaultz murer.pin.padlock.vaultz
ed:    vaultz dextra dec  test.txt.enc.vaultz -p                       raoni.pin.padlock.vaultz murer.pin.padlock.vaultz -o test.txt

murer: vaultz dextra enc  test.txt            -o test.txt.enc.vaultz       -g dsa    # -l 1
ed:    vaultz dextra dec  test.txt.enc.vaultz -o test.txt

vaultz dextra self gen -u murer
vaultz dextra self export pub -o murer.pubkey.vaultz
vaultz dextra self export priv -o murer.privkey.vaultz

vaultz dextra keyring import -i alvaro.pubkey.vaultz # confirm user interactive
vaultz dextra keyring delete alvaro

vaultz dextra padlock dsa create
vaultz dextra padlock dsa delete
vaultz dextra padlock dsa writer/read add/remove murer ed
vaultz dextra padlock dsa writer/read list
vaultz dextra padlock dsa export -o dsa.ids.padlock.vaultz # signed by murer
vaultz dextra padlock dsa import -i dsa.ids.padlock.vaultz # confirm user interactive

vaultz dextra export -o dextra.namespace.vaultz
vaultz dextra import -i dextra.namespace.vaultz
```