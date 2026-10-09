/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * k6 API for [pkg.GetDesTcFromKsn] port from moov-io,
 * gets the transaction counter from the key serial number.
 * @param ksn The key serial number as a base64 string. Must be at least 4 bytes (stricter than moov-io, which reads counter 0).
 * @returns The 21 bit transaction counter.
 */
export declare function getDesTcFromKsnAsBase64(ksn: string): number;

/**
 * k6 API for [pkg.GenerateNextDesKsn] port from moov-io,
 * generates the next key serial number as a base64 encoded string.
 * @param ksn The key serial number as a base64 string. Must be at least 4 bytes (stricter than moov-io).
 * @returns The next key serial number as a base64 string.
 */
export declare function generateNextDesKsnAsBase64(ksn: string): string;
