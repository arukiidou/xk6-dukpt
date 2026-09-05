# xk6-dukpt

**k6 DUKPT extension porting with moov-io**

```typescript file=dukpt.ts
import { check } from "k6";
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64 } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ikExpected = "asKS+qExW02FirOj19WTOg==";//"0123456789ABCDEFFEDCBA9876543210";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";

  const bdk = "ASNFZ4mrze/+3LqYdlQyEA=="; //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = "//+YdlQyEOAAAQ=="; //"FFFF9876543210E00001";

  const ik = derivationOfInitialKeyAsBase64(bdk, ksn)
  console.log("Actual: " + ik);
  console.log("Expected: " + ikExpected);

  const ck = deriveCurrentTransactionKeyAsBase64(ik, ksn)
  console.log("Actual: " + ck);
  console.log("Expected: " + ckExpected);

  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => ik === ikExpected,
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => ck === ckExpected,
  });
}

```

## Download

Building a custom k6 binary with the `xk6-example` extension is necessary for its use. You can download pre-built k6 binaries from the [Releases page](https://github.com/grafana/xk6-example/releases/).

## Build

Use the [xk6](https://github.com/grafana/xk6) tool to build a custom k6 binary with the `xk6-example` extension. Refer to the [xk6 documentation](https://github.com/grafana/xk6) for more information.

## Contribute

If you wish to contribute to this project, please start by reading the [Contributing Guidelines](CONTRIBUTING.md).
