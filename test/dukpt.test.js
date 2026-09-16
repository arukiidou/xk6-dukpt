import { check } from "k6";
import { getDesTcFromKsn, generateNextDesKsn } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ=="); //"FFFF9876543210E00001";

  check(null, {
    'getDesTcFromKsn(ksn)': () => getDesTcFromKsn(ksn) === 1,
  });

  const nextKsnExpected = "//+YdlQyEOAAAg=="; //"FFFF9876543210E00002";

  const nextKsn = generateNextDesKsn(ksn);

  check(null, {
    'generateNextDesKsn(ksn)': () => new Uint8Array(nextKsn).toBase64() === nextKsnExpected,
    'generateNextDesKsn(ksn) keeps ksn': () => ksn.toBase64() === "//+YdlQyEOAAAQ==",
    'generateNextDesKsn(ksn) returns an ArrayBuffer': () => nextKsn instanceof ArrayBuffer,
  });
}
