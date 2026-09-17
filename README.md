# xk6-dukpt

[![Go Reference](https://pkg.go.dev/badge/github.com/arukiidou/xk6-dukpt.svg)](https://pkg.go.dev/github.com/arukiidou/xk6-dukpt)

- xk6-dukpt is a [k6](https://github.com/grafana/k6) [extension](https://github.com/grafana/xk6) for DUKPT(Derived Unique Key Per Transaction) cryptography operations. port from [moov-io](https://pkg.go.dev/github.com/moov-io/dukpt)

# Important

As with moov-io/dukpt, this extention is mostly for validation and debugging purposes.

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

## Requrements

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
import { derivationOfInitialKey, deriveCurrentTransactionKey } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  // this example uses well-known BDK, but you should replace it with your own.
  example("6AC292FAA1315B4D858AB3A3D7D5933A", "FFFF9876543210E00001");
}

function example(ik: string, ksn: string) {
  const ckExpected = "042666B49184CFA368DE9628D0397BC9"; //"BCZmtJGEz6No3pYo0Dl7yQ==";

  const ck = deriveCurrentTransactionKey(Uint8Array.fromHex(ik), Uint8Array.fromHex(ksn));
  check(null, {
    'deriveCurrentTransactionKey(ik, ksn)': () => new Uint8Array(ck).toHex().toUpperCase() === ckExpected,
  });
}
```

## MAC generation

`generateMac` / `generateMacAsBase64` port [des.GenerateMac](https://pkg.go.dev/github.com/moov-io/dukpt/pkg/des#GenerateMac)
and return the **full 8 bytes** produced by moov-io. Per ANSI X9.24-1, the MAC is
the **first 4 bytes** of that result, so truncate before sending it:

```typescript
import { generateMacAsBase64 } from "k6/x/dukpt/des";

const mac = generateMacAsBase64(ck, "4012345678909D987", "request");
// mac === "nMx4Fz/E+2Q="  (9CCC78173FC4FB64)
// the X9.24-1 mac is the first 4 bytes: 9CCC7817
```

## Base64 API

Every ArrayBuffer API has a base64 counterpart with the `AsBase64` suffix
(`derivationOfInitialKeyAsBase64`, `deriveCurrentTransactionKeyAsBase64`, `encryptPinAsBase64`,
`decryptPinAsBase64`, `encryptDataAsBase64`, `decryptDataAsBase64`, `generateMacAsBase64`).
Inputs and outputs are standard base64 strings with padding (RFC 4648), so values can be passed around as plain strings:

```typescript
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64, generateMacAsBase64 } from "k6/x/dukpt/des";

const ik = "...";
const ck = deriveCurrentTransactionKeyAsBase64(ik, "//+YdlQyEOAAAQ==");
// ck === "BCZmtJGEz6No3pYo0Dl7yQ=="

const mac = generateMacAsBase64(ck, "4012345678909D987", "request");
// mac === "nMx4Fz/E+2Q="
```

## Hex API

Every Base64 API has a hex counterpart with the `AsHex` suffix
(`derivationOfInitialKeyAsHex`, `deriveCurrentTransactionKeyAsHex`, `encryptPinAsHex`,
`decryptPinAsHex`, `encryptDataAsHex`, `decryptDataAsHex`, `generateMacAsHex`).
Inputs are case-insensitive hex strings; outputs are **uppercase** hex, so test vectors can be used as-is:

```typescript
import { derivationOfInitialKeyAsHex, deriveCurrentTransactionKeyAsHex, generateMacAsHex } from "k6/x/dukpt/des";

const ik = "...";
const ck = deriveCurrentTransactionKeyAsHex(ik, "FFFF9876543210E00001");
// ck === "042666B49184CFA368DE9628D0397BC9"

const mac = generateMacAsHex(ck, "4012345678909D987", "request");
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
