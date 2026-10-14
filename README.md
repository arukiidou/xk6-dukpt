# xk6-dukpt

[![Go Reference](https://pkg.go.dev/badge/github.com/arukiidou/xk6-dukpt.svg)](https://pkg.go.dev/github.com/arukiidou/xk6-dukpt)

- xk6-dukpt is a [k6](https://github.com/grafana/k6) [extension](https://github.com/grafana/xk6) for DUKPT(Derived Unique Key Per Transaction) cryptography operations. port from [moov-io](https://pkg.go.dev/github.com/moov-io/dukpt)
- It is meant to be called from k6 scripts. Every API that returns an `ArrayBuffer` needs the k6 VU runtime,
  so it lives on the module instance rather than as a package-level Go function; use [moov-io/dukpt](https://pkg.go.dev/github.com/moov-io/dukpt) directly from Go programs.

# Important

As with moov-io/dukpt, this extension is mostly for validation and debugging purposes.

> [!CAUTION]
> Use this extension with great care.
> **Never use production BDKs or IPEKs.** 
> Create keys dedicated to k6 testing, or use well-known test keys (e.g. the ANSI X9.24 test BDK `0123456789ABCDEFFEDCBA9876543210`).
> 
> - Wherever possible, follow these best practices:
>   - Avoid using the BDK directly. Wherever possible, derive keys only from the IPEK.
>   - Protect keys with [`crypto`(WebCrypto)](https://grafana.com/docs/k6/latest/javascript-api/crypto/).
>   - Do not hardcode key protection keys; load them with [`k6/secrets`](https://grafana.com/docs/k6/latest/using-k6/secret-source/).
>   - Never write keys to logs.

See [examples/dukpt_rsa.ts](examples/dukpt_rsa.ts) for a script that loads the BDK via `k6/secrets` and wraps the derived IK with RSA-OAEP using WebCrypto.

## Project Status

> [!WARNING]
> This project is under active development.
> Breaking changes may still occur before a stable release.

## Requirements

- k6 2.3.0+
  - From this version on, you can combine with `Uint8Array.fromHex()` and `Uint8Array.fromBase64()`.

## How to Build

```bash
go install go.k6.io/xk6@latest
xk6 build --k6-version latest --os linux --cgo 0 --with github.com/arukiidou/xk6-dukpt@latest
# or
# go get -tool go.k6.io/xk6@latest
# go tool xk6 build --k6-version latest --os linux --cgo 0 --with github.com/arukiidou/xk6-dukpt@latest
```

# Example script

```typescript file=dukpt.ts
import { check } from "k6";
import { deriveCurrentTransactionKey } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  // this example uses well-known IPEK, but you should replace it with your own.
  example("6AC292FAA1315B4D858AB3A3D7D5933A", "FFFF9876543210E00001");
}

function example(ik: string, ksn: string) {
  const ckExpected = "042666B49184CFA368DE9628D0397BC9";
  const ckBase64Expected = "BCZmtJGEz6No3pYo0Dl7yQ==";

  const ck = deriveCurrentTransactionKey(Uint8Array.fromHex(ik), Uint8Array.fromHex(ksn));
  const ckBase64 = new Uint8Array(ck).toBase64();

  check(null, {
    'deriveCurrentTransactionKey(ik, ksn) as hex': () => new Uint8Array(ck).toHex().toUpperCase() === ckExpected,
    'deriveCurrentTransactionKey(ik, ksn) as base64': () => ckBase64 === ckBase64Expected,
    // Uint8Array.fromBase64() converts a base64 key back into a call argument.
    'Uint8Array.fromBase64(ck)': () => Uint8Array.fromBase64(ckBase64).toHex().toUpperCase() === ckExpected,
  });
}
```

## MAC generation

`generateMac` / `generateMacAsBase64` port [des.GenerateMac](https://pkg.go.dev/github.com/moov-io/dukpt/pkg/des#GenerateMac)
and return the **full 8 bytes** produced by moov-io. Per ANSI X9.24-1, the MAC is
the **first 4 bytes** of that result, so truncate before sending it:

```typescript
import { generateMacAsBase64 } from "k6/x/dukpt/des";

const data = new TextEncoder().encode("4012345678909D987").toBase64();
// data === "NDAxMjM0NTY3ODkwOUQ5ODc="
const mac = generateMacAsBase64(ck, data, "request");
// mac === "nMx4Fz/E+2Q="  (9CCC78173FC4FB64)
// the X9.24-1 mac is the first 4 bytes: 9CCC7817
```

## Base64 API

Every ArrayBuffer API has a base64 counterpart with the `AsBase64` suffix
(`derivationOfInitialKeyAsBase64`, `deriveCurrentTransactionKeyAsBase64`, `encryptPinAsBase64`,
`decryptPinAsBase64`, `encryptDataAsBase64`, `decryptDataAsBase64`, `generateMacAsBase64`).
Inputs and outputs are standard base64 strings with padding (RFC 4648), so values can be passed around as plain strings.
`decryptPinAsBase64` is the one exception and returns the PIN itself
(see [Breaking changes in v0.8.0](#breaking-changes-in-v080)):

```typescript
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64, generateMacAsBase64 } from "k6/x/dukpt/des";

const ik = "...";
const ck = deriveCurrentTransactionKeyAsBase64(ik, "//+YdlQyEOAAAQ==");
// ck === "BCZmtJGEz6No3pYo0Dl7yQ=="

const mac = generateMacAsBase64(ck, "NDAxMjM0NTY3ODkwOUQ5ODc=", "request"); // "4012345678909D987"
// mac === "nMx4Fz/E+2Q="
```

## Hex API

Every Base64 API has a hex counterpart with the `AsHex` suffix
(`derivationOfInitialKeyAsHex`, `deriveCurrentTransactionKeyAsHex`, `encryptPinAsHex`,
`decryptPinAsHex`, `encryptDataAsHex`, `decryptDataAsHex`, `generateMacAsHex`).
Inputs are case-insensitive hex strings; outputs are **uppercase** hex, so test vectors can be used as-is.
`decryptPinAsHex` is the one exception and returns the PIN itself
(see [Breaking changes in v0.8.0](#breaking-changes-in-v080)):

```typescript
import { derivationOfInitialKeyAsHex, deriveCurrentTransactionKeyAsHex, generateMacAsHex } from "k6/x/dukpt/des";

const ik = "...";
const ck = deriveCurrentTransactionKeyAsHex(ik, "FFFF9876543210E00001");
// ck === "042666B49184CFA368DE9628D0397BC9"

const mac = generateMacAsHex(ck, "3430313233343536373839303944393837", "request"); // "4012345678909D987"
// mac === "9CCC78173FC4FB64"
// the X9.24-1 mac is the first 8 hex characters: mac.slice(0, 8) === "9CCC7817"
```

## Download

Building a custom k6 binary with the `xk6-dukpt` extension is necessary for its use. You can download pre-built k6 binaries from the [Releases page](https://github.com/arukiidou/xk6-dukpt/releases/).

## Build

Use the [xk6](https://github.com/grafana/xk6) tool to build a custom k6 binary with the `xk6-dukpt` extension. Refer to the [xk6 documentation](https://github.com/grafana/xk6) for more information.

## License

`xk6-dukpt` is licensed under the [Apache License 2.0](LICENSE).

## Contribute

If you wish to contribute to this project, please start by reading the [Contributing Guidelines](CONTRIBUTING.md).

### Breaking changes in v0.8.0

The decryptData APIs now return the decrypted bytes encoded, instead of the raw padded string:

| API | v0.7.1 | v0.8.0 |
| --- | --- | --- |
| `decryptData` | string | **`ArrayBuffer`** |
| `decryptDataAsBase64` | string | **base64 string** |
| `decryptDataAsHex` | string | **uppercase hex string** |

`decryptPin`, `decryptPinAsBase64` and `decryptPinAsHex` are unchanged and still return the PIN itself,
because a PIN is already a digit string.

Scripts that sliced the old plain text **do not throw** after the upgrade, they just stop matching:
`ArrayBuffer.prototype.slice` returns another ArrayBuffer, so `decryptData(...).slice(0, n) === data`
silently becomes `false`. Decode the result instead:

```typescript
const out = decryptDataAsBase64(ck, ciphertext, "", "request");
// out === "NDAxMjM0NTY3ODkwOUQ5ODcAAAAAAAAA"
new TextDecoder().decode(Uint8Array.fromBase64(out).slice(0, 17)); // "4012345678909D987"
```

### Breaking changes in v0.10.0: stricter input validation

Inputs that moov-io silently accepts now return an error before moov-io is called.
This applies to all variants (ArrayBuffer, Base64, Hex), and is stricter than upstream.
Scripts that relied on the lenient behavior now fail.

| Argument | Accepted | moov-io behavior | APIs |
| --- | --- | --- | --- |
| `action` | exactly `"request"` or `"response"` | any other value (e.g. `"Response"`) is silently processed as `"request"` | `encryptData`, `decryptData`, `generateMac` |
| `bdk`, `ik`, `currentKey` | exactly 16 bytes | zero padded if short, truncated if long | `deriveCurrentTransactionKey`, `derivationOfInitialKey`, `encryptPin`, `decryptPin`, `encryptData`, `decryptData`, `generateMac` |
| `ksn` | exactly 10 bytes | not validated | `deriveCurrentTransactionKey`, `derivationOfInitialKey` |
| `iv` | empty (zero vector) or exactly 8 bytes | padded, truncated, or panics | `encryptData`, `decryptData` |
| `ciphertext` | non-empty multiple of 8 bytes | zero padded, then decrypted into garbage | `decryptData` |

Errors name the argument, e.g. `ksn must be 10 bytes, got 9` or
`ciphertext length must be a non-zero multiple of 8 bytes, got 15`.

### Breaking changes in v0.11.0: KSN length validation for the transaction counter

`getDesTcFromKsn` and `generateNextDesKsn` now return an error for a `ksn` shorter than 4 bytes.
This applies to all variants (ArrayBuffer, Base64, Hex), and is stricter than upstream.
Scripts that relied on the lenient behavior now fail.

| Argument | Accepted | moov-io behavior | APIs |
| --- | --- | --- | --- |
| `ksn` | at least 4 bytes | `getDesTcFromKsn` reads counter 0 below 4 bytes; `generateNextDesKsn` restarts the counter at 1 below 4 bytes and panics below 3 | `getDesTcFromKsn`, `generateNextDesKsn` |

The raw `getDesTcFromKsn` now throws instead of returning `0`; the Base64 and Hex variants already returned an error for undecodable input.

### Breaking changes in v0.12.0: plainText takes bytes

`encryptData` and `generateMac` now take `plainText` as bytes, encoded the same way as every other argument.
The signatures change in place, string input is no longer accepted.
This makes `encryptData` the inverse of `decryptData`, so arbitrary binary data can be encrypted and MACed.

| API | v0.11.0 `plainText` | v0.12.0 `plainText` |
| --- | --- | --- |
| `encryptData`, `generateMac` | string | **`ArrayBuffer` or `Uint8Array`** |
| `encryptDataAsBase64`, `generateMacAsBase64` | raw string | **base64 string** |
| `encryptDataAsHex`, `generateMacAsHex` | raw string | **hex string** (case-insensitive) |

Passing a string to `encryptData` or `generateMac` throws `plainText must be an ArrayBuffer or Uint8Array`.
A raw string passed to the Base64 and Hex variants is usually rejected as invalid encoding,
but one that happens to be valid base64 or hex (e.g. `"1234"`) is **silently decoded** into different bytes.
Encode the text first:

```typescript
const text = "4012345678909D987";
const bytes = new TextEncoder().encode(text);

encryptData(ck, null, bytes, "request");
encryptDataAsBase64(ck, "", bytes.toBase64(), "request");
encryptDataAsHex(ck, "", bytes.toHex(), "request");
```

An empty `plainText` is now an error in all variants (`plainText must not be empty`), stricter than upstream:
moov-io encrypts it into an empty ciphertext that `decryptData` rejects, and MACs it as a zero block.

## Breaking changes in v0.13.0: 

### transaction counter validation

`deriveCurrentTransactionKey` now returns an error for a `ksn` whose 21 bit transaction counter no ANSI X9.24-1 originator produces.
This applies to all variants (ArrayBuffer, Base64, Hex), and is stricter than upstream.
Scripts that relied on the lenient behavior now fail.

| Argument | Accepted | moov-io behavior | APIs |
| --- | --- | --- | --- |
| `ksn` transaction counter | non-zero, at most 10 one-bits | counter 0 returns the **initial key itself**; a counter with more than 10 one-bits derives a key | `deriveCurrentTransactionKey` |

Errors read `ksn transaction counter must not be 0` or
`ksn transaction counter must have at most 10 one-bits, got 11`.

`derivationOfInitialKey` is unchanged: it clears the counter, so a KSN with counter 0 still derives the initial key.
`getDesTcFromKsn` and `generateNextDesKsn` are unchanged as well; `generateNextDesKsn` from counter 0 returns counter 1.

### byte arguments take only ArrayBuffer or Uint8Array

Every byte argument of the raw APIs now rejects anything but an `ArrayBuffer` or `Uint8Array`,
as `plainText` already did since v0.12.0 and as the type definitions always declared.
Before, sobek converted a string into its UTF-8 bytes and a number array into one byte per element without complaint,
so e.g. `deriveCurrentTransactionKey("0123456789ABCDEF", "0123456789")` derived a key from the ASCII codes.
The Base64 and Hex variants are unchanged.

| Argument | Accepted | APIs |
| --- | --- | --- |
| `bdk`, `ik`, `ksn` | `ArrayBuffer` or `Uint8Array` | `derivationOfInitialKey`, `deriveCurrentTransactionKey` |
| `currentKey` | `ArrayBuffer` or `Uint8Array` | `encryptPin`, `decryptPin`, `encryptData`, `decryptData`, `generateMac` |
| `ciphertextPin`, `ciphertext` | `ArrayBuffer` or `Uint8Array` | `decryptPin`, `decryptData` |
| `iv` | `ArrayBuffer`, `Uint8Array`, `null` or `undefined` (zero vector) | `encryptData`, `decryptData` |
| `ksn` | `ArrayBuffer` or `Uint8Array` | `getDesTcFromKsn`, `generateNextDesKsn` (`k6/x/dukpt`) |

Rejected values include strings (also `""` as `iv`), number arrays, other typed arrays such as `Uint16Array`, and `DataView`.
Errors name the argument, e.g. `ksn must be an ArrayBuffer or Uint8Array` or
`iv must be an ArrayBuffer, Uint8Array, null or undefined`.
Convert the value first:

```typescript
deriveCurrentTransactionKey(Uint8Array.fromHex(ik), Uint8Array.fromHex(ksn));
deriveCurrentTransactionKey(new Uint8Array([/* ... */]), ksn); // instead of a number array
encryptData(ck, null, data, "request"); // instead of "" as iv
```
