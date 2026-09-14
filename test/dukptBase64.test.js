import { check } from "k6";
import { getDesTcFromKsnAsBase64, generateNextDesKsnAsBase64 } from "k6/x/dukpt";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

export default function () {
  const ksn = "//+YdlQyEOAAAQ=="; //"FFFF9876543210E00001";

  check(null, {
    'getDesTcFromKsnAsBase64(ksn)': () => getDesTcFromKsnAsBase64(ksn) === 1,
    'generateNextDesKsnAsBase64(ksn)': () => generateNextDesKsnAsBase64(ksn) === "//+YdlQyEOAAAg==", //"FFFF9876543210E00002"
  });
}
