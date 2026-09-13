import desBase64Test from "./desBase64.test.js"
import desHexTest from "./desHex.test.js"
import desTest from "./des.test.js"

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  desBase64Test()
  desHexTest()
  desTest()
}
