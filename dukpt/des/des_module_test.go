// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.k6.io/k6/v2/js/modulestest"
)

func Test_module(t *testing.T) { //nolint:tparallel
	t.Parallel()

	runtime := modulestest.NewRuntime(t)

	err := runtime.SetupModuleSystem(map[string]any{ImportPath: new(DesModule)}, nil, nil)
	require.NoError(t, err)

	_, err = runtime.RunOnEventLoop(`let mod = require("` + ImportPath + `")`)
	require.NoError(t, err)

	// The raw APIs take ArrayBuffer or Uint8Array and hand back an ArrayBuffer,
	// so they are spelled through Uint8Array here.
	const (
		jsBdk = `Uint8Array.fromBase64("ASNFZ4mrze/+3LqYdlQyEA==")`
		jsKsn = `Uint8Array.fromBase64("//+YdlQyEOAAAQ==")`
		jsIk  = `Uint8Array.fromBase64("asKS+qExW02FirOj19WTOg==")`
		jsCk  = `Uint8Array.fromBase64("BCZmtJGEz6No3pYo0Dl7yQ==")`
		// "4012345678909D987"
		jsData = `Uint8Array.fromHex("3430313233343536373839303944393837")`
	)

	tests := []struct {
		name  string
		check string
	}{
		{name: "derivationOfInitialKey(bdk, ksn)", check: `new Uint8Array(mod.derivationOfInitialKey(` + jsBdk + `, ` + jsKsn + `)).toBase64() === "asKS+qExW02FirOj19WTOg=="`},
		{name: "deriveCurrentTransactionKey(ik, ksn)", check: `new Uint8Array(mod.deriveCurrentTransactionKey(` + jsIk + `, ` + jsKsn + `)).toBase64() === "BCZmtJGEz6No3pYo0Dl7yQ=="`},
		{name: "derivationOfInitialKey(bdk, ksn) returns an ArrayBuffer", check: `mod.derivationOfInitialKey(` + jsBdk + `, ` + jsKsn + `) instanceof ArrayBuffer`},
		{name: "encryptPin(currentKey, pin, pan, format)", check: `new Uint8Array(mod.encryptPin(` + jsCk + `, "1234", "4012345678909", "ISO-0")).toBase64() === "G5wYReuZOno="`},
		{name: "decryptPin(currentKey, ciphertext, pan, format)", check: `mod.decryptPin(` + jsCk + `, Uint8Array.fromBase64("G5wYReuZOno="), "4012345678909", "ISO-0") === "1234"`},
		{name: "encryptData(currentKey, null, plainText, action)", check: `new Uint8Array(mod.encryptData(` + jsCk + `, null, ` + jsData + `, "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, iv, plainText, action)", check: `new Uint8Array(mod.encryptData(` + jsCk + `, Uint8Array.fromHex("0000000000000000"), ` + jsData + `, "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, null, plainText, action) accepts an ArrayBuffer", check: `new Uint8Array(mod.encryptData(` + jsCk + `, null, ` + jsData + `.buffer, "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, null, string, action) throws", check: `(() => { try { mod.encryptData(` + jsCk + `, null, "4012345678909D987", "request"); return false } catch (e) { return String(e).includes("plainText must be an ArrayBuffer or Uint8Array") } })()`},
		{name: "generateMac(currentKey, string, action) throws", check: `(() => { try { mod.generateMac(` + jsCk + `, "4012345678909D987", "request"); return false } catch (e) { return String(e).includes("plainText must be an ArrayBuffer or Uint8Array") } })()`},
		{name: "encryptData(currentKey, null, subarray, action) honors the view offset", check: `new Uint8Array(mod.encryptData(` + jsCk + `, null, Uint8Array.fromHex("FF3430313233343536373839303944393837FF").subarray(1, 18), "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, null, Uint16Array, action) throws", check: `(() => { try { mod.encryptData(` + jsCk + `, null, new Uint16Array(9), "request"); return false } catch (e) { return String(e).includes("plainText must be an ArrayBuffer or Uint8Array") } })()`},
		{name: "generateMac(currentKey, null, action) throws", check: `(() => { try { mod.generateMac(` + jsCk + `, null, "request"); return false } catch (e) { return String(e).includes("plainText must be an ArrayBuffer or Uint8Array") } })()`},
		{name: "encryptData(currentKey, null, decryptData(...), action) round trips", check: `(() => { const ct = Uint8Array.fromBase64("/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"); const pt = mod.decryptData(` + jsCk + `, ct, null, "request"); return new Uint8Array(mod.encryptData(` + jsCk + `, null, pt, "request")).toBase64() === ct.toBase64() })()`},
		{name: "encryptDataAsHex(currentKey, iv, decryptDataAsHex(...), action) round trips", check: `mod.encryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "", mod.decryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F", "", "request"), "request") === "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"`},
		{name: "encryptDataAsBase64(currentKey, iv, decryptDataAsBase64(...), action) round trips", check: `mod.encryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "", mod.decryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP", "", "request"), "request") === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, shortIv, plainText, action) throws", check: `(() => { try { mod.encryptData(` + jsCk + `, Uint8Array.fromHex("0102"), ` + jsData + `, "request"); return false } catch (e) { return String(e).includes("iv must be 8 bytes, got 2") } })()`},
		{name: "decryptData(currentKey, ciphertext, null, action)", check: `new Uint8Array(mod.decryptData(` + jsCk + `, Uint8Array.fromBase64("/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"), null, "request")).toBase64() === "NDAxMjM0NTY3ODkwOUQ5ODcAAAAAAAAA"`},
		{name: "decryptData(currentKey, ciphertext, null, action) returns an ArrayBuffer", check: `mod.decryptData(` + jsCk + `, Uint8Array.fromBase64("/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"), null, "request") instanceof ArrayBuffer`},
		{name: `generateMac(currentKey, plainText, "request")`, check: `new Uint8Array(mod.generateMac(` + jsCk + `, ` + jsData + `, "request")).toBase64() === "nMx4Fz/E+2Q="`},
		{name: `generateMac(currentKey, plainText, "response")`, check: `new Uint8Array(mod.generateMac(` + jsCk + `, ` + jsData + `, "response")).toBase64() === "IDZCI8H/APo="`},
		{name: "derivationOfInitialKey(bdk, ksn) keeps the input", check: `(() => { const bdk = ` + jsBdk + `; mod.derivationOfInitialKey(bdk, ` + jsKsn + `); return bdk.toBase64() === "ASNFZ4mrze/+3LqYdlQyEA==" })()`},
		{name: "derivationOfInitialKeyAsBase64(bdk, ksn)", check: `mod.derivationOfInitialKeyAsBase64("ASNFZ4mrze/+3LqYdlQyEA==", "//+YdlQyEOAAAQ==") === "asKS+qExW02FirOj19WTOg=="`},
		{name: "deriveCurrentTransactionKeyAsBase64(ik, ksn)", check: `mod.deriveCurrentTransactionKeyAsBase64("asKS+qExW02FirOj19WTOg==", "//+YdlQyEOAAAQ==") === "BCZmtJGEz6No3pYo0Dl7yQ=="`},
		{name: "encryptPinAsBase64(currentKey, pin, pan, format)", check: `mod.encryptPinAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "1234", "4012345678909", "ISO-0") === "G5wYReuZOno="`},
		{name: "decryptPinAsBase64(currentKey, ciphertext, pan, format)", check: `mod.decryptPinAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "G5wYReuZOno=", "4012345678909", "ISO-0") === "1234"`},
		{name: "encryptDataAsBase64(currentKey, iv, plainText, action)", check: `mod.encryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "", "NDAxMjM0NTY3ODkwOUQ5ODc=", "request") === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "decryptDataAsBase64(currentKey, ciphertext, iv, action)", check: `mod.decryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP", "", "request") === "NDAxMjM0NTY3ODkwOUQ5ODcAAAAAAAAA"`},
		{name: `generateMacAsBase64(currentKey, plainText, "request")`, check: `mod.generateMacAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "NDAxMjM0NTY3ODkwOUQ5ODc=", "request") === "nMx4Fz/E+2Q="`},
		{name: `generateMacAsBase64(currentKey, plainText, "response")`, check: `mod.generateMacAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "NDAxMjM0NTY3ODkwOUQ5ODc=", "response") === "IDZCI8H/APo="`},
		{name: "derivationOfInitialKeyAsHex(bdk, ksn)", check: `mod.derivationOfInitialKeyAsHex("0123456789ABCDEFFEDCBA9876543210", "FFFF9876543210E00001") === "6AC292FAA1315B4D858AB3A3D7D5933A"`},
		{name: "deriveCurrentTransactionKeyAsHex(ik, ksn)", check: `mod.deriveCurrentTransactionKeyAsHex("6AC292FAA1315B4D858AB3A3D7D5933A", "FFFF9876543210E00001") === "042666B49184CFA368DE9628D0397BC9"`},
		{name: "encryptPinAsHex(currentKey, pin, pan, format)", check: `mod.encryptPinAsHex("042666B49184CFA368DE9628D0397BC9", "1234", "4012345678909", "ISO-0") === "1B9C1845EB993A7A"`},
		{name: "decryptPinAsHex(currentKey, ciphertext, pan, format)", check: `mod.decryptPinAsHex("042666B49184CFA368DE9628D0397BC9", "1B9C1845EB993A7A", "4012345678909", "ISO-0") === "1234"`},
		{name: "encryptDataAsHex(currentKey, iv, plainText, action)", check: `mod.encryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "", "3430313233343536373839303944393837", "request") === "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"`},
		{name: "decryptDataAsHex(currentKey, ciphertext, iv, action)", check: `mod.decryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F", "", "request") === "343031323334353637383930394439383700000000000000"`},
		{name: `generateMacAsHex(currentKey, plainText, "request")`, check: `mod.generateMacAsHex("042666B49184CFA368DE9628D0397BC9", "3430313233343536373839303944393837", "request") === "9CCC78173FC4FB64"`},
		{name: `generateMacAsHex(currentKey, plainText, "response")`, check: `mod.generateMacAsHex("042666B49184CFA368DE9628D0397BC9", "3430313233343536373839303944393837", "response") === "20364223C1FF00FA"`},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.RunOnEventLoop(tt.check)

			require.NoError(t, err)
			require.True(t, got.ToBoolean())
		})
	}
}

