import { check } from "k6";
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64, generateMacAsBase64 } from "k6/x/dukpt";

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

  // generateMacAsBase64 returns 8 bytes; the ANSI X9.24-1 mac is the first 4.
  const macExpected = "nMx4Fz/E+2Q="; //"9CCC78173FC4FB64";
  const mac = generateMacAsBase64(ck, "4012345678909D987", "request")
  console.log("Actual: " + mac);
  console.log("Expected: " + macExpected);

  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => ik === ikExpected,
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => ck === ckExpected,
    'generateMacAsBase64(ck, data, "request")': () => mac === macExpected,
  });
}
