# Piknik

Forked personal version of [piknik](https://github.com/jedisct1/piknik) with codepaths I do not use removed and support for client integration with the OSX clipboard.

## License

[ISC](https://en.wikipedia.org/wiki/ISC_license).

## Credits

Upstream [Piknik](https://github.com/jedisct1/piknik) repository.

## Protocol, re-written for clarity (to me)

Common definitions:

```text
v: protocol version
Hk,s: BLAKE2b(domain="SK", key=k, salt=s, size=32)
r: random 256-bit client nonce
r': random 256-bit server nonce
```

Initial exchange:

```text
h0 = Hk,0(v || r)
-> v || r || h0

h1 = Hk,1(v || r' || h0)
<- v || r' || h1

```

```text
ek: 256-bit symmetric encryption key
ekid: encryption key id encoded as a 64-bit little endian integer

Len(x): x encoded as a 64-bit little endian unsigned integer

ct: XChaCha20 ek,n (plaintext message)

n: random 192-bit nonce
s = Ed25519(ekid || n || ct)

ts: Unix timestamp as a 64-bit little endian integer

Copy:

h2 = Hk,2(h1 || opcode || ts || s)
-> opcode || h2 || Len(ekid || n || ct) || ts || s || ekid || n || ct
<- Hk,3(h2)

Move/Paste:

h2 = Hk,2(h1 || opcode)
-> opcode || h2
<- Hk,3(h2 || ts || s) || Len(ekid || n || ct) || ts || s || ekid || n || ct
```
