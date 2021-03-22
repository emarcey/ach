// Licensed to The Moov Authors under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. The Moov Authors licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package examples

import (
	"fmt"
	"log"

	"github.com/moov-io/ach"
)

func Example_bocWriteDebit() {
	fh := mockFileHeader()

	bh := ach.NewBatchHeader()
	bh.ServiceClassCode = ach.DebitsOnly
	bh.CompanyName = "Payee Name"
	bh.CompanyIdentification = fh.ImmediateOrigin
	bh.StandardEntryClassCode = ach.BOC
	bh.CompanyEntryDescription = "ACH BOC"
	bh.EffectiveEntryDate = "190816" // need EffectiveEntryDate to be fixed so it can match output
	bh.ODFIIdentification = "121042882"

	entry := ach.NewEntryDetail()
	entry.TransactionCode = ach.CheckingDebit
	entry.SetRDFI("231380104")
	entry.DFIAccountNumber = "12345678"
	entry.Amount = 250000
	entry.SetCheckSerialNumber("123879654")
	entry.SetReceivingCompany("ABC Company")
	entry.SetTraceNumber(bh.ODFIIdentification, 1)

	// build the batch
	batch := ach.NewBatchBOC(bh)
	batch.AddEntry(entry)
	if err := batch.Create(); err != nil {
		log.Fatalf("Unexpected error building batch: %s\n", err)
	}

	// build the file
	file := ach.NewFile()
	file.SetHeader(fh)
	file.AddBatch(batch)
	if err := file.Create(); err != nil {
		log.Fatalf("Unexpected error building file: %s\n", err)
	}

	fmt.Println(file.Header.String())
	fmt.Println(file.Batches[0].GetHeader().String())
	fmt.Println(file.Batches[0].GetEntries()[0].String())
	fmt.Println(file.Batches[0].GetControl().String())
	fmt.Println(file.Control.String())

	// Output:
	// 101 031300012 2313801041908161055A094101Federal Reserve Bank   My Bank Name           12345678
	// 5225Payee Name                          231380104 BOCACH BOC         190816   1121042880000001
	// 62723138010412345678         0000250000123879654      ABC Company             0121042880000001
	// 82250000010023138010000000250000000000000000231380104                          121042880000001
	// 9000001000001000000010023138010000000250000000000000000
}
