import { check } from "k6";
import { b32encode, b32decode } from "k6/x/example";

export const options = {
  thresholds: {
    checks: ["rate==1"],
  },
};

function toString(ab) {
  return Array.from(new Uint8Array(ab), (c) => String.fromCharCode(c)).join('')
}

export default function () {
  const input = "The quick brown fox jumps over the lazy dog.";
  const std = "KRUGKIDROVUWG2ZAMJZG653OEBTG66BANJ2W24DTEBXXMZLSEB2GQZJANRQXU6JAMRXWOLQ="
  const stdraw = "KRUGKIDROVUWG2ZAMJZG653OEBTG66BANJ2W24DTEBXXMZLSEB2GQZJANRQXU6JAMRXWOLQ"
  const hex = "AHK6A83HELKM6QP0C9P6UTRE41J6UU10D9QMQS3J41NNCPBI41Q6GP90DHGNKU90CHNMEBG="
  const hexraw = "AHK6A83HELKM6QP0C9P6UTRE41J6UU10D9QMQS3J41NNCPBI41Q6GP90DHGNKU90CHNMEBG"

  // b32encode

  check(b32encode(input, "std"), {
    'b32encode(input, "std")': (str) => str == std,
  });

  check(b32encode(input, "stdraw"), {
    'b32encode(input, "stdraw")': (str) => str == stdraw,
  });

  check(b32encode(input, "hex"), {
    'b32encode(input, "hex")': (str) => str == hex,
  });

  check(b32encode(input, "hexraw"), {
    'b32encode(input, "hexraw")': (str) => str == hexraw,
  });

  // b32decode

  check(toString(b32decode(std, "std")), {
    'b32decode(input, "std")': (str) => str == input,
  });

  check(b32decode(std, "std", "s"), {
    'b32decode(input, "std", "s")': (str) => str == input,
  });

  check(toString(b32decode(stdraw, "stdraw")), {
    'b32decode(input, "stdraw")': (str) => str == input,
  });

  check(b32decode(stdraw, "stdraw", "s"), {
    'b32decode(input, "stdraw", "s")': (str) => str == input,
  });

  check(toString(b32decode(hex, "hex")), {
    'b32decode(input, "hex")': (str) => str == input,
  });

  check(b32decode(hexraw, "hexraw", "s"), {
    'b32decode(input, "hexraw", "s")': (str) => str == input,
  });
}
