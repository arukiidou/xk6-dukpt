import { check } from "k6";
import encoding from 'k6/encoding';
import { derivationOfInitialKey, deriveCurrentTransactionKey, encryptPin, decryptPin, encryptData, decryptData, generateMac } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default async function () {
  const ikExpected = "asKS+qExW02FirOj19WTOg==";//"0123456789ABCDEFFEDCBA9876543210";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";
  const pinEncExpected = "G5wYReuZOno="; //"1B9C1845EB993A7A";

  // TODO: wait for merge Uint8Array.fromBase64() to sobek
  const bdk = encoding.b64decode("ASNFZ4mrze/+3LqYdlQyEA=="); //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = encoding.b64decode("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  const ik = derivationOfInitialKey(bdk, ksn);
  const ck = deriveCurrentTransactionKey(ik, ksn);

  check(null, {
    'derivationOfInitialKey(bdk, ksn)': () => encoding.b64encode(ik) === ikExpected,
    'deriveCurrentTransactionKey(ik, ksn)': () => encoding.b64encode(ck) === ckExpected,
  });

  const pin = "1234";
  const pan = "4012345678909";
  const format = "ISO-0";

  const pinEnc = encryptPin(ck, pin, pan, format);

  check(null, {
    'encryptPin(ck, pin, pan, format)': () => encoding.b64encode(pinEnc) === pinEncExpected,
    'decryptPin(ck, pinEnc, pan, format)': () => decryptPin(ck, pinEnc, pan, format) === pin,
  });

  const dataReqEncExpected = "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"; //"FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F";
  const dataResEncExpected = "H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl"; //"1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25";

  const data = "4012345678909D987";

  const reqEnc = encryptData(ck, null, data, "request");
  const resEnc = encryptData(ck, null, data, "response");

  check(null, {
    'encryptData(ck, null, data, "request")': () => encoding.b64encode(reqEnc) === dataReqEncExpected,
    'encryptData(ck, null, data, "response")': () => encoding.b64encode(resEnc) === dataResEncExpected,
    'decryptData(ck, reqEnc, null, "request")': () => decryptData(ck, reqEnc, null, "request").slice(0, data.length) === data,
    'decryptData(ck, resEnc, null, "response")': () => decryptData(ck, resEnc, null, "response").slice(0, data.length) === data,
  });

  const reqMacExpected = "nMx4Fz/E+2Q="; //"9CCC78173FC4FB64";
  const resMacExpected = "IDZCI8H/APo="; //"20364223C1FF00FA";

  const reqMac = generateMac(ck, data, "request");
  const resMac = generateMac(ck, data, "response");

  check(null, {
    'generateMac(ck, data, "request")': () => encoding.b64encode(reqMac) === reqMacExpected,
    'generateMac(ck, data, "response")': () => encoding.b64encode(resMac) === resMacExpected,
  });
}
