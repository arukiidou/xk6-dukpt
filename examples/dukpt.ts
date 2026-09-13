import { check } from "k6";
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64, generateMacAsHex } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  // this example uses well-known BDK, but you should replace it with your own.
  const ikExpected = "asKS+qExW02FirOj19WTOg=="; //"6AC292FAA1315B4D858AB3A3D7D5933A";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";

  const bdk = "ASNFZ4mrze/+3LqYdlQyEA=="; //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = "//+YdlQyEOAAAQ=="; //"FFFF9876543210E00001";

  const ik = derivationOfInitialKeyAsBase64(bdk, ksn)
  const ck = deriveCurrentTransactionKeyAsBase64(ik, ksn)

  // generateMacAsBase64 returns 8 bytes; the ANSI X9.24-1 mac is the first 4.
  const macExpected = "9CCC78173FC4FB64"; //"nMx4Fz/E+2Q=";
  const mac = generateMacAsHex(Uint8Array.fromBase64(ck).toHex(), "4012345678909D987", "request")

  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => ik === ikExpected,
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => ck === ckExpected,
    'generateMacAsHex(ck, data, "request")': () => mac === macExpected,
  });
}
