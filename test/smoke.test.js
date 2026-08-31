import greeting from "./greeting.test.js"
import base32 from "./base32.test.js"
import random from "./random.test.js"

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  greeting()
  base32()
  random()
}
