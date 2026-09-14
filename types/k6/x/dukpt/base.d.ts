/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * k6 API for [pkg.GetDesTcFromKsn] port from moov-io,
 * gets the transaction counter from the key serial number.
 * KSN length is not validated.
 * @param ksn The key serial number.
 * @returns The 21 bit transaction counter.
 */
export declare function getDesTcFromKsn(
  ksn: ArrayBuffer | Uint8Array<ArrayBufferLike>
): number;

/**
 * k6 API for [pkg.GenerateNextDesKsn] port from moov-io,
 * generates the next key serial number.
 * The input is not modified.
 * @param ksn The key serial number.
 * @returns The next key serial number.
 */
export declare function generateNextDesKsn(
  ksn: ArrayBuffer | Uint8Array<ArrayBufferLike>
): ArrayBuffer;
