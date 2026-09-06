/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * **k6 DUKPT extension porting with moov-io**
 *
 * @module "k6/x/dukpt"
 */
export as namespace dukpt;

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
 * k6 API for [des.DerivationOfInitialKey] port from moov-io,
 * derives the initial key.
 * @param bdk The base derivation key.
 * @param ksn The key serial number.
 * @returns The derived initial key.
 */
export declare function derivationOfInitialKey(
  bdk: ArrayBuffer,
  ksn: ArrayBuffer
): ArrayBuffer;

/**
 * k6 API for [des.DeriveCurrentTransactionKey] compat with moov-io,
 * derives the current transaction key.
 * @param ik The initial key.
 * @param ksn The key serial number.
 * @returns The derived transaction key.
 */
export declare function deriveCurrentTransactionKey(
  ik: ArrayBuffer,
  ksn: ArrayBuffer
): ArrayBuffer;
