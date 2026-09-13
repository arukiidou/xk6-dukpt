import { check } from "k6";
import { derivationOfInitialKeyAsHex, deriveCurrentTransactionKeyAsHex, encryptPinAsHex, decryptPinAsHex, encryptDataAsHex, decryptDataAsHex, generateMacAsHex } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ikExpected = "6AC292FAA1315B4D858AB3A3D7D5933A";
  const ckExpected = "042666B49184CFA368DE9628D0397BC9";
  const pinEncExpected = "1B9C1845EB993A7A";

  const bdk = "0123456789ABCDEFFEDCBA9876543210";
  const ksn = "FFFF9876543210E00001";

  const ik = derivationOfInitialKeyAsHex(bdk, ksn)
  const ck = deriveCurrentTransactionKeyAsHex(ik, ksn)

  check(null, {
    'derivationOfInitialKeyAsHex(bdk, ksn)': () => ik === ikExpected,
    'deriveCurrentTransactionKeyAsHex(ik, ksn)': () => ck === ckExpected,
  });

  const pin = "1234";
  const pan = "4012345678909";
  const format = "ISO-0";

  const pinEnc = encryptPinAsHex(ck, pin, pan, format)

  check(null, {
    'encryptPinAsHex(ck, pin, pan, format)': () => pinEnc === pinEncExpected,
    'decryptPinAsHex(ck, pinEnc, pan, format)': () => decryptPinAsHex(ck, pinEnc, pan, format) === pin,
  });

  const dataReqEncExpected = "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F";
  const dataResEncExpected = "1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25";

  const data = "4012345678909D987";

  const reqEnc = encryptDataAsHex(ck, "", data, "request")
  const resEnc = encryptDataAsHex(ck, "", data, "response")

  check(null, {
    'encryptDataAsHex(ck, "", data, "request")': () => reqEnc === dataReqEncExpected,
    'encryptDataAsHex(ck, "", data, "response")': () => resEnc === dataResEncExpected,
    'decryptDataAsHex(ck, reqEnc, "", "request")': () => decryptDataAsHex(ck, reqEnc, "", "request").slice(0, data.length) === data,
    'decryptDataAsHex(ck, resEnc, "", "response")': () => decryptDataAsHex(ck, resEnc, "", "response").slice(0, data.length) === data,
  });

  const reqMacExpected = "9CCC78173FC4FB64";
  const resMacExpected = "20364223C1FF00FA";

  check(null, {
    'generateMacAsHex(ck, data, "request")': () => generateMacAsHex(ck, data, "request") === reqMacExpected,
    'generateMacAsHex(ck, data, "response")': () => generateMacAsHex(ck, data, "response") === resMacExpected,
  });
}
