import { check } from "k6";
import { derivationOfInitialKeyAsBase64, deriveCurrentTransactionKeyAsBase64, encryptPinAsBase64, decryptPinAsBase64, encryptDataAsBase64, decryptDataAsBase64, generateMacAsBase64 } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ikExpected = "asKS+qExW02FirOj19WTOg==";//"0123456789ABCDEFFEDCBA9876543210";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";
  const pinEncExpected = "G5wYReuZOno="; //"1B9C1845EB993A7A";

  const bdk = "ASNFZ4mrze/+3LqYdlQyEA=="; //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = "//+YdlQyEOAAAQ=="; //"FFFF9876543210E00001";

  const ik = derivationOfInitialKeyAsBase64(bdk, ksn)
  const ck = deriveCurrentTransactionKeyAsBase64(ik, ksn)

  check(null, {
    'derivationOfInitialKeyAsBase64(bdk, ksn)': () => ik === ikExpected,
    'deriveCurrentTransactionKeyAsBase64(ik, ksn)': () => ck === ckExpected,
  });

  const pin = "1234";
  const pan = "4012345678909";
  const format = "ISO-0";
  
  const pinEnc = encryptPinAsBase64(ck, pin, pan, format)

  check(null, {
    'encryptPinAsBase64(ck, pin, pan, format)': () => pinEnc === pinEncExpected,
    'decryptPinAsBase64(ck, pinEnc, pan, format)': () => decryptPinAsBase64(ck, pinEnc, pan, format) === pin,
  });

  const dataReqEncExpected = "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"; //"FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F";
  const dataResEncExpected = "H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl"; //"1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25";

  const data = "4012345678909D987";

  const reqEnc = encryptDataAsBase64(ck, "", data, "request")
  const resEnc = encryptDataAsBase64(ck, "", data, "response")

  check(null, {
    'encryptDataAsBase64(ck, "", data, "request")': () => reqEnc === dataReqEncExpected,
    'encryptDataAsBase64(ck, "", data, "response")': () => resEnc === dataResEncExpected,
    'decryptDataAsBase64(ck, reqEnc, "", "request")': () => decryptDataAsBase64(ck, reqEnc, "", "request").slice(0, data.length) === data,
    'decryptDataAsBase64(ck, resEnc, "", "response")': () => decryptDataAsBase64(ck, resEnc, "", "response").slice(0, data.length) === data,
  });

  const reqMacExpected = "nMx4Fz/E+2Q="; //"9CCC78173FC4FB64";
  const resMacExpected = "IDZCI8H/APo="; //"20364223C1FF00FA";

  check(null, {
    'generateMacAsBase64(ck, data, "request")': () => generateMacAsBase64(ck, data, "request") === reqMacExpected,
    'generateMacAsBase64(ck, data, "response")': () => generateMacAsBase64(ck, data, "response") === resMacExpected,
  });
}
