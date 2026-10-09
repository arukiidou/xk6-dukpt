/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * k6 API for [pkg.GetDesTcFromKsn] port from moov-io,
 * gets the transaction counter from the key serial number.
 * @param ksn The key serial number as a hex string (case-insensitive). Must be at least 4 bytes (stricter than moov-io, which reads counter 0).
 * @returns The 21 bit transaction counter.
 */
export declare function getDesTcFromKsnAsHex(ksn: string): number;

/**
 * k6 API for [pkg.GenerateNextDesKsn] port from moov-io,
 * generates the next key serial number as a hex encoded string.
 * @param ksn The key serial number as a hex string (case-insensitive). Must be at least 4 bytes (stricter than moov-io).
 * @returns The next key serial number as an uppercase hex string.
 */
export declare function generateNextDesKsnAsHex(ksn: string): string;
