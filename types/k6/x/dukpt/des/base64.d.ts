/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

import type { DukptAction, PinBlockFormat } from "./base";

/**
 * k6 API for [des.DerivationOfInitialKey] port from moov-io,
 * derives the initial key as a base64 encoded string.
 * @param bdk The base derivation key as a base64 string.
 * @param ksn The key serial number as a base64 string.
 * @returns The derived initial key as a base64 string.
 */
export declare function derivationOfInitialKeyAsBase64(
  bdk: string,
  ksn: string
): string;

/**
 * k6 API for [des.DeriveCurrentTransactionKey] compat with moov-io,
 * derives the current transaction key as a base64 encoded string.
 * @param ik The initial key as a base64 string.
 * @param ksn The key serial number as a base64 string.
 * @returns The derived transaction key as a base64 string.
 */
export declare function deriveCurrentTransactionKeyAsBase64(
  ik: string,
  ksn: string
): string;

/**
 * k6 API for [des.EncryptPin] port from moov-io,
 * encrypts a PIN block using the DUKPT transaction key.
 * @param currentKey The transaction key as a base64 string.
 * @param pin The not formatted pin string.
 * @param pan The not formatted pan string.
 * @param format The pinblock format.
 * @returns The cipher text as a base64 string.
 */
export declare function encryptPinAsBase64(
  currentKey: string,
  pin: string,
  pan: string,
  format: PinBlockFormat
): string;

/**
 * k6 API for [des.DecryptPin] port from moov-io,
 * decrypts a PIN block using the DUKPT transaction key.
 * @param currentKey The transaction key as a base64 string.
 * @param ciphertext The encrypted pin block as a base64 string.
 * @param pan The not formatted pan string.
 * @param format The pinblock format.
 * @returns The pin string.
 */
export declare function decryptPinAsBase64(
  currentKey: string,
  ciphertext: string,
  pan: string,
  format: PinBlockFormat
): string;

/**
 * k6 API for [des.EncryptData] port from moov-io,
 * encrypts transaction data using the DUKPT transaction key.
 * @param currentKey The transaction key as a base64 string.
 * @param iv The initial vector as a base64 string, empty for the default zero vector.
 * @param plainText The transaction request data.
 * @param action The request or response action.
 * @returns The encrypted data as a base64 string, zero padded to a multiple of 8 bytes.
 */
export declare function encryptDataAsBase64(
  currentKey: string,
  iv: string,
  plainText: string,
  action: DukptAction
): string;

/**
 * k6 API for [des.DecryptData] port from moov-io,
 * decrypts transaction data using the DUKPT transaction key.
 * @param currentKey The transaction key as a base64 string.
 * @param ciphertext The encrypted data as a base64 string.
 * @param iv The initial vector as a base64 string, empty for the default zero vector.
 * @param action The request or response action.
 * @returns The transaction request data, zero padded to a multiple of 8 bytes.
 */
export declare function decryptDataAsBase64(
  currentKey: string,
  ciphertext: string,
  iv: string,
  action: DukptAction
): string;

/**
 * k6 API for [des.GenerateMac] port from moov-io,
 * generates a MAC using the DUKPT transaction key.
 * @param currentKey The transaction key as a base64 string.
 * @param plainText The transaction request data.
 * @param action The request or response action.
 * @returns The 8 byte mac as a base64 string. Use the first 4 bytes as the ANSI X9.24-1 mac.
 */
export declare function generateMacAsBase64(
  currentKey: string,
  plainText: string,
  action: DukptAction
): string;
