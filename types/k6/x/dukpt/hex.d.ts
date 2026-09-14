/*
 * SPDX-FileCopyrightText: 2026 arukiidou <arukiidou@yahoo.co.jp>
 * SPDX-License-Identifier: Apache-2.0
 */

/**
 * k6 API for [pkg.GetDesTcFromKsn] port from moov-io,
 * gets the transaction counter from the key serial number.
 * KSN length is not validated.
 * @param ksn The key serial number as a hex string (case-insensitive).
 * @returns The 21 bit transaction counter.
 */
export declare function getDesTcFromKsnAsHex(ksn: string): number;
