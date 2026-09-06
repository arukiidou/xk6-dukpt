# xk6-dukpt

**[k6](https://github.com/grafana/k6) DUKPT [extension](https://github.com/grafana/xk6) port from [moov-io](https://pkg.go.dev/github.com/moov-io/dukpt)**

[![Go Reference](https://pkg.go.dev/badge/github.com/arukiidou/xk6-dukpt.svg)](https://pkg.go.dev/github.com/arukiidou/xk6-dukpt)

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

## Download

Building a custom k6 binary with the `xk6-example` extension is necessary for its use. You can download pre-built k6 binaries from the [Releases page](https://github.com/grafana/xk6-example/releases/).

## Build

Use the [xk6](https://github.com/grafana/xk6) tool to build a custom k6 binary with the `xk6-example` extension. Refer to the [xk6 documentation](https://github.com/grafana/xk6) for more information.

## Contribute

If you wish to contribute to this project, please start by reading the [Contributing Guidelines](CONTRIBUTING.md).
