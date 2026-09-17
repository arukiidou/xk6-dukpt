// SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
// SPDX-License-Identifier: Apache-2.0

package des

import (
	_ "embed"
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
		{name: "encryptData(currentKey, null, plainText, action)", check: `new Uint8Array(mod.encryptData(` + jsCk + `, null, "4012345678909D987", "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, iv, plainText, action)", check: `new Uint8Array(mod.encryptData(` + jsCk + `, Uint8Array.fromHex("0000000000000000"), "4012345678909D987", "request")).toBase64() === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "encryptData(currentKey, shortIv, plainText, action) throws", check: `(() => { try { mod.encryptData(` + jsCk + `, Uint8Array.fromHex("0102"), "4012345678909D987", "request"); return false } catch (e) { return String(e).includes("iv must be 8 bytes, got 2") } })()`},
		{name: "decryptData(currentKey, ciphertext, null, action)", check: `mod.decryptData(` + jsCk + `, Uint8Array.fromBase64("/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"), null, "request").slice(0, 17) === "4012345678909D987"`},
		{name: `generateMac(currentKey, plainText, "request")`, check: `new Uint8Array(mod.generateMac(` + jsCk + `, "4012345678909D987", "request")).toBase64() === "nMx4Fz/E+2Q="`},
		{name: `generateMac(currentKey, plainText, "response")`, check: `new Uint8Array(mod.generateMac(` + jsCk + `, "4012345678909D987", "response")).toBase64() === "IDZCI8H/APo="`},
		{name: "derivationOfInitialKey(bdk, ksn) keeps the input", check: `(() => { const bdk = ` + jsBdk + `; mod.derivationOfInitialKey(bdk, ` + jsKsn + `); return bdk.toBase64() === "ASNFZ4mrze/+3LqYdlQyEA==" })()`},
		{name: "derivationOfInitialKeyAsBase64(bdk, ksn)", check: `mod.derivationOfInitialKeyAsBase64("ASNFZ4mrze/+3LqYdlQyEA==", "//+YdlQyEOAAAQ==") === "asKS+qExW02FirOj19WTOg=="`},
		{name: "deriveCurrentTransactionKeyAsBase64(ik, ksn)", check: `mod.deriveCurrentTransactionKeyAsBase64("asKS+qExW02FirOj19WTOg==", "//+YdlQyEOAAAQ==") === "BCZmtJGEz6No3pYo0Dl7yQ=="`},
		{name: "encryptPinAsBase64(currentKey, pin, pan, format)", check: `mod.encryptPinAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "1234", "4012345678909", "ISO-0") === "G5wYReuZOno="`},
		{name: "decryptPinAsBase64(currentKey, ciphertext, pan, format)", check: `mod.decryptPinAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "G5wYReuZOno=", "4012345678909", "ISO-0") === "1234"`},
		{name: "encryptDataAsBase64(currentKey, iv, plainText, action)", check: `mod.encryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "", "4012345678909D987", "request") === "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP"`},
		{name: "decryptDataAsBase64(currentKey, ciphertext, iv, action)", check: `mod.decryptDataAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "/A1Tt+of2p7miq8ucNm5UGIpviqpk/BP", "", "request").slice(0, 17) === "4012345678909D987"`},
		{name: `generateMacAsBase64(currentKey, plainText, "request")`, check: `mod.generateMacAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "4012345678909D987", "request") === "nMx4Fz/E+2Q="`},
		{name: `generateMacAsBase64(currentKey, plainText, "response")`, check: `mod.generateMacAsBase64("BCZmtJGEz6No3pYo0Dl7yQ==", "4012345678909D987", "response") === "IDZCI8H/APo="`},
		{name: "derivationOfInitialKeyAsHex(bdk, ksn)", check: `mod.derivationOfInitialKeyAsHex("0123456789ABCDEFFEDCBA9876543210", "FFFF9876543210E00001") === "6AC292FAA1315B4D858AB3A3D7D5933A"`},
		{name: "deriveCurrentTransactionKeyAsHex(ik, ksn)", check: `mod.deriveCurrentTransactionKeyAsHex("6AC292FAA1315B4D858AB3A3D7D5933A", "FFFF9876543210E00001") === "042666B49184CFA368DE9628D0397BC9"`},
		{name: "encryptPinAsHex(currentKey, pin, pan, format)", check: `mod.encryptPinAsHex("042666B49184CFA368DE9628D0397BC9", "1234", "4012345678909", "ISO-0") === "1B9C1845EB993A7A"`},
		{name: "decryptPinAsHex(currentKey, ciphertext, pan, format)", check: `mod.decryptPinAsHex("042666B49184CFA368DE9628D0397BC9", "1B9C1845EB993A7A", "4012345678909", "ISO-0") === "1234"`},
		{name: "encryptDataAsHex(currentKey, iv, plainText, action)", check: `mod.encryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "", "4012345678909D987", "request") === "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F"`},
		{name: "decryptDataAsHex(currentKey, ciphertext, iv, action)", check: `mod.decryptDataAsHex("042666B49184CFA368DE9628D0397BC9", "FC0D53B7EA1FDA9EE68AAF2E70D9B9506229BE2AA993F04F", "", "request").slice(0, 17) === "4012345678909D987"`},
		{name: `generateMacAsHex(currentKey, plainText, "request")`, check: `mod.generateMacAsHex("042666B49184CFA368DE9628D0397BC9", "4012345678909D987", "request") === "9CCC78173FC4FB64"`},
		{name: `generateMacAsHex(currentKey, plainText, "response")`, check: `mod.generateMacAsHex("042666B49184CFA368DE9628D0397BC9", "4012345678909D987", "response") === "20364223C1FF00FA"`},
	}
	for _, tt := range tests { //nolint:paralleltest
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.RunOnEventLoop(tt.check)

			require.NoError(t, err)
			require.True(t, got.ToBoolean())
		})
	}
}
