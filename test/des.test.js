import { check } from "k6";
import { derivationOfInitialKey, deriveCurrentTransactionKey, encryptPin, decryptPin, encryptData, decryptData, generateMac } from "k6/x/dukpt/des";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default async function () {
  const ikExpected = "asKS+qExW02FirOj19WTOg=="; //"6AC292FAA1315B4D858AB3A3D7D5933A";
  const ckExpected = "BCZmtJGEz6No3pYo0Dl7yQ=="; //"042666B49184CFA368DE9628D0397BF9";
  const pinEncExpected = "G5wYReuZOno="; //"1B9C1845EB993A7A";

  const bdk = Uint8Array.fromBase64("ASNFZ4mrze/+3LqYdlQyEA=="); //"0123456789ABCDEFFEDCBA9876543210";
  const ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  const ik = derivationOfInitialKey(bdk, ksn);
  const ck = deriveCurrentTransactionKey(ik, ksn);

  check(null, {
    'derivationOfInitialKey(bdk, ksn)': () => new Uint8Array(ik).toBase64() === ikExpected,
    'deriveCurrentTransactionKey(ik, ksn)': () => new Uint8Array(ck).toBase64() === ckExpected,
  });

  const pin = "1234";
  const pan = "4012345678909";
  const format = "ISO-0";

  const pinEnc = encryptPin(ck, pin, pan, format);

  check(null, {
    'encryptPin(ck, pin, pan, format)': () => new Uint8Array(pinEnc).toBase64() === pinEncExpected,
    'decryptPin(ck, pinEnc, pan, format)': () => decryptPin(ck, pinEnc, pan, format) === pin,
  });

  const dataReqEncExpected = "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"; //"FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F";
  const dataResEncExpected = "H8yJr2YiLye5A4mLsryFic2/3l7Gr8wl"; //"1FCC89AF66222F27B903898BB2BC8589CDBFDE5EC6AFCC25";

  const data = "4012345678909D987";

  const reqEnc = encryptData(ck, null, data, "request");
  const resEnc = encryptData(ck, null, data, "response");

  check(null, {
    'encryptData(ck, null, data, "request")': () => new Uint8Array(reqEnc).toBase64() === dataReqEncExpected,
    'encryptData(ck, null, data, "response")': () => new Uint8Array(resEnc).toBase64() === dataResEncExpected,
    'decryptData(ck, reqEnc, null, "request")': () => decryptData(ck, reqEnc, null, "request").slice(0, data.length) === data,
    'decryptData(ck, resEnc, null, "response")': () => decryptData(ck, resEnc, null, "response").slice(0, data.length) === data,
  });

  const reqMacExpected = "nMx4Fz/E+2Q="; //"9CCC78173FC4FB64";
  const resMacExpected = "IDZCI8H/APo="; //"20364223C1FF00FA";

  const reqMac = generateMac(ck, data, "request");
  const resMac = generateMac(ck, data, "response");

  check(null, {
    'generateMac(ck, data, "request")': () => new Uint8Array(reqMac).toBase64() === reqMacExpected,
    'generateMac(ck, data, "response")': () => new Uint8Array(resMac).toBase64() === resMacExpected,
  });
}