// sobek converts strings, number arrays and other typed arrays into []byte
// without complaint, so every byte argument of the raw APIs has to reject them.
func Test_module_nonByteArguments(t *testing.T) { //nolint:tparallel
	t.Parallel()

	runtime := modulestest.NewRuntime(t)

	err := runtime.SetupModuleSystem(map[string]any{ImportPath: new(DesModule)}, nil, nil)
	require.NoError(t, err)

	_, err = runtime.RunOnEventLoop(`
		var mod = require("` + ImportPath + `");
		var key = Uint8Array.fromBase64("BCZmtJGEz6No3pYo0Dl7yQ==");
		var ksn = Uint8Array.fromBase64("//+YdlQyEOAAAQ==");
		var pinBlock = Uint8Array.fromBase64("G5wYReuZOno=");
		var data = Uint8Array.fromHex("3430313233343536373839303944393837");
		var ct = Uint8Array.fromBase64("/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP");
		// rejected returns the error a call throws, or "" when it does not throw.
		var rejected = (f) => { try { f(); return "" } catch (e) { return String(e) } };
	`)
	require.NoError(t, err)

	// Each call takes the value under test in place of $.
	args := []struct {
		name string
		call string
	}{
		{name: "bdk", call: `mod.derivationOfInitialKey($, ksn)`},
		{name: "ksn", call: `mod.derivationOfInitialKey(key, $)`},
		{name: "ik", call: `mod.deriveCurrentTransactionKey($, ksn)`},
		{name: "ksn", call: `mod.deriveCurrentTransactionKey(key, $)`},
		{name: "currentKey", call: `mod.encryptPin($, "1234", "4012345678909", "ISO-0")`},
		{name: "currentKey", call: `mod.decryptPin($, pinBlock, "4012345678909", "ISO-0")`},
		{name: "ciphertextPin", call: `mod.decryptPin(key, $, "4012345678909", "ISO-0")`},
		{name: "currentKey", call: `mod.encryptData($, null, data, "request")`},
		{name: "plainText", call: `mod.encryptData(key, null, $, "request")`},
		{name: "currentKey", call: `mod.decryptData($, ct, null, "request")`},
		{name: "ciphertext", call: `mod.decryptData(key, $, null, "request")`},
		{name: "currentKey", call: `mod.generateMac($, data, "request")`},
		{name: "plainText", call: `mod.generateMac(key, $, "request")`},
	}
	values := []string{
		`"0123456789ABCDEF"`,
		`""`,
		`[0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]`,
		`new Uint16Array(16)`,
		`new DataView(new ArrayBuffer(16))`,
		`42`,
		`{}`,
		`null`,
		`undefined`,
	}
	for _, arg := range args {
		for _, value := range values {
			call := strings.ReplaceAll(arg.call, "$", value)
			t.Run(call, func(t *testing.T) {
				got, err := runtime.RunOnEventLoop(`rejected(() => ` + call + `)`)
				require.NoError(t, err)
				require.Contains(t, got.String(), arg.name+" must be an ArrayBuffer or Uint8Array")
			})
		}
	}

	ivCalls := []string{
		`mod.encryptData(key, $, data, "request")`,
		`mod.decryptData(key, ct, $, "request")`,
	}
	for _, ivCall := range ivCalls {
		// null and undefined stand for the zero vector.
		for _, value := range []string{`null`, `undefined`} {
			call := strings.ReplaceAll(ivCall, "$", value)
			t.Run(call, func(t *testing.T) {
				got, err := runtime.RunOnEventLoop(`new Uint8Array(` + call + `).toBase64() === new Uint8Array(` +
					strings.ReplaceAll(ivCall, "$", `new Uint8Array(8)`) + `).toBase64()`)
				require.NoError(t, err)
				require.True(t, got.ToBoolean())
			})
		}
		for _, value := range values[:7] {
			call := strings.ReplaceAll(ivCall, "$", value)
			t.Run(call, func(t *testing.T) {
				got, err := runtime.RunOnEventLoop(`rejected(() => ` + call + `)`)
				require.NoError(t, err)
				require.Contains(t, got.String(), "iv must be an ArrayBuffer, Uint8Array, null or undefined")
			})
		}
	}
}
