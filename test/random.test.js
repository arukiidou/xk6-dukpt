import { check } from "k6";
import { Random } from "k6/x/example";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
}

export default function () {
  const seed = 42

  // without seed

  let rnd = new Random()

  check(rnd, {
    "new Random()": (rnd) => rnd.seed != 0
  })

  check(rnd, {
    "rnd.int() before seed change": (rnd) => rnd.int() != 8836438174961104
  })

  rnd.seed = seed

  check(rnd, {
    "rnd.int() after seed change": (rnd) => rnd.int() == 8836438174961104
  })

  // with seed

  check(new Random(seed), {
    "new Random(seed)": (rnd) => rnd.seed == 42
  })

  check(new Random(seed), {
    "int()": (rnd) => rnd.int() == 8836438174961104
  })

  check(new Random(seed), {
    "int(n)": (rnd) => rnd.int(2000) == 675
  })

  check(new Random(seed), {
    "float()": (rnd) => Math.round(rnd.float() * 100) == 37
  })

  check(new Random(seed), {
    "float(n)": (rnd) => Math.round(rnd.float(2000)) == 746
  })
}
