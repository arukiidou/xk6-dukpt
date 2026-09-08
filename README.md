# xk6-dukpt

- [k6](https://github.com/grafana/k6) DUKPT [extension](https://github.com/grafana/xk6)
- port from [moov-io](https://pkg.go.dev/github.com/moov-io/dukpt)

[![Go Reference](https://pkg.go.dev/badge/github.com/arukiidou/xk6-dukpt.svg)](https://pkg.go.dev/github.com/arukiidou/xk6-dukpt)

## Requrements

- k6 2.2.0+

## How to Build

```bash
go install go.k6.io/xk6@latest
xk6 build --k6-version latest --os linux --cgo 0 --with github.com/arukiidou/xk6-dukpt@latest
# or
# go get -tool go.k6.io/xk6@latest
# go run go.k6.io/xk6 build --k6-version latest --os linux --cgo 0 --with github.com/arukiidou/xk6-dukpt@latest
```

# Example script

```typescript file=dukpt.ts
import { check } from "k6";
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64 } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default async function () {
  ...  
}

async function example(bdk: string, ksn: string) {

  const ik = derivationOfInitialKeyAsBase64(bdk, ksn)
  const ck = deriveCurrentTransactionKeyAsBase64(ik, ksn)

  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => ik === "...",
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => ck === "...",
  });
}

```

## MAC generation

`generateMac` / `generateMacAsBase64` port [des.GenerateMac](https://pkg.go.dev/github.com/moov-io/dukpt/pkg/des#GenerateMac)
and return the **full 8 bytes** produced by moov-io. Per ANSI X9.24-1, the MAC is
the **first 4 bytes** of that result, so truncate before sending it:

```typescript
import { generateMacAsBase64 } from "k6/x/dukpt";

const mac = generateMacAsBase64(ck, "4012345678909D987", "request");
// mac === "nMx4Fz/E+2Q="  (9CCC78173FC4FB64)
// the X9.24-1 mac is the first 4 bytes: 9CCC7817
```

## Download

Building a custom k6 binary with the `xk6-dukpt` extension is necessary for its use. You can download pre-built k6 binaries from the [Releases page](https://github.com/arukiidou/xk6-dukpt/releases/).

## Build

Use the [xk6](https://github.com/grafana/xk6) tool to build a custom k6 binary with the `xk6-dukpt` extension. Refer to the [xk6 documentation](https://github.com/grafana/xk6) for more information.

## Contribute

If you wish to contribute to this project, please start by reading the [Contributing Guidelines](CONTRIBUTING.md).
